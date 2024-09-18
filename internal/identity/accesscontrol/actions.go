package accesscontrol

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// Define permissions
var rolePermissions = map[string]map[string]bool{
	RoleAdmin: {
		"create": true,
		"read":   true,
		"update": true,
		"delete": true,
	},
	RoleUser: {
		"read": true,
	},
}

// Check if role has permission
func HasPermission(role, permission string) bool {
	perms, ok := rolePermissions[role]
	if !ok {
		return false
	}

	return perms[permission]
}
