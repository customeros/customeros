package enum

type UserRole string

const (
	UserRoleUser          UserRole = "USER"
	UserRoleOwner         UserRole = "OWNER"
	UserRoleAdmin         UserRole = "ADMIN"
	UserRolePlatformOwner UserRole = "PLATFORM_OWNER"
	UserRoleImpersonated  UserRole = "IMPERSONATED"
)

var AllUserRoles = []UserRole{
	UserRoleUser,
	UserRoleOwner,
	UserRoleAdmin,
	UserRolePlatformOwner,
	UserRoleImpersonated,
}

func (r UserRole) String() string {
	return string(r)
}

func DecodeUserRole(s string) UserRole {
	if IsValidUserRole(s) {
		return UserRole(s)
	}
	return UserRoleUser // Default to USER if invalid
}

func IsValidUserRole(s string) bool {
	for _, role := range AllUserRoles {
		if role == UserRole(s) {
			return true
		}
	}
	return false
}
