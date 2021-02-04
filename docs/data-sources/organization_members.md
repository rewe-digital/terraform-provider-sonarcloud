# Data Source `sonarcloud_organization_members`

Data source that retrieves a list of user that are a member of the configured organization.

## Example usage

```hcl
data "sonarcloud_organization_members" "users" {}

output "users" {
  value = data.sonarcloud_organization_members.users
}
```

## Schema

### Read-only

- **users** (List of Object, Read-only) The members of this organization. (see [below for nested schema](#nestedatt--users))

<a id="nestedatt--users"></a>
### Nested Schema for `users`

- **login** (String) Login of the user.


