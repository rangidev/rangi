package admin

import "fmt"

const (
	LoginPath     = "/admin/login"
	DashboardPath = "/admin/dashboard"

	collectionPathTemplate = "/admin/collections/%s"
)

func CollectionPath(collectionName string) string {
	return fmt.Sprintf(collectionPathTemplate, collectionName)
}
