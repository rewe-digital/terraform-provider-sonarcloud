package api

// AUTOMATICALLY GENERATED, DO NOT EDIT BY HAND!

// PermissionsAddGroup: Add permission to a group.<br /> This service defaults to global permissions, but can be limited to project permissions by providing project id or project key.<br /> The group name or group id must be provided. <br />Requires the permission 'Administer' on the specified project.
type PermissionsAddGroup struct {
	GroupId      string `form:"groupId,omitempty"`      // Group id
	GroupName    string `form:"groupName,omitempty"`    // Group name or 'anyone' (case insensitive)
	Organization string `form:"organization,omitempty"` // Key of organization, used when group name is set
	Permission   string `form:"permission,omitempty"`   // Permission<ul><li>Possible values for global permissions: admin, profileadmin, gateadmin, scan, provisioning</li><li>Possible values for project permissions admin, codeviewer, issueadmin, securityhotspotadmin, scan, user</li></ul>
	ProjectId    string `form:"projectId,omitempty"`    // Project id
	ProjectKey   string `form:"projectKey,omitempty"`   // Project key
}

// PermissionsAddGroupToTemplate: Add a group to a permission template.<br /> The group id or group name must be provided. <br />Requires the permission 'Administer' on the organization.
type PermissionsAddGroupToTemplate struct {
	GroupId      string `form:"groupId,omitempty"`      // Group id
	GroupName    string `form:"groupName,omitempty"`    // Group name or 'anyone' (case insensitive)
	Organization string `form:"organization,omitempty"` // Key of organization, used when group name is set
	Permission   string `form:"permission,omitempty"`   // Permission<ul><li>Possible values for project permissions admin, codeviewer, issueadmin, securityhotspotadmin, scan, user</li></ul>
	TemplateId   string `form:"templateId,omitempty"`   // Template id
	TemplateName string `form:"templateName,omitempty"` // Template name
}

// PermissionsAddProjectCreatorToTemplate: Add a project creator to a permission template.<br>Requires the permission 'Administer' on the organization.
type PermissionsAddProjectCreatorToTemplate struct {
	Organization string `form:"organization,omitempty"` // Key of organization, used when group name is set
	Permission   string `form:"permission,omitempty"`   // Permission<ul><li>Possible values for project permissions admin, codeviewer, issueadmin, securityhotspotadmin, scan, user</li></ul>
	TemplateId   string `form:"templateId,omitempty"`   // Template id
	TemplateName string `form:"templateName,omitempty"` // Template name
}

// PermissionsAddUser: Add permission to a user.<br /> This service defaults to global permissions, but can be limited to project permissions by providing project id or project key.<br />Requires the permission 'Administer' on the specified project.
type PermissionsAddUser struct {
	Login        string `form:"login,omitempty"`        // User login
	Organization string `form:"organization,omitempty"` // Key of organization, used when group name is set
	Permission   string `form:"permission,omitempty"`   // Permission<ul><li>Possible values for global permissions: admin, profileadmin, gateadmin, scan, provisioning</li><li>Possible values for project permissions admin, codeviewer, issueadmin, securityhotspotadmin, scan, user</li></ul>
	ProjectId    string `form:"projectId,omitempty"`    // Project id
	ProjectKey   string `form:"projectKey,omitempty"`   // Project key
}

// PermissionsAddUserToTemplate: Add a user to a permission template.<br /> Requires the permission 'Administer' on the organization.
type PermissionsAddUserToTemplate struct {
	Login        string `form:"login,omitempty"`        // User login
	Organization string `form:"organization,omitempty"` // Key of organization, used when group name is set
	Permission   string `form:"permission,omitempty"`   // Permission<ul><li>Possible values for project permissions admin, codeviewer, issueadmin, securityhotspotadmin, scan, user</li></ul>
	TemplateId   string `form:"templateId,omitempty"`   // Template id
	TemplateName string `form:"templateName,omitempty"` // Template name
}

// PermissionsApplyTemplate: Apply a permission template to one project.<br>The project id or project key must be provided.<br>The template id or name must be provided.<br>Requires the permission 'Administer' on the organization.
type PermissionsApplyTemplate struct {
	Organization string `form:"organization,omitempty"` // Key of organization, used when group name is set
	ProjectId    string `form:"projectId,omitempty"`    // Project id
	ProjectKey   string `form:"projectKey,omitempty"`   // Project key
	TemplateId   string `form:"templateId,omitempty"`   // Template id
	TemplateName string `form:"templateName,omitempty"` // Template name
}

