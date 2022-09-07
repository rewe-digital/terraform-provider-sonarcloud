resource "sonarcloud_quality_gate" "awesome" {
  name       = "My Awesome Quality Gate"
  is_default = true
  conditions = [
    // Less than 100% coverage
    {
      metric = "coverage"
      error  = 100
      op     = "LT"
    },
    // Less than 100% coverage on new code
    {
      metric = "new_coverage"
      error  = 100
      op     = "LT"
    }
  ]

}
