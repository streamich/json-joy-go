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
		fmt.Fprintf(os.Stderr, "Could not read input: %v\n", err)
		os.Exit(1)
	}

	var doc interface{}
	err = json.Unmarshal(bytes, &doc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid input JSON: %v\n", err)
		os.Exit(1)
	}

	if (len(os.Args)) < 2 {
		fmt.Fprintf(os.Stderr, "JSON Patch argument not provided")
		os.Exit(1)
	}

	var patch interface{}
	err = json.Unmarshal([]byte(os.Args[1]), &patch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid patch JSON: %v\n", err)
		os.Exit(1)
	}

	ops, _, err := jsonjoy.CreateOps(patch)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	err = jsonjoy.ApplyOps(&doc, ops)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	result, err := json.Marshal(doc)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not serialize output: %v\n", err)
		os.Exit(1)
	}

	os.Stdout.Write(result)
	os.Stdout.WriteString("\n")
}
