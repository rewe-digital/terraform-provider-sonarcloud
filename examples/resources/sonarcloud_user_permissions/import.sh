# import user permissions for the whole organization using <login>
terraform import "sonarcloud_user_permissions.example_user" "user@github.com"

# import user permissions for a specific project using <login>,<project_key>
terraform import "sonarcloud_user_permissions.example_user" "user@github.com,example_project"