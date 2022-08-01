package main

import "github.com/nepomueno/ea-cli/cmd"

// Import key modules.

// Define the function to create a resource group.

func main() {
	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}
