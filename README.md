# Ansible Terraform AWS - Dynamic Inventory
Combining Terraform an Ansible together greats a powerful team, but they're unlinkable without Ansible being able to read Terraform's state files in the form of an inventory. This code will allow you to use your AWS focused Terraform state files as dynamic inventories within Ansible, with some extra magic thrown in.

## Overview of Features
* Reads in `.tfstate` files and produces an Ansible inventory;
* Can download the state from an AWS S3 bucket (and Terraform can put one there for you);
* Has a simple JSON configuration file for managing some options;
* Has configurable tags for detecting the instance human-readable name and Ansible group;

## Installation
You will need to build the binaries your self, for your required platforms, using a Go installation. I do not have the capacility or time to manage pre-compiled binaries, unless someone can make a suggestion for making that an easy task? Perhaps some TravisCI work? Let me know.

`go install github.com/mrcrilly/ansible-inventory-terraform-aws`

You should now have `$GOBIN/ansible-inventory-terraform-aws` or some binary that matches your platform.

## Usage
Ansible expects two flags to be present on a dynamic inventory:

* `--list`
* `--host <inventory-hostname>`

These are the only flags provided to you. Further enhancements are made via the configuration file, as documented below.

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

That being said, we don't currently set the `ansible_ssh_host` variable to whatever we find as this is rather intrusive and which IP do we pick? Instead you should use the host variables provided above to set the host IP you want to Ansible to SSH too or DNS, which is more desirable. 

## Configuration
A JSON file can be used to configure how this utility behaves. It looks like this, in its entirety, at this point in time:

```json
{
    "state_file": "/tmp/terraform.tfstate",
    "s3": {
        "bucket_name": "my-awesome-bucket",
        "bucket_key": "states/terraform.tfstate"
    },
    "options": {
        "group_name_tag": "Group",
        "instance_name_tag": "Name"
    }
}
```

The `s3` section can be omitted, and AWS S3 won't be used at all for retreiving state. The options key can also be omitted and defaults will be used in code (`Group` and `Name`.)

The `state_file` key is required and the file must exist already, and be a valid Terraform state, if you're not pulling it from S3.

## Author
Michael Crilly.

## Licence
MIT