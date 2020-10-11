terraform {
  required_providers {
    sonarcloud = {
      versions = ["0.1"]
      source   = "rewe-digital.com/platform/sonarcloud"
    }
  }
}

provider "sonarcloud" {
  organization = var.organization
  token        = var.token
}

data "sonarcloud_user_groups" "groups" {}

resource "sonarcloud_user_group" "test_group" {
  name        = "example_group"
  description = "Example group"
}

output "groups" {
  value = { for k, group in data.sonarcloud_user_groups.groups.groups : lower(group.name) => group }
}