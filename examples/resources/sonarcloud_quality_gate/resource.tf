resource "sonarcloud_quality_gate" "awesome" {
  name      = "My Awesome Quality Gate"
  isDefault = true
  conditions = [
    // Less than 100% coverage
    {
      metric = "coverage"
      error  = 100
      Op     = "LT"
    },
    // Less than 100% coverage on new code
    {
      metric = "new_coverage"
      error  = 100
      Op     = "LT"
    }
  ]

}