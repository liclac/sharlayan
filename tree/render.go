package tree

import (
	"encoding/json"
)

var _ Node = JSONNode{}

type JSONNode struct {
	NodeInfo
	V interface{}
}

func JSON(info NodeInfo, v interface{}) JSONNode {
	return JSONNode{info, v}
}

func (n JSONNode) Render(t Tree, self *NodeWrapper, path string) error {
	f, err := t.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(n.V)
}
