package dynamicClient

import (
	"ginskeleton/app/utils/kube_client"
)

func CreateDynamicClientFactory() *DynamicClient {
	return &DynamicClient{}
}

type DynamicClient struct {
	*kube_client.KubeControllerClient
}
