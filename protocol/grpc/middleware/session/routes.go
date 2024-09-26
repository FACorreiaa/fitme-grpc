package session

// rolePermissions Define permissions for each role
var rolePermissions = map[string][]string{
	"USER":  {"view_profile", "view_workouts, calculator_service"},
	"PT":    {"view_profile", "view_workouts", "manage_clients, calculator_service"},
	"ADMIN": {"view_profile", "view_workouts", "manage_clients", "admin_dashboard", "view_all_users, calculator_service"},
}

// MethodPermissions Define required permissions for each gRPC method
var MethodPermissions = map[string]string{
	"/auth.Auth/GetUserInfo":    "view_profile",
	"/auth.Auth/ManageClients":  "manage_clients",
	"/auth.Auth/ViewWorkouts":   "view_workouts",
	"/auth.Auth/AdminDashboard": "admin_dashboard",
	//"/auth.Auth/GetAllUsers":    "view_all_users",
}

func GetUserPermissions(roles []string) []string {
	permissionsSet := make(map[string]bool)
	for _, role := range roles {
		rolePerms, exists := rolePermissions[role]
		if exists {
			for _, perm := range rolePerms {
				permissionsSet[perm] = true
			}
		}
	}

	var permissions []string
	for perm := range permissionsSet {
		permissions = append(permissions, perm)
	}

	return permissions
}
