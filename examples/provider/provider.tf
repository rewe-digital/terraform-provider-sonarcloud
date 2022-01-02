terraform {
  required_providers {
    sonarcloud = {
      source  = "rewe-digital/sonarcloud"
      version = "0.1.1"
    }
  }
}

provider "sonarcloud" {
  organization = var.organization
  token        = var.token
}
