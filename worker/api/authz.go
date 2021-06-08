package api

var authzUsers = map[string][]string{
	"client_read":                     []string{"read"},
	"client_write":                    []string{"read", "write"},
	"client_unauthorized_permissions": []string{},
}

var permissions = map[string][]string{
	"/WorkerService/Start":  []string{"write"},
	"/WorkerService/Stop":   []string{"write"},
	"/WorkerService/Stream": []string{"read"},
	"/WorkerService/Query":  []string{"read"},
}

// Can returns whether the specified user can perform the method operation.
func Can(method string, user string) bool {
	methodPermission, ok := permissions[method]
	if !ok {
		return false
	}

	userPermission, ok := authzUsers[user]
	if !ok {
		return false
	}

	for _, up := range userPermission {
		for _, mp := range methodPermission {
			if up == mp {
				return true
			}
		}
	}

	return false
}
