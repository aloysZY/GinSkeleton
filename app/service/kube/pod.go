package kube

import (
	"context"
	"ginskeleton/app/global/variable"
	"ginskeleton/app/model/api"

	corev1 "k8s.io/api/core/v1"
)

// 构造函数
func CreatePodFactory(ctx context.Context) *pod {
	p := &pod{ctx: ctx}
	p.clientset = api.CreateClientsetFactory()
	return p

}

// 定义一个空结构体，主要是为了实现一些方法
type pod struct {
	ctx       context.Context
	clientset *api.Clientset
}

type PodResp struct {
	Total int       `json:"total"`
	Items []PodList `json:"items"`
}

// 封装一个二级目录列表,返回信息
type PodList struct {
	Name       string              `json:"name"`
	Namespace  string              `json:"namespace"`
	Status     string              `json:"status"`
	Containers []map[string]string `json:"containers"`
}

// 获取pod列表,支持过滤,排序,分页,模糊查找
func (p *pod) List(namespace string) ([]*corev1.Pod, error) {
	// 调用 model 层查询 pod
	podList, err := p.clientset.List(namespace)
	//podList, err := api_model.CreateClientsetFactory().List(namespace)
	if err != nil {
		variable.ZapLog.Sugar().Info("List pod failed error: %v\n", err)
		return nil, err
	}
	return podList, nil
}

/*// 获取单个pod
func (p *pod) Detail(namespace, podName string) (*corev1.Pod, error) {
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
