package api

import (
	"ginskeleton/app/global/variable"
	"ginskeleton/app/utils/kube_client"

	"k8s.io/apimachinery/pkg/labels"

	corev1 "k8s.io/api/core/v1"
)

func CreateClientsetFactory() *Clientset {
	return &Clientset{variable.KubeControllerClientset}
}

type Clientset struct {
	*kube_client.KubeControllerClientset
}

func (c *Clientset) List(namespace string) ([]*corev1.Pod, error) {
	// metav1.ListOptions{} 用于过滤List数据,如label,field等
	//return c.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})

	//labels.Everything() 所有标签
	podList, err := c.PodsLister.Pods(namespace).List(labels.Everything())
	if err != nil {
		return nil, err
	}
	return podList, nil
}
