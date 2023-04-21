package kube_client

import (
	"fmt"
	"time"

	informers_core_v1 "k8s.io/client-go/informers/core/v1"

	util_runtime "k8s.io/apimachinery/pkg/util/runtime"

	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	listers_core_v1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

// 获取配置文件
func getKubeClientset(cfgpath string) (*kubernetes.Clientset, error) {
	var kubeconfig *rest.Config
	var err error
	if cfgpath == "" {
		kubeconfig, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
		if err != nil {
			// variable.ZapLog.Error("BuildConfigFromFlags kube clientset err:", zap.Error(err))
			kubeconfig, err = rest.InClusterConfig()
			if err != nil {
				return nil, err
			}
		}
	} else {
		kubeconfig, err = clientcmd.BuildConfigFromFlags("", cfgpath)
		if err != nil {
			return nil, err
		}
	}

	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		// variable.ZapLog.Error("NewForConfig err:%s", zap.Error(err))
		return nil, err
	}
	// variable.ZapLog.Info("getKubeClientset OK")
	return clientset, nil
}

// 创建KController对象
func NewKubeControllerclientset(cfgpath string, defaultResync time.Duration) (*KubeControllerClientset, error) {
	clientset, err := getKubeClientset(cfgpath)
	if err != nil {
		return nil, err
	}

	kc := &KubeControllerClientset{Clientset: clientset}
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

// KController 对象
type KubeControllerClientset struct {
	// kubeConfig *rest.Config
	Status int32
	// clusterId  []string
	// env        []string
	Clientset *kubernetes.Clientset
	factory   informers.SharedInformerFactory
	// 定义Deployment、Pod、Service等资源对象的Informer、Lister以及HasSynce
	// ......
	podInformer informers_core_v1.PodInformer
	PodsLister  listers_core_v1.PodLister
	podsSynced  cache.InformerSynced
	// ......
}

// 启动Factory，获取缓存
func (kc *KubeControllerClientset) Run(stopPodch chan struct{}) {
	// defer close(stopPodCh)
	defer util_runtime.HandleCrash()
	// defer variable.ZapLog.Error("KubeControllerClientset shutdown")
	// 传入停止的stopCh
	kc.factory.Start(stopPodch)
	// 等待资源查询（List）完成后同步到缓存
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
