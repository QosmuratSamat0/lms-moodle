// internal/shared/constants/roles.go
package constants

const (
	RoleStudent = "student"
	RoleTeacher = "teacher"
	RoleManager = "manager"
	RoleAdmin   = "admin"
)

var AllRoles = []string{
	RoleStudent,
	RoleTeacher,
	RoleManager,
	RoleAdmin,
}

func IsValidRole(role string) bool {
	for _, r := range AllRoles {
		if r == role {
			return true
		}
	}
	return false
}
