package clientset

import (
	"context"
	model_clientset "ginskeleton/app/model/api/clientset"
)

// 构造函数
func CreateClientsetFactory(ctx context.Context) *Service {
	s := &Service{ctx: ctx}
	//其实这里返回的是所有的 client
	s.client = model_clientset.CreateClientsetFactory()
	return s
}

type Service struct {
	ctx    context.Context
	client *model_clientset.Clientset
}
