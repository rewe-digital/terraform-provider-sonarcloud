# Data Source `sonarcloud_user_groups`

Data source that retrieves a list of user groups for the configured organization.

## Example usage

```hcl
data "sonarcloud_user_groups" "groups" {}

output "groups" {
  value = { for _, group in data.sonarcloud_user_groups.groups.groups : lower(group.name) => group }
}
```

## Schema

### Read-only

- **groups** (List of Object, Read-only) The groups of this organization. (see [below for nested schema](#nestedatt--groups))

<a id="nestedatt--groups"></a>
### Nested Schema for `groups`

- **default** (Boolean) Whether new members are added to this user group per default or not.
- **description** (String) Description of the user group.
- **id** (Number) Numerical ID of the user group.
- **members_count** (Number) Number of members in this user group.
- **name** (String) Name of the user group.


