package roles

const (
	SuperAdmin = "super_admin"
	OrgAdmin   = "org_admin"
	Any        = "any"
	User       = "user"
	Admin      = "admin"
)

func CheckRoles(requiredRoles []string, grantedRoles []string) bool {
	if len(requiredRoles) == 0 {
		return true
	}

	for _, grantedRole := range grantedRoles {
		if grantedRole == SuperAdmin || grantedRole == Admin {
			return true
		}
		for _, requiredRole := range requiredRoles {
			if requiredRole == Any {
				return true
			}
			if grantedRole == requiredRole {
				return true
			}
		}
	}

	return false
}
