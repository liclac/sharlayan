package server

import (
	"context"
	"fmt"
	"sync"

	"github.com/spf13/afero"
)

// A Server function implements serving on a specific protocol - see eg. ServeHTTP and ServeSFTP.
// The server is responsible for its own listener and shutting down when the context expires.
type Server interface {
	Run(ctx context.Context, fs afero.Fs) error
}

// Runs several servers concurrently, until the context expires or one returns. Once one returns,
// the context given to all servers is cancelled. A nil server is silently ignored - this is the
// usual way to disable a server, as just returning immediately would trigger a shutdown, and
// sleeping until the context returns is a waste of a perfectly good goroutine.
func Serve(ctx context.Context, fs afero.Fs, srvs ...Server) <-chan error {
	errC := make(chan error, len(srvs)) // Buffer prevents deadlock if not fully consumed.

	go func() {
		defer close(errC) // Prevent the channel from leaking.

		ctx, cancel := context.WithCancel(ctx)
		defer cancel() // Prevent a context goroutine leak.

		var wg sync.WaitGroup
		defer wg.Wait() // Wait for all goroutines to terminate.

		var num int // Number of non-nil servers.
		for _, srv := range srvs {
			if srv == nil {
				continue
			}
			num++

			wg.Add(1) // Inform the group of an outstanding task.
			go func(srv Server) {
				defer cancel()  // Cancel the context if an error is received.
				defer wg.Done() // Signal task completion.
				if err := srv.Run(ctx, fs); err != nil {
					errC <- err // Buffered channel makes this non-blocking.
				}
			}(srv)
		}

		// If no servers are enabled, return an error, else it'd be confusing.
		if num == 0 {
			errC <- fmt.Errorf("no protocols are enabled")
		}
	}()

	return errC
}
