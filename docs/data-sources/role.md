---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gocd_role Data Source - terraform-provider-gocd"
subcategory: ""
description: |-

---

# gocd_role (Data Source)
Fetches the role information of specified role from GoCD by interacting with GoCD roles [api](https://api.gocd.org/current/#roles).

## Example Usage
```terraform
data "gocd_role" "sample" {
    name = "sample"
}
```


<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) The name of the role.

### Optional

- `auth_config_id` (String) The authorization configuration identifier.
- `etag` (String) Etag used to track the role
- `policy` (List of Map of String) Policy is fine-grained permissions attached to the users belonging to the current role.
- `properties` (Block List) Attributes are used to describes the configuration for gocd role or plugin role. (see [below for nested schema](#nestedblock--properties))
- `type` (String) Type of the role. Use GoCD to create core role and plugin to create plugin role.
- `users` (List of String) The list of users belongs to the role.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--properties"></a>
### Nested Schema for `properties`

Required:

- `key` (String) the name of the property key.

Optional:

- `encrypted_value` (String) The encrypted value of the property
- `is_secure` (Boolean) Specify whether the given property is secure or not. If true and encrypted_value is not specified, GoCD will store the value in encrypted format.
- `value` (String) The value of the property

