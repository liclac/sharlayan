package ssh

import (
	"net"

	"github.com/spf13/afero"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
)

type SSHConnShared struct {
	Server *SSHServer
	Fs     afero.Fs

	ServerConfig *ssh.ServerConfig
}

type SSHConn struct {
	*SSHConnShared
	Conn *ssh.ServerConn
	L    *zap.Logger
}

func accept(shared *SSHConnShared, rawConn net.Conn) error {
	sconn, chans, reqs, err := ssh.NewServerConn(rawConn, shared.ServerConfig)
	if err != nil {
		return err
	}
	c := SSHConn{shared, sconn, shared.Server.L.With(zap.String("user", sconn.User()))}
	go c.serveChans(chans)
	go ssh.DiscardRequests(reqs)
	return nil
}

func (c *SSHConn) serveChans(chans <-chan ssh.NewChannel) {
	for newChan := range chans {
		chType := newChan.ChannelType()
		L := c.L.With(zap.String("type", chType))
		switch chType {
		case "session":
			L.Debug("Opening channel")
			ch, reqs, err := newChan.Accept()
			if err != nil {
				L.Error("Failed to accept channel", zap.Error(err))
				continue
			}
			go c.serveSession(ch, reqs)
		default:
			L.Warn("Refusing to open channel")
			if err := newChan.Reject(ssh.UnknownChannelType, "unknown type: "+chType); err != nil {
				L.Error("Failed to send reply", zap.Error(err))
			}
		}
	}
}

func (c *SSHConn) serveSession(ch ssh.Channel, reqs <-chan *ssh.Request) {
	L := c.L.Named("session")
	defer func() {
		if err := ch.Close(); err != nil {
			L.Error("Failed to close channel", zap.Error(err))
		}
	}()
	for req := range reqs {
		switch req.Type {
		case "subsystem":
			// Payload is 4 length bytes, then an ASCII subsystem ID.
			ok := false
			if len(req.Payload) > 4 {
				// If we have a subsystem with this ID, start it and affirm.
				id := string(req.Payload[4:])
				if sub := c.Server.sub[id]; sub != nil {
					L.Debug("Starting subsystem", zap.String("id", id))
					ok = true
					go sub.Serve(c, ch)
				} else {
					L.Warn("Unknown subsystem requested", zap.String("id", id))
				}
			} else {
				L.Debug("Malformed subsystem ID", zap.String("id", string(req.Payload)))
			}
			reply(L, req, ok, nil)
		default:
			DeclineUnknownRequest(L, req)
		}
	}
}

func DeclineUnknownRequest(L *zap.Logger, req *ssh.Request) {
	L.Debug("Unknown request", zap.String("type", req.Type),
		zap.String("payload", string(req.Payload)))
	if req.WantReply {
		reply(L, req, false, nil)
	}
}

func reply(L *zap.Logger, req *ssh.Request, ok bool, payload []byte) {
	if err := req.Reply(ok, payload); err != nil {
		L.Warn("Failed to send reply", zap.String("type", req.Type), zap.Error(err))
	}
}
