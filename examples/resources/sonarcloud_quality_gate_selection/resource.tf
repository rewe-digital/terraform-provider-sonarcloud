data "sonarcloud_quality_gate" "awesome_qg" {
  name = "my_awesome_quality_gate"
}

data "sonarcloud_projects" "all" {}

resource "sonarcloud_quality_gate_selection" "example_quality_gate_selection" {
  gate_id   = data.sonarcloud_quality_gate.awesome_qg.gate_id
  project_keys = [for project in data.sonarcloud_projects.all : project.key if project.name == "My Awesome Project"]
}