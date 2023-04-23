package kube_client

import (
	"fmt"

	util_runtime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/informers"
	informers_core_v1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	listers_core_v1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
)

// KController 对象
// 如果结构体内的字段，进行初始化的时候赋值，大小写都可以，但是额外进行赋值和使用，必须要大写开头
type KubeControllerClient struct {
	// kubeConfig *rest.Config
	Status int32
	// clusterId  []string
	// env        []string
	//初始化连接
	Clientset     Clientset
	DynamicClient DynamicClient
	factory       informers.SharedInformerFactory
	// 定义Deployment、Pod、Service等资源对象的Informer、Lister以及HasSynce
	// ......
	podInformer informers_core_v1.PodInformer
	PodsLister  listers_core_v1.PodLister
	podsSynced  cache.InformerSynced
	// ......
}

type Clientset struct {
	Client *kubernetes.Clientset //和DynamicClient 重名，是要显示调用
}

type DynamicClient struct {
	Client *dynamic.DynamicClient
}

// 启动Factory，获取缓存
func (kc *KubeControllerClient) Run(stopPodch chan struct{}) {
	// defer close(stopPodCh)
	defer util_runtime.HandleCrash()
	// defer variable.ZapLog.Error("KubeControllerClientset shutdown")
	// 传入停止的stopCh
	kc.factory.Start(stopPodch)
	// 等待资源查询（ClientsetList）完成后同步到缓存
	if !cache.WaitForCacheSync(stopPodch,
		// kc.nodesSynced, kc.deploymentsSynced,
		kc.podsSynced,
		// kc.ingressesSynced, kc.servicesSynced, kc.configMapsSynced, kc.namespaceSynced
	) {
		util_runtime.HandleError(fmt.Errorf("timed out waiting for kuberesource caches to sync"))
		return
	}
	// 同步成功，设置标志位status 为1
	kc.Status = 1
	// variable.ZapLog.Info("KubeControllerClientset start")
	<-stopPodch
}