// PermissionsBulkApplyTemplate: Apply a permission template to several projects.<br />The template id or name must be provided.<br />Requires the permission 'Administer' on the organization.
type PermissionsBulkApplyTemplate struct {
	AnalyzedBefore    string `form:"analyzedBefore,omitempty"`    // Filter the projects for which last analysis is older than the given date (exclusive).<br> Either a date (server timezone) or datetime can be provided.
	OnProvisionedOnly string `form:"onProvisionedOnly,omitempty"` // Filter the projects that are provisioned
	Organization      string `form:"organization,omitempty"`      // Key of organization, used when group name is set
	Projects          string `form:"projects,omitempty"`          // Comma-separated list of project keys
	Q                 string `form:"q,omitempty"`                 // Limit search to: <ul><li>project names that contain the supplied string</li><li>project keys that are exactly the same as the supplied string</li></ul>
	Qualifiers        string `form:"qualifiers,omitempty"`        // Comma-separated list of component qualifiers. Filter the results with the specified qualifiers. Possible values are:<ul><li>TRK - Projects</li></ul>
	TemplateId        string `form:"templateId,omitempty"`        // Template id
	TemplateName      string `form:"templateName,omitempty"`      // Template name
}

// PermissionsCreateTemplate: Create a permission template.<br />Requires the permission 'Administer' on the organization.
type PermissionsCreateTemplate struct {
	Description       string `form:"description,omitempty"`       // Description
	Name              string `form:"name,omitempty"`              // Name
	Organization      string `form:"organization,omitempty"`      // Key of organization, used when group name is set
	ProjectKeyPattern string `form:"projectKeyPattern,omitempty"` // Project key pattern. Must be a valid Java regular expression
}

// PermissionsDeleteTemplate: Delete a permission template.<br />Requires the permission 'Administer' on the organization.
type PermissionsDeleteTemplate struct {
	Organization string `form:"organization,omitempty"` // Key of organization, used when group name is set
	TemplateId   string `form:"templateId,omitempty"`   // Template id
	TemplateName string `form:"templateName,omitempty"` // Template name
}

// PermissionsRemoveGroup: Remove a permission from a group.<br /> This service defaults to global permissions, but can be limited to project permissions by providing project id or project key.<br /> The group id or group name must be provided, not both.<br />Requires the permission 'Administer' on the specified project.
type PermissionsRemoveGroup struct {
	GroupId      string `form:"groupId,omitempty"`      // Group id
	GroupName    string `form:"groupName,omitempty"`    // Group name or 'anyone' (case insensitive)
	Organization string `form:"organization,omitempty"` // Key of organization, used when group name is set
	Permission   string `form:"permission,omitempty"`   // Permission<ul><li>Possible values for global permissions: admin, profileadmin, gateadmin, scan, provisioning</li><li>Possible values for project permissions admin, codeviewer, issueadmin, securityhotspotadmin, scan, user</li></ul>
	ProjectId    string `form:"projectId,omitempty"`    // Project id
	ProjectKey   string `form:"projectKey,omitempty"`   // Project key
}

// PermissionsRemoveGroupFromTemplate: Remove a group from a permission template.<br /> The group id or group name must be provided. <br />Requires the permission 'Administer' on the organization.
type PermissionsRemoveGroupFromTemplate struct {
	GroupId      string `form:"groupId,omitempty"`      // Group id
	GroupName    string `form:"groupName,omitempty"`    // Group name or 'anyone' (case insensitive)
	Organization string `form:"organization,omitempty"` // Key of organization, used when group name is set
	Permission   string `form:"permission,omitempty"`   // Permission<ul><li>Possible values for project permissions admin, codeviewer, issueadmin, securityhotspotadmin, scan, user</li></ul>
	TemplateId   string `form:"templateId,omitempty"`   // Template id
	TemplateName string `form:"templateName,omitempty"` // Template name
}

