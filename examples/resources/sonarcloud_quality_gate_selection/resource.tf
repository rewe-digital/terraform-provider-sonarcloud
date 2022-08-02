data sonarcloud_quality_gate "awesome_qg" {
    name = "my_awesome_quality_gate"
}

data sonarcloud_projects "all" {}

resource sonarcloud_quality_gate_selection "example_quality_gate_selection" {    
    gate_id = awesome_qg.gate_id
    selection = [for project in all : project.key if project.name == "My Awesome Project"]
}