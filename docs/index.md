# SonarCloud Provider

This provider manages resources likes user groups and permissions for an organization in SonarCloud.

## Example usage

```hcl
provider "sonarcloud" {
  organization = var.organization
  token        = var.token
}
```

## Argument Reference

- **organization** (String, Required) The name of the organization which you want to manage. 
- **token** (String, Required) The token used to communicate with the SonarCloud API. For most actions this token should 
    belong to an account with admin permissions on the organization.
