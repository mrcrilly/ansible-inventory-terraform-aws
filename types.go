package main

type AnsibleInventoryGroup struct {
	Hosts     []string          `json:"hosts"`
	Variables map[string]string `json:"vars"`
}

type AnsibleInventory map[string]AnsibleInventoryGroup
