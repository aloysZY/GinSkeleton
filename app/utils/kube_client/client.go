package kube_client

import (
	"time"

	"k8s.io/client-go/dynamic"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func getConfig(cfgpath string) (*rest.Config, error) {
	if cfgpath != "" {
		kubeconfig, err := clientcmd.BuildConfigFromFlags("", cfgpath)
		if err != nil {
			return nil, err
		}
		return kubeconfig, nil
	}
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		// variable.ZapLog.Error("BuildConfigFromFlags kube Clientset err:", zap.Error(err))
		kubeconfig, err := rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
		return kubeconfig, nil
	}
	return kubeconfig, nil
}

func getClient(config *rest.Config) (*kubernetes.Clientset, *dynamic.DynamicClient, error) {
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		// variable.ZapLog.Error("NewForConfig err:%s", zap.Error(err))
		return nil, nil, err
	}
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		// variable.ZapLog.Error("NewForConfig err:%s", zap.Error(err))
		return nil, nil, err
	}
	return clientset, dynamicClient, nil
}

// 创建KController对象
func NewKubeControllerclient(cfgpath string, defaultResync time.Duration) (*KubeControllerClient, error) {

	config, err := getConfig(cfgpath)
	if err != nil {
		return nil, err
	}
	clientset, dynamicClient, err := getClient(config)
	if err != nil {
		return nil, err
	}

	kc := &KubeControllerClient{}
	kc.Clientset.Client = clientset
	kc.DynamicClient.Client = dynamicClient //暂时用不上informers这步先不管了

	//tianjia dyn
	// 通过Clientset生成SharedInformerFactory
	// https://blog.csdn.net/qq_17305249/article/details/107890063
	// https://github.com/cloudnativeto/sig-kubernetes/issues/11
	// Resync 机制的引入，定时将 Indexer 缓存事件重新同步到 Delta FIFO 队列中，在处理 SharedInformer 事件回调时，让处理失败的事件得到重新处理。并且通过入队前判断 FIFO 队列中是否已经有了更新版本的 event，来决定是否丢弃 Indexer 缓存不进行 Resync 入队。在处理 Delta FIFO 队列中的 Resync 的事件数据时，触发 onUpdate 回调来让事件重新处理。
	kc.factory = informers.NewSharedInformerFactory(clientset, defaultResync)
	// 生成Deployment、Pod、Service等资源对象的Informer、Lister以及HasSynced
	// ......
	kc.podInformer = kc.factory.Core().V1().Pods()
	kc.PodsLister = kc.podInformer.Lister()
	kc.podsSynced = kc.podInformer.Informer().HasSynced

	// ......
	return kc, nil
}
