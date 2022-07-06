package resources

type ComponentManifest struct {
	Deployment          string
	Service             string
	ServiceAccount      string
	Role                []string
	RoleBinding         []string
	ClusterRoles        []string
	ClusterRoleBindings []string
}
