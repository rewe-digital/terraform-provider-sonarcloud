resource "sonarcloud_user_token" "test_token" {
  name  = "EXAMPLE_TOKEN"
  login = var.token_owner
}
