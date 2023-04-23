package clientset

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (c *Clientset) ListPod(namespace string) ([]*corev1.Pod, error) {
	// metav1.ListOptions{} 用于过滤List数据,如label,field等
	//return c.CoreV1().Pods(namespace).ClientsetList(context.TODO(), metav1.ListOptions{})
	//labels.Everything() 所有标签
	podList, err := c.kc.PodsLister.Pods(namespace).List(labels.Everything())
	if err != nil {
		return nil, err
	}
	return podList, nil
}

func (c *Clientset) DetailPod(namespace, name string) (*corev1.Pod, error) {
	pod, err := c.kc.PodsLister.Pods(namespace).Get(name)
	if err != nil {
		return nil, err
	}
	return pod, nil
}
