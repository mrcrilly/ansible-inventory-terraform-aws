package main

import (
	"flag"
	"fmt"
	"os"
)

var list = flag.Bool("list", false, "list mode")
var host = flag.String("host", "", "host mode")

func main() {
	flag.Parse()

	configFile := "./aita.json"
	config, err := parseConfiguration(configFile)
	checkError(err)

	if config.StateFile == "" {
		fmt.Fprint(os.Stderr, "State file location needed.")
		os.Exit(1)
	}

	if config.S3 != nil {
		if config.S3.BucketName != "" && config.S3.BucketKey != "" {
			err := downloadStateFromS3(config.S3.BucketName, config.S3.BucketKey, config.StateFile)
			checkError(err)
		}
	}

	if _, OK := config.Options["instance_name_tag"]; !OK {
		config.Options["instance_name_tag"] = "Name"
	}

	if _, OK := config.Options["group_name_tag"]; !OK {
		config.Options["group_name_tag"] = "Group"
	}

	state, err := parseState(config.StateFile)
	checkError(err)

	if *list {
		fullInventory, err := listInventory(state, config.Options["group_name_tag"], config.Options["instance_name_tag"])
		checkError(err)

		outputJson(fullInventory)
		os.Exit(0)
	}

	if !*list && *host != "" {
		hostVariables, err := listHost(state, *host, config.Options["instance_name_tag"])
		checkError(err)

		outputJson(hostVariables)
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "No action given. Look at the help output.\n\n")
	os.Exit(1)
}
