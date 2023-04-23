package dynamicClient

import (
	"context"
	model_dynamicClient "ginskeleton/app/model/api/dynamicClient"
	"ginskeleton/app/utils/kube_client"
)

// 构造函数
func CreateDynamicClientFactory(ctx context.Context) *Service {
	s := &Service{ctx: ctx}
	s.client = model_dynamicClient.CreateDynamicClientFactory()
	return s
}

type Service struct {
	ctx    context.Context
	client *model_dynamicClient.DynamicClient
}

type DynamicClient struct {
	*kube_client.KubeControllerClient
}
