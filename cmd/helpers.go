package cmd

import (
	"encoding/json"
	"fmt"
)

func dump(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(data))
	return nil
}
