package deployment

import (
	"time"

	appsv1 "k8s.io/api/apps/v1"
)

// deployCell
type deploymentCell appsv1.Deployment

func (d deploymentCell) GetCreation() time.Time {
	return d.CreationTimestamp.Time
}
func (d deploymentCell) GetName() string {
	return d.Name
}