// PermissionsRemoveProjectCreatorFromTemplate: Remove a project creator from a permission template.<br>Requires the permission 'Administer' on the organization.
type PermissionsRemoveProjectCreatorFromTemplate struct {
	Organization string `form:"organization,omitempty"` // Key of organization, used when group name is set
	Permission   string `form:"permission,omitempty"`   // Permission<ul><li>Possible values for project permissions admin, codeviewer, issueadmin, securityhotspotadmin, scan, user</li></ul>
	TemplateId   string `form:"templateId,omitempty"`   // Template id
	TemplateName string `form:"templateName,omitempty"` // Template name
}

// PermissionsRemoveUser: Remove permission from a user.<br /> This service defaults to global permissions, but can be limited to project permissions by providing project id or project key.<br /> Requires the permission 'Administer' on the specified project.
type PermissionsRemoveUser struct {
	Login        string `form:"login,omitempty"`        // User login
	Organization string `form:"organization,omitempty"` // Key of organization, used when group name is set
	Permission   string `form:"permission,omitempty"`   // Permission<ul><li>Possible values for global permissions: admin, profileadmin, gateadmin, scan, provisioning</li><li>Possible values for project permissions admin, codeviewer, issueadmin, securityhotspotadmin, scan, user</li></ul>
	ProjectId    string `form:"projectId,omitempty"`    // Project id
	ProjectKey   string `form:"projectKey,omitempty"`   // Project key
}

// PermissionsRemoveUserFromTemplate: Remove a user from a permission template.<br /> Requires the permission 'Administer' on the organization.
type PermissionsRemoveUserFromTemplate struct {
	Login        string `form:"login,omitempty"`        // User login
	Organization string `form:"organization,omitempty"` // Key of organization, used when group name is set
	Permission   string `form:"permission,omitempty"`   // Permission<ul><li>Possible values for project permissions admin, codeviewer, issueadmin, securityhotspotadmin, scan, user</li></ul>
	TemplateId   string `form:"templateId,omitempty"`   // Template id
	TemplateName string `form:"templateName,omitempty"` // Template name
}

// PermissionsSearchGlobalPermissions: List global permissions. <br />Requires the following permission: 'Administer System'
type PermissionsSearchGlobalPermissions struct {
	Organization string `form:"organization,omitempty"` // Key of organization, used when group name is set
}

// PermissionsSearchProjectPermissions: List project permissions. A project can be a technical project, a view or a developer.<br />Requires the permission 'Administer' on the specified project.
type PermissionsSearchProjectPermissions struct {
	P          string `form:"p,omitempty"`          // 1-based page number
	ProjectId  string `form:"projectId,omitempty"`  // Project id
	ProjectKey string `form:"projectKey,omitempty"` // Project key
	Ps         string `form:"ps,omitempty"`         // Page size. Must be greater than 0.
	Q          string `form:"q,omitempty"`          // Limit search to: <ul><li>project names that contain the supplied string</li><li>project keys that are exactly the same as the supplied string</li></ul>
	Qualifier  string `form:"qualifier,omitempty"`  // Project qualifier. Filter the results with the specified qualifier. Possible values are:<ul><li>TRK - Projects</li></ul>
}

// PermissionsSearchTemplates: List permission templates.<br />Requires the permission 'Administer' on the organization.
type PermissionsSearchTemplates struct {
	Organization string `form:"organization,omitempty"` // Key of organization, used when group name is set
	Q            string `form:"q,omitempty"`            // Limit search to permission template names that contain the supplied string.
}

// PermissionsSetDefaultTemplate: Set a permission template as default.<br />Requires the permission 'Administer' on the organization.
type PermissionsSetDefaultTemplate struct {
	Organization string `form:"organization,omitempty"` // Key of organization, used when group name is set
	Qualifier    string `form:"qualifier,omitempty"`    // Project qualifier. Filter the results with the specified qualifier. Possible values are:<ul><li>TRK - Projects</li></ul>
	TemplateId   string `form:"templateId,omitempty"`   // Template id
	TemplateName string `form:"templateName,omitempty"` // Template name
}

// PermissionsUpdateTemplate: Update a permission template.<br />Requires the permission 'Administer' on the organization.
type PermissionsUpdateTemplate struct {
	Description       string `form:"description,omitempty"`       // Description
	Id                string `form:"id,omitempty"`                // Id
	Name              string `form:"name,omitempty"`              // Name
	ProjectKeyPattern string `form:"projectKeyPattern,omitempty"` // Project key pattern. Must be a valid Java regular expression
}
