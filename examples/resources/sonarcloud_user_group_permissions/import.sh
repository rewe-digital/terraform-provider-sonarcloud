# import user group permissions for the whole organization using <name>
terraform import "sonarcloud_user_group_permissions.example_group" "Example Group"

# import user group permissions for a specific project using <name>,<project_key>
terraform import "sonarcloud_user_group_permissions.example_group" "Example Group,example_project"