# Resource `sonarcloud_user_group_permissions`

This resource manages the permissions of a user group for the whole organization or a specific project.

## Example Usage

### Global permissions

```hcl
resource "sonarcloud_user_group_permission" "global" {
  group = "my-group"
  permissions = ["scan"]
}
```

-> Possible values for global permissions are `admin`, `profileadmin`, `gateadmin`, `scan`, and `provisioning`.

### Permissions on a specific project

```hcl
resource "sonarcloud_user_group_permission" "project" {
  project = "my-project"
  group = "my-group"
  permissions = ["admin", "scan"]
}
```

->Possible values for project permissions are `admin`, `codeviewer`, `issueadmin`, `securityhotspotadmin`, `scan`, 
and `user`.

!>Some permissions (like codeviewer) are active by default on public projects and are not returned when reading the resource.
Using those on a public project will result in an unstable state when applying. 

## Schema

### Required

- **group** (String, Required) User group to set the permissions for.
- **permissions** (List of String, Required) List of permissions to grant.

### Optional

- **project** (String, Optional) The key of the project to restrict the permissions to.


