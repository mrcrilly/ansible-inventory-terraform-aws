# Ansible Terraform AWS - Dynamic Inventory
Combining Terraform an Ansible together greats a powerful team, but they're unlinkable without Ansible being able to read Terraform's state files. This code will allow you to use your AWS focused Terraform state files as dynamic inventories within Ansible, with some extra magic thrown in.

## Installation
You will need to build the binaries your self, for your required platforms, using a Go installation. I do not have the capacility or time to manage pre-compiled binaries, unless someone can make a suggestion for making that an easy task? Perhaps some TravisCI work? Let me know.

`go install github.com/mrcrilly/ansible-inventory-terraform-aws`

You should now have `$GOBIN/ansible-inventory-terraform-aws(.exe?)`

## Usage
Ansible expects two flags to be present on a dynamic inventory:

* `--list`
* `--host <inventory-hostname>`

These are the only flags provided to you. Further enhancements are made via environment variables, as documented below.

### Listing the whole inventory

The `--list` flag will generate an Ansible compatiable, complete inventory based on the tags provided. It may look something like this:

```json
{
  "dev-mssql": {
    "hosts": [
      "dev-mssql-001",
      "dev-mssql-002"
    ],
    "vars": {}
  }
}
```

This is inventory's group names and host names are based on two assumptions about how you're tagging your instances in AWS:

* The group name `dev-mssql` is based on a tag called `Group`;
* The host name `dev-mssql-001`, and so forth, is based on a tag called `Name`;

This means when you're creating `aws_instance` resources in your Terraform code, that tags section should look a bit like this:

```hcl
tags {
    Name = "dev-mssql-${format("%03d", count.index + 1)}"
    Group = "dev-mssql"
}
```

This behaviour can be overridden, as discussed below.

You'll also note that `"vars":{}` is an empty hash/map/dict. This is because there's no real sensible way to populate the group variables here. One idea is to perhaps use the same tags/variables that are assigned to hosts when you use `--host` (see below), but this may not be desirable.

### Listing host variables

The `--host` flag will do what you expect it to: it will return the host variables for the hostname provided.

The result will look something like:

```json
{
  "Environment": "dev",
  "Group": "dev-mssql",
  "Name": "dev-mssql-002",
}
```

The host level variables are parsed using a simple regular expression and follow this logic:

* Loop over all the attributes in the state file, looking for `tag.*` (regexp used is: `^tags\..+`);
* Any matches, split the key name, such as `tags.Environment`, and disregard the `tags.` element;
* Store the above value as a key in a dict, setting its value to the value of the attribute from the state file;

For you, this simply means any tags you define in Terraform will be taken as host level variables and provided as such to Ansible. Hopefully this is a sensible way of doing this and you'll agree it's useful.

## IP Addresses
When pulling in the IP addresses from the state file, we simply use, at this point in time, whatever is provided by the following attribute keys:

* private_ip
* private_dns
* public_ip
* public_dns

We basically map these, one-to-one to the host variables returned to Ansible. Also, this should include EIPs you attach to instances via Terraform, as you're reading the state file after everything has been built and linked together.

That being said, we don't currently set the `ansible_ssh_host` variable to whatever we find as this is rather intrusive and which IP do we pick? Instead you should use the host variables provided above to set the host IP you want to Ansible to SSH too.

## Overriding Tags
As weve said above, two tags are important to make this dynamic inventory work: `Name` and `Group`. That being said, these two tags can be overridden using two environment variables:

* `TF_STATE_GROUP_TAG`
* `TF_STATE_INSTANCE_TAG`

Settings to something else, such as `group_name` and `vm_name`, will force this code base ot search `tags.group_name` and `tags.vm_name` when determining what group a VM belongs to, and what that VM is called.

An example of executing this might be:

```
TF_STATE_GROUP_NAME="group_name" ./ansible-inventory-terraform-aws.exe --list
```

This is a classic and well known means of providing a neone-time environment variable value directly to a running command.

### State File Locations
You can also override where the Terraform state file will be located. At this point in time the code base simply looks for `./terraform.tfstate`, but if you use the `TF_STATE` environment variable, and set it to a relative or absolute path, the code base will look there for its state file, open the file, and do its work.

## Author
Michael Crilly.

## Licence
MIT