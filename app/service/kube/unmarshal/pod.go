package unmarshal

import (
	"k8s.io/apimachinery/pkg/util/yaml"

	corev1 "k8s.io/api/core/v1"
)

func UnmarshalPod(content string) (*corev1.Pod, error) {
	pod := &corev1.Pod{}
	if err := yaml.Unmarshal([]byte(content), pod); err != nil {
		return nil, err
	}
	return pod, nil
}
