package response

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type PodResp struct {
	Total int        `json:"total"`
	Items []*Message `json:"items"`
}

// 封装一个二级目录列表,返回信息
type Message struct {
	Name         string              `json:"name,omitempty"` //,omitempty若为空值，则字符串中不会包含它
	Namespace    string              `json:"namespace,omitempty"`
	Replicas     uint8               `json:"replicas,omitempty"`
	Status       string              `json:"status,omitempty"`
	CreationTime metav1.Time         `json:"creation_Time,omitempty"`
	Containers   []map[string]string `json:"containers,omitempty"`
}
