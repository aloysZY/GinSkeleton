package api

import (
	"ginskeleton/app/global/consts"
	"ginskeleton/app/service/dataselector"
	"ginskeleton/app/service/kube"
	"ginskeleton/app/utils/response"

	"github.com/gin-gonic/gin"
)

type Pod struct{}

func (p *Pod) List(context *gin.Context) {

	//  由于本项目骨架已经将表单验证器的字段(成员)绑定在上下文，因此可以按照 GetString()、GetInt64()、GetFloat64（）等快捷获取需要的数据类型 ,因为在进行封装的时候上下文写入到空接口上，gin 默认解析为GetString()、GetInt64()、GetFloat64（）
	// 当然也可以通过gin框架的上下文原原始方法获取，例如： context.PostForm("name") 获取，这样获取的数据格式为文本，需要自己继续转换
	filterName := context.GetString(consts.ValidatorPrefix + "filter_name")
	namespace := context.GetString(consts.ValidatorPrefix + "namespace")
	limit := context.GetFloat64(consts.ValidatorPrefix + "limit")
	page := context.GetFloat64(consts.ValidatorPrefix + "page")

	//根据命名空间查询 pod
	dataCell, err := kube.CreatePodFactory().List(namespace)
	if err != nil {
		response.ErrorSystem(context, "获取 pod 列表失败", nil) // 返回前端，还是错误封装一下
		return
	}

	//PodList 进行类型转换，排序和根据 pood 名称过滤
	data := dataselector.CreateDataSelectorFactory(dataCell, filterName, int(limit), int(page)).PodList()

	// // 排序和分页,转化为 pod 类型
	// data := dataSelector.Sort().paginate().FromCells()
	//
	// // 修改为返回特定字段
	// /*	datalen := len(data)
	// 	dataMapList := make([]map[string]string, datalen)
	// 	// 添加特定字段
	// 	// 这个地方统筹规划一下吧，全都是切片了，太乱了,直接使用一个结构体处理吧
	// 	for k, v := range data {
	// 		dataMap := make([]map[string]string, datalen)
	// 		dataMap[k]["pod_name"] = v.Name
	// 		for x, y := range v.Spec.Containers {
	// 			dataC := make([]map[string]string, len(v.Spec.Containers))
	// 			dataC[x]["containers_name"] = y.Name
	// 			dataC[x]["containers_image"] = y.Image
	//
	// 		}
	// 		dataMap["status"] = string(v.Status.Phase)
	// 		dataMapList[k] = dataMap
	// 	}*/
	// data1 := make([]kube.PodList, len(data))
	// for k, pod := range data {
	// 	data1[k].Name = pod.Name
	// 	data1[k].Status = string(pod.Status.Phase)
	// 	for _, container := range pod.Spec.Containers {
	// 		// yyyMap :=new()
	// 		// yyy := make([]map[string]string, len(pod.Spec.Containers)) //这样写就多了一层列表，所以不需要了
	// 		containerMsg := make(map[string]string) // 对 []map[string]string 中的 map 进行初始化
	// 		containerMsg["name"] = container.Name
	// 		containerMsg["Image"] = container.Image
	// 		data1[k].Containers = append(data1[k].Containers, containerMsg)
	// 	}
	// }

	response.Success(context, "podList", &kube.PodResp{
		Total: len(dataCell), // 返回的 pod 总数量，data 数据的长度不准了，应为处理了数据
		Items: data,
	})
}

//func (p *Pod) Detail(context *gin.Context) {
//	podName := context.GetString(consts.ValidatorPrefix + "pod_name")
//	namespace := context.GetString(consts.ValidatorPrefix + "namespace")
//
//	// 获取相关命名空间下所有 pod，namespace 不能为空
//	pod, err := kube.CreatePodFactory().Detail(namespace, podName)
//	if err != nil {
//		if strings.HasSuffix(err.Error(), "not found") {
//			response.Success(context, err.Error(), nil)
//			return
//		}
//		response.ErrorSystem(context, "获取 pod 详情失败", nil) // 返回前端，还是错误封装一下
//		return
//	}
//	response.Success(context, "Detail", pod)
//}
//
//func (p *Pod) Delete(context *gin.Context) {
//	podName := context.GetString(consts.ValidatorPrefix + "pod_name")
//	namespace := context.GetString(consts.ValidatorPrefix + "namespace")
//
//	// 获取相关命名空间下所有 pod，namespace 不能为空
//	err := kube.CreatePodFactory().Delete(namespace, podName)
//	if err != nil {
//		if strings.HasSuffix(err.Error(), "not found") {
//			response.Success(context, err.Error(), nil)
//			return
//		}
//		response.ErrorSystem(context, "删除 pod 失败", nil) // 返回前端，还是错误封装一下
//		return
//	}
//	response.Success(context, "Delete", "删除成功")
//}
//
//func (p *Pod) Update(context *gin.Context) {
//	namespace := context.GetString(consts.ValidatorPrefix + "namespace")
//	podName := context.GetString(consts.ValidatorPrefix + "pod_name")
//	content := context.GetString(consts.ValidatorPrefix + "content")
//	// content, ok := context.Get(consts.ValidatorPrefix + "content")
//	// if !ok {
//	// 	response.Success(context, "missing content", "")
//	// 	return
//	// }
//	if err := kube.CreatePodFactory().Update(namespace, podName, content); err != nil {
//		response.Success(context, "更新 pod 成功", "")
//		return
//	}
//	response.ErrorSystem(context, "更新 pod 失败", "")
//}
//
//func (p *Pod) Create(context *gin.Context) {
//
//	// namespace := context.GetString(consts.ValidatorPrefix + "namespace")
//	// podName := context.GetString(consts.ValidatorPrefix + "pod_name")
//	content := context.GetString(consts.ValidatorPrefix + "content")
//	pod, err := unmarshal.UnmarshalPod(content)
//	if err != nil {
//		variable.ZapLog.Error(" Create pod yaml.Unmarshal failed: ", zap.Error(err))
//		response.ErrorSystem(context, "json.Unmarshal pod 失败", "")
//		return
//	}
//	if pod.Namespace == "" {
//		// pod.Namespace = "default"
//		pod.Namespace = corev1.NamespaceDefault
//	}
//
//	// 1.判断 pod 是否存在
//	_, err = kube.CreatePodFactory().Detail(pod.Namespace, pod.Name)
//	if err != nil && strings.HasSuffix(err.Error(), "not found") {
//		// 2.不存在创建
//		if err := kube.CreatePodFactory().Create(pod.Namespace, pod); err != nil {
//			variable.ZapLog.Info("Create pod error :", zap.Error(err))
//			response.ErrorSystem(context, "Create pod error"+err.Error(), "")
//			return
//		}
//		response.Success(context, "Create pod success!", "")
//		return
//	}
//	response.Success(context, "pod already exists "+pod.Name, "")
//}
