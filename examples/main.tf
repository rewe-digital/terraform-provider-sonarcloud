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

output "groups" {
  value = { for k, group in data.sonarcloud_user_groups.groups.groups : lower(group.name) => group }
}