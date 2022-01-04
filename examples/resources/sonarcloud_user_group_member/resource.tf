data "sonarcloud_user_group" "awesome" {
  name = "My awesome group"
}

resource "sonarcloud_user_group_member" "example_member" {
  group = data.sonarcloud_user_group.awesome.name
  login = var.example_member_login
}

