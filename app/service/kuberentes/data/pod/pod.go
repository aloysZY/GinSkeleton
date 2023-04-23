package pod

import (
	"ginskeleton/app/service/kuberentes/data/data_interf"
	"time"

	corev1 "k8s.io/api/core/v1"
)

// 定义podCell, 重写GetCreation和GetName 方法后,可以进行数据转换
// covev1.Pod --> PodCell  --> dataCell
// appsv1.Deployment --> deployCell --> dataCell
// 这里为什么要转化类型和实现一些方法（这些方法 pod 存在）返回的数据都转化为一个类型DataCell，在根据需要，转化为不同的类型
type PodCell corev1.Pod

// 重写DataCell接口的两个方法
func (p PodCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}
func (p PodCell) GetName() string {
	return p.Name
}

// 每种自定义都要实现一个将自己类型转换为data_interf.DataCell类型的方法
// 类型转换方法corev1.Pod --> DataCell,DataCell-->corev1.Pod
func PodToCells(pods []*corev1.Pod) []data_interf.DataCell {
	cells := make([]data_interf.DataCell, len(pods))
	for i := range pods {
		cells[i] = PodCell(*pods[i])
	}
	return cells
}
