package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	jsonjoy "github.com/streamich/json-joy-go"
)

func main() {
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	var doc interface{}
	json.Unmarshal(bytes, &doc)

	if (len(os.Args)) < 2 {
		fmt.Fprintf(os.Stderr, "error: JSON Pointer argument not provided")
		os.Exit(1)
	}

	tokens, err := jsonjoy.NewJSONPointer(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	value, err := tokens.Get(doc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	result, err := json.Marshal(value)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	os.Stdout.Write(result)
	os.Stdout.WriteString("\n")
}
