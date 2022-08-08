# Terraform Provider SonarCloud

A Terraform provider for managing SonarCloud user groups and their permissions.

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) >= 1.x
-	[Go](https://golang.org/doc/install) >= 1.17
-   [GoReleaser](https://goreleaser.com/) >= 0.153.x

## Installing the Provider locally

1. Clone the repository
1. Enter the repository directory
1. Run `make install` to build the terraform provider and store it in `~/.terraform/plugins`. 

**Note**: this uses Goreleaser under the hood. Alternatively you can use `go build` and move the binary to the correct location yourself.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

Run `make test` to run all unit tests. This should work without further config and not touch any infrastructure.

Run `make testacc` to run all acceptance tests. This relies on quite a specific test-organization being available in SonarCloud.
The project should have the following 3 groups:
  
- Members (Default) - with 2 members
- Owners (Should be auto-created as well) - with 1 member
- TEST_DONT_REMOVE - with 0 members (see environment variables below to see how to customize this)

Acceptance tests also rely on the following preconfigured SonarCloud Resources:

- One test Quality Gate
- One test Project
  
 Set the following environment variables before running the acceptance test: 

| Variable | Description |
|---|---|
| `SONARCLOUD_ORGANIZATION` | The name of the org to run tests against. |
| `SONARCLOUD_TOKEN` | A token with admin permissions for the organization. |
| `SONARCLOUD_TEST_USER_LOGIN` | The login for testing `sonarcloud_user_group_member`. Must be an existing member of the org and in the form of `<github_handle>@github` if you have imported the user via GitHub. |
| `SONARCLOUD_TEST_GROUP_NAME` | The name of an existing group to which the test-user will be added and removed from. | 
| `SONARCLOUD_TOKEN_TEST_USER_LOGIN` | The login for testing `sonarcloud_user_token`. This must be the login that also has the existing `SONARCLOUD_TOKEN`. |
| `SONARCLOUD_PROJECT_KEY` | The Key of a test `project` for testing the `sonarcloud_quality_gate_selection` resource. |
| `SONARCLOUD_QUALITY_GATE_ID` | The `GateId` of a test `Quality Gate` for testing `sonarcloud_qualtiy_gate_selection` resource. |
| `SONARCLOUD_QUALITY_GATE_NAME` | The `name` of a test `Quality Gate` for testing the `sonarcloud_qualtiy_gate` data source. |
