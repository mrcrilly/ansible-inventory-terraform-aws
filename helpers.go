package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n\n", err.Error())
		os.Exit(1)
	}
}

func outputJson(raw interface{}) error {
	rawToJson, err := json.Marshal(raw)

	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "%s", string(rawToJson))
	return nil
}
