package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
)

import (
	"github.com/hashicorp/terraform/terraform"
)

type AnsibleInventoryGroup struct {
	Hosts     []string          `json:"hosts"`
	Variables map[string]string `json:"vars"`
}

type AnsibleInventory map[string]AnsibleInventoryGroup

var version = flag.Bool("version", false, "print version information and exit")
var list = flag.Bool("list", false, "list mode")
var host = flag.String("host", "", "host mode")

var tagsRegexp = regexp.MustCompile(`^tags\..+`)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}

func parseState(file string) (*terraform.State, error) {
	var tfstate *terraform.State

	fd, err := os.Open(file)

	if err != nil {
		return nil, err
	}

	jsonDecoder := json.NewDecoder(fd)
	err = jsonDecoder.Decode(&tfstate)

	if err != nil {
		return nil, err
	}

	return tfstate, nil
}

func listInventory(tfstate *terraform.State) (map[string]AnsibleInventoryGroup, error) {
	if tfstate == nil {
		return nil, errors.New("No state file provided")
	}

	var inv map[string]AnsibleInventoryGroup = make(map[string]AnsibleInventoryGroup, 0)

	for _, M := range tfstate.Modules {
		for _, R := range M.Resources {
			if R.Type == "aws_instance" {
				var groupName, instanceName string
				var OK bool

				if groupName, OK = R.Primary.Attributes["tags.Group"]; !OK {
					return nil, errors.New("We need to see Group in the instance tags")
				}

				if instanceName, OK = R.Primary.Attributes["tags.Name"]; !OK {
					return nil, errors.New("We need to see Name in the instance tags")
				}

				var ansibleGroup AnsibleInventoryGroup

				if _, OK := inv[groupName]; !OK {
					ansibleGroup = AnsibleInventoryGroup{}
					ansibleGroup.Hosts = []string{}
					ansibleGroup.Variables = make(map[string]string, 1)

					inv[groupName] = ansibleGroup
				} else {
					ansibleGroup = inv[groupName]
				}

				ansibleGroup.Hosts = append(inv[groupName].Hosts, instanceName)
				inv[groupName] = ansibleGroup
			}
		}
	}

	return inv, nil
}

func listHost(tfstate *terraform.State, hostname string) (map[string]string, error) {
	if tfstate == nil {
		return nil, errors.New("No state file provided")
	}

	var inv map[string]string = make(map[string]string, 0)

	for _, M := range tfstate.Modules {
		for _, R := range M.Resources {
			if R.Type == "aws_instance" {
				var instanceName string
				var OK bool

				if instanceName, OK = R.Primary.Attributes["tags.Name"]; !OK {
					return nil, errors.New("We need to see Name in the instance tags")
				}

				if instanceName == hostname {
					if tagCount, OK := R.Primary.Attributes["tags.#"]; OK {
						if tagCount == "0" {
							return inv, nil
						}
					} else {
						return inv, nil // inv == empty map
					}

					for AK, A := range R.Primary.Attributes {
						if tagsRegexp.MatchString(AK) {
							variableName := strings.SplitN(AK, ".", 2)[1]

							if variableName == "#" {
								continue
							}

							inv[variableName] = A
						}
					}
				}
			}
		}
	}

	return inv, nil
}

func outputJson(raw interface{}) error {
	rawToJson, err := json.Marshal(raw)

	if err != nil {
		return err
	}

	fmt.Printf("%s", string(rawToJson))
	return nil
}

func main() {
	flag.Parse()

	var stateFile string
	fromTheEnvironment := os.Getenv("TF_STATE")

	if fromTheEnvironment != "" {
		stateFile = fromTheEnvironment
	} else {
		stateFile = "./terraform.tfstate"
	}

	state, err := parseState(stateFile)
	checkError(err)

	if *list {
		fullInventory, err := listInventory(state)
		checkError(err)

		outputJson(fullInventory)
		os.Exit(0)
	}

	if !*list && *host != "" {
		hostVariables, err := listHost(state, *host)
		checkError(err)

		outputJson(hostVariables)
		os.Exit(0)
	}

	fmt.Println("No action given...")
	os.Exit(1)
}
