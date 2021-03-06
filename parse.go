package main

import (
	"encoding/json"
	"errors"
	"os"
	"regexp"
	"strings"
)

import (
	"github.com/hashicorp/terraform/terraform"
)

var tagsRegexp = regexp.MustCompile(`^tags\..+`)

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

func listInventory(tfstate *terraform.State, expectedGroupName, expectedInstanceName string) (map[string]AnsibleInventoryGroup, error) {
	if tfstate == nil {
		return nil, errors.New("No state file provided")
	}

	var inv map[string]AnsibleInventoryGroup = make(map[string]AnsibleInventoryGroup, 0)

	for _, M := range tfstate.Modules {
		for _, R := range M.Resources {
			if R.Type == "aws_instance" {
				var groupName, instanceName string
				var OK bool

				if groupName, OK = R.Primary.Attributes["tags."+expectedGroupName]; !OK {
					continue
				}

				if instanceName, OK = R.Primary.Attributes["tags."+expectedInstanceName]; !OK {
					continue
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

func listHost(tfstate *terraform.State, hostname, expectedInstanceName string) (map[string]string, error) {
	if tfstate == nil {
		return nil, errors.New("No state file provided")
	}

	var inv map[string]string = make(map[string]string, 0)

	for _, M := range tfstate.Modules {
		for _, R := range M.Resources {
			if R.Type == "aws_instance" {
				var instanceName string
				var OK bool

				if instanceName, OK = R.Primary.Attributes["tags."+expectedInstanceName]; !OK {
					continue
				}

				if instanceName == hostname {
					inv = listHostTags(R.Primary.Attributes)
				}
			}
		}
	}

	return inv, nil
}

func listHostTags(attributes map[string]string) map[string]string {
	if attributes == nil {
		return nil
	}

	var tags map[string]string = make(map[string]string, 0)

	if tagCount, OK := attributes["tags.#"]; OK {
		if tagCount == "0" {
			return nil
		}
	} else {
		return nil
	}

	for AK, A := range attributes {
		if tagsRegexp.MatchString(AK) {
			variableName := strings.SplitN(AK, ".", 2)[1]

			if variableName == "#" {
				continue
			}

			tags[variableName] = A
		}
	}

	tags["private_ip"] = attributes["private_ip"]
	tags["private_dns"] = attributes["private_dns"]
	tags["public_ip"] = attributes["public_ip"]
	tags["public_dns"] = attributes["public_dns"]

	return tags
}
