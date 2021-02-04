# Resource `sonarcloud_user_token`

This resource manages tokens for a user.

## Example Usage

The following example shows how to create token for a user. The token used in the provider must belong to 
the same user as the login in this resource.

```hcl
resource "sonarcloud_user_token" "my_token" {
  login = var.my_login
  name  = "MY_TOKEN"
}
```

## Schema

### Required

- **login** (String, Required) Login of the user to add the token to.
- **name** (String, Required) Name of the token.
