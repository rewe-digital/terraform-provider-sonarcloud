---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "sonarcloud_quality_gate Data Source - terraform-provider-sonarcloud"
subcategory: ""
description: |-
  This Data Source retrieves a single Quality Gate for the configured Organization.
---

# sonarcloud_quality_gate (Data Source)

This Data Source retrieves a single Quality Gate for the configured Organization.

## Example Usage

```terraform
data "sonarcloud_quality_gate" "awesome" {
  name = "my_awesome_quality_gate"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `name` (String) Name of the Quality Gate

### Optional

- `conditions` (Attributes List) The conditions of this quality gate. (see [below for nested schema](#nestedatt--conditions))

### Read-Only

- `gate_id` (Number) Id created by SonarCloud
- `id` (String) Id for Terraform backend
- `is_built_in` (Boolean) Is this Quality gate built in?
- `is_default` (Boolean) Is this the default Quality gate for this project?

<a id="nestedatt--conditions"></a>
### Nested Schema for `conditions`

Optional:

- `op` (String) Operation on which the metric is evaluated must be either: LT, GT

Read-Only:

- `error` (String) The value on which the condition errors.
- `id` (Number) ID of the Condition.
- `metric` (String) The metric on which the condition is based. Must be one of: https://docs.sonarqube.org/latest/user-guide/metric-definitions/


