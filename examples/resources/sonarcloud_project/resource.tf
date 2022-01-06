resource "sonarcloud_project" "example_project" {
  key        = "my-unique-project-key"
  name       = "My not-unique project name"
  visibility = "private"
}
