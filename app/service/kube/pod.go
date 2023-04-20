package kube

import (
	"time"

	"ginskeleton/app/global/variable"
	"ginskeleton/app/model/api_model"
	"ginskeleton/app/service/interf"

	corev1 "k8s.io/api/core/v1"
)

// 构造函数
func CreatePodFactory() *pod {
	return &pod{}
}

// 定义一个空结构体，主要是为了实现一些方法
type pod struct{}

// 定义一个新的类型，基于corev1.pod类型
type PodCell corev1.Pod

// 针对自定义类型实现方法，实现DataCell接口  排序的时候需要用到

// 重写DataCell接口的两个方法
func (p PodCell) GetCreation() time.Time {
	return p.CreationTimestamp.Time
}

func (p PodCell) GetName() string {
	return p.Name
}

// 定义列表的返回内容,Items是pod元素列表,Total是元素数量，返回前端的结构体，根据前端需求返回（这里是直接全部返回 pod 字段）
//
//	type PodResp struct {
//		Total int          `json:"total"`
//		Items []corev1.Pod `json:"items"`
//	}
type PodResp struct {
	Total int       `json:"total"`
	Items []PodList `json:"items"`
}

// 封装一个二级目录列表
type PodList struct {
	Name       string              `json:"name"`
	Namespace  string              `json:"namespace"`
	Status     string              `json:"status"`
	Containers []map[string]string `json:"containers"`
}

// 类型转换方法corev1.Pod --> DataCell,DataCell-->corev1.Pod
func (p *pod) toCells(pods []*corev1.Pod) []interf.DataCell {
	cells := make([]interf.DataCell, len(pods))
	for i := range pods {
		cells[i] = PodCell(*pods[i])
	}
	return cells
}

// 获取pod列表,支持过滤,排序,分页,模糊查找
func (p *pod) List(namespace string) ([]interf.DataCell, error) {
	// context.TODO()  用于声明一个空的context上下文,用于List方法内设置这个请求超时
	//调用 model 层查询 pod
	podList, err := api_model.CreateClientsetFactory().List(namespace)
	if err != nil {
		variable.ZapLog.Sugar().Info("List pod failed error: %v\n", err)
		return nil, err
	}
	/*	podList, err := variable.Clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			variable.ZapLog.Info("获取 pod 列表失败", zap.Error(err))
			// 返回给上一层,最终返回给前端,前端捕获到后打印出来
			return nil, errors.New(fmt.Sprintf("获取 pod 列表失败, namespace: %s ,pod: %s", namespace, err))
		}*/
	return p.toCells(podList), nil
}

/*// 获取单个pod
func (p *pod) Detail(namespace, podName string) (*corev1.Pod, error) {
	// context.TODO()  用于声明一个空的context上下文,用于List方法内设置这个请求超时
	// metav1.ListOptions{} 用于过滤List数据,如label,field等
	// 获取单独的 pod，要指定命名空间，不能为空
	pod, err := variable.Clientset.CoreV1().Pods(namespace).Get(context.TODO(), podName, metav1.GetOptions{})
	if err != nil {
		variable.ZapLog.Info("获取 pod 详情失败", zap.Error(err))
		// variable.ZapLog.Info("获取 pod 详情失败->", zap.String("namespace", namespace), zap.Any("pod", err))
		// 返回给上一层,最终返回给前端,前端捕获到后打印出来
		return nil, errors.New(fmt.Sprintf("获取 pod 详情失败, namespace: %s, pod: %s", namespace, err))
	}
	return pod, nil
}

func (p *pod) Delete(namespace, podName string) error {
	err := variable.Clientset.CoreV1().Pods(namespace).Delete(context.TODO(), podName, metav1.DeleteOptions{})
	if err != nil {
		variable.ZapLog.Info("删除 pod 失败", zap.Error(err))
		// variable.ZapLog.Info("获取 pod 详情失败->", zap.String("namespace", namespace), zap.Any("pod", err))
		// 返回给上一层,最终返回给前端,前端捕获到后打印出来
		return errors.New(fmt.Sprintf("删除 pod 失败, namespace: %s, pod: %s", namespace, err))
	}
	return nil
}

func (p *pod) Update(namespace, podName, content string) error {
	// 总体步骤应该是
	// 1.查询 pod 是否存在，存在进行更新，不存在进行创建，保存查询到的 pod 句柄
	// 2.更新，或创建操作都是让用户输入的 yaml 进行反序列化到查询到的 pod 中

	// 这里的代码还是存在一些问题，要仔细的处理一下要修改那些字段

	// 想法是content不是一个整个的 pod，是一个 key,value，替换相关
	// contentStr, _ := content.(string)
	// pod := &corev1.Pod{} // 感觉这样传入的时候太麻烦了，需要传入所有pod 信息
	// if err := json.Unmarshal([]byte(content), pod); err != nil {
	// 	variable.ZapLog.Sugar().Info("json.Unmarshal pod %v failed: %v", podName, err)
	// 	return err
	// }
	// 尝试根据 pod 信息获取 pod，然后将content反序列化到新 pod 中（修改对于的配置）
	// detail, err := p.Detail(namespace, podName)
	// if err != nil {
	// 	return err
	// }
	// if err := json.Unmarshal([]byte(content), detail); err != nil {
	// 	variable.ZapLog.Sugar().Info("json.Unmarshal pod %v failed: %v", podName, err)
	// 	return err
	// }
	//
	// _, err = variable.Clientset.CoreV1().Pods(namespace).Update(context.TODO(), detail, metav1.UpdateOptions{})
	// if err != nil {
	// 	variable.ZapLog.Sugar().Info("Update pod %v failed: %v", podName, err)
	// 	return err
	// }
	return nil
}

func (p *pod) Create(namespace string, pod *corev1.Pod) error {
	_, err := variable.Clientset.CoreV1().Pods(namespace).Create(context.TODO(), pod, metav1.CreateOptions{})
	if err != nil {
		// variable.ZapLog.Info("Create pod err :", zap.Error(err))
		return err
	}
	return nil
}
*/
