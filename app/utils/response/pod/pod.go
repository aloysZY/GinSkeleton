package pod

import (
	"ginskeleton/app/utils/response"

	corev1 "k8s.io/api/core/v1"
)

// 返回信息的列表
func Pod(podList []*corev1.Pod) []*response.Message {
	// 这里初始化的容量，应该是要返回到前端，已经处理过的数据的长度
	data := make([]*response.Message, len(podList))
	for k, pod := range podList { // 过滤对应的 pod名称 排序、分页、类型转化
		//data[k] = new(response.Message) //初始化在赋值
		//data[k].Name = pod.Name
		//data[k].Namespace = pod.Namespace
		//data[k].Status = string(pod.Status.Phase)
		//data[k].CreationTime = pod.CreationTimestamp
		data[k] = &response.Message{ //初始化赋值
			Name:         pod.Name,
			Namespace:    pod.Namespace,
			Status:       string(pod.Status.Phase),
			CreationTime: pod.CreationTimestamp,
		} //这两种有什么区别？
		for _, container := range pod.Spec.Containers {
			containerMsg := make(map[string]string) // 对 []map[string]string 中的 map 进行初始化
			containerMsg["name"] = container.Name
			containerMsg["Image"] = container.Image
			data[k].Containers = append(data[k].Containers, containerMsg)
		} //这里的循环，当Containers元素特别多的时候，最好使用的是下标进行循环
		//下标循坏
		//for i, tl := 0, len(pod.Spec.Containers); i < tl; i++ {
		//	containerMsg := make(map[string]string)
		//	containerMsg["name"] = pod.Spec.Containers[i].Name
		//	containerMsg["Image"] = pod.Spec.Containers[i].Image
		//	data[k].Containers = append(data[k].Containers, containerMsg)
		//}
	}
	return data
}
