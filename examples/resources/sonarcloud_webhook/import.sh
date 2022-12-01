# import a webhook for the whole organization using <id>
terraform import "sonarcloud_webhook.example" "ABCDEFGHIJKLMNOPQRST"

# import a webhook for a specific project using <id>,<project_key>
terraform import "sonarcloud_webhook.example" "ABCDEFGHIJKLMNOPQRST,example_project"