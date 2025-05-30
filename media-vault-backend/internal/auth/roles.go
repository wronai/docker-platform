package auth

// User roles
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// HasAdminRole checks if the user has admin role
func HasAdminRole(roles []string) bool {
	for _, role := range roles {
		if role == RoleAdmin {
			return true
		}
	}
	return false
}

// HasUserRole checks if the user has user role
func HasUserRole(roles []string) bool {
	for _, role := range roles {
		if role == RoleUser {
			return true
		}
	}
	return false
}

// HasAnyRole checks if the user has any of the specified roles
func HasAnyRole(userRoles, requiredRoles []string) bool {
	for _, requiredRole := range requiredRoles {
		for _, userRole := range userRoles {
			if userRole == requiredRole {
				return true
			}
		}
	}
	return false
}