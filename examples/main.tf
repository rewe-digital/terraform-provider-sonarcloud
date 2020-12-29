terraform {
  required_providers {
    sonarcloud = {
      version  = "~> 0.1"
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

resource "sonarcloud_user_group_member" "example_member" {
  group = sonarcloud_user_group.test_group.name
  login = var.test_member_login
}

output "groups" {
  value = { for k, group in data.sonarcloud_user_groups.groups.groups : lower(group.name) => group }
}

resource "sonarcloud_permission" "global" {
  group = sonarcloud_user_group.test_group.name
  permissions = ["scan"]
}

resource "sonarcloud_permission" "project" {
  project = "${var.organization}_test"
  group = sonarcloud_user_group.test_group.name
  permissions = ["admin", "scan"]
}
