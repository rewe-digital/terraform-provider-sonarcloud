# Resource `sonarcloud_user_group_member`

This resource manages the members of a specific user group.

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

- **login** (String, Required) User login

### Optional

- **group** (String, Optional) Group name
