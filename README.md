# terraform-provider-sonarcloud

A Terraform provider for managing SonarCloud groups and users.

## Installing

Run `make install` to build the terraform provider and store it in `~/.terraform/plugins`. 

## Testing

Run `make test` to run all unit tests. This should work without further config and not touch any infrastructure.

Run 'make testacc' to run all acceptance tests. This relies on quite a specific test-organization being available in SonarCloud.
The project should have: 

- 3 groups:
  - Members (Default) - with 2 members
  - Owners (Should be auto-created as well) - with 1 member
  - TEST_DONT_REMOVE - with 0 members (see environment variables below to see how to customize this)
  
 Set the following environment variables before running the acceptance test: 
| Variable | Description |
|---|---|
| `SONARCLOUD_ORGANIZATION` | The name of the org to run tests against. |
| `SONARCLOUD_TOKEN` | A token with admin permissions for the organization. |
| `SONARCLOUD_TEST_USER_LOGIN` | The login for testing `sonarcloud_user_group_member`. Must be an existing member of the org and in the form of `<github_handle>@github` if you have imported the user via GitHub. |
| `SONARCLOUD_TEST_GROUP_NAME` | The name of an existing group to which the test-user will be added and removed from. |
 
 ## Development
 
 The source files for sending requests to SonarCloud are generated and stored in `pkg/api`.
 The API is generated based on the contents of `gen/services.json`, which is the output of `https://sonarcloud.io/api/webservices/list`.
 See the `AllowedEndpoints` in `gen/main.go` for the list of endpoints that is used for creating API source files.
 Run `make gen` to (re)generate the API source files.
