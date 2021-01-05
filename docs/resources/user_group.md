# Resource `sonarcloud_user_group`

This resource manages user groups for the organization.

## Example Usage

The following example shows how to create a group and add a specific member to this group.

```hcl
resource "sonarcloud_user_group" "my_group" {
  name        = "my-group"
  description = "My group's description."
}

resource "sonarcloud_user_group_member" "my_group_user" {
  group = sonarcloud_user_group.my_group.name
  login = "some-synchronized-user@github"
}
```

## Schema

### Required

- **name** (String, Required) Name of the user group

### Optional

- **description** (String, Optional) Description for the user group
