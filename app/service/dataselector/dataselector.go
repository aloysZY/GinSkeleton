package dataselector

import (
	"sort"
	"strings"

	"ginskeleton/app/service/interf"
	"ginskeleton/app/service/kube"

	corev1 "k8s.io/api/core/v1"
)

// 实例化dataselector结构体,组装数据
func CreateDataSelectorFactory(dataCell []interf.DataCell, filterName string, limit, page int) *dataSelector {
	return &dataSelector{
		genericDataList: dataCell,
		dataSelect: &dataSelectQuery{
			filter: &filterQuery{name: filterName},
			paginate: &paginateQuery{
				limit: limit,
				page:  page,
			},
		},
	}
}

type dataSelector struct {
	genericDataList []interf.DataCell
	dataSelect      *dataSelectQuery
}

// dataSelectQuery 定义过滤和分页的结构体,过滤:name 分页:Limit和Page
type dataSelectQuery struct {
	filter   *filterQuery
	paginate *paginateQuery
}

// filterQuery 用于查询 过滤:name
type filterQuery struct {
	name string
}

// 分页:Limit和Page Limit是单页的数据条数,Page是第几页
type paginateQuery struct {
	page  int
	limit int
}

// Filter方法用于过滤,比较数据Name属性,若包含则返回
func (d *dataSelector) filter() *dataSelector {
	if d.dataSelect.filter.name == "" {
		return d
	}
	// 匹配 pod 名称，将符合名称的进行返回
	var filtered []interf.DataCell
	for _, value := range d.genericDataList {
		// 定义是否匹配标签变量,默认是匹配的
		// matches := true
		objName := value.GetName()
		// 如果 pod 名称中包含了要查找的名字，就添加到列表
		// if !strings.Contains(objName, d.dataSelect.filter.name) {
		// 	matches = false
		// 	continue // 跳出这个 if 判断，继续执行代码
		// }
		// if matches {
		// 	filtered = append(filtered, value)
		// }
		// 如果 pod 名称中包含了要查找的名字，就添加到列表
		if strings.Contains(objName, d.dataSelect.filter.name) {
			filtered = append(filtered, value)
		}
	}
	d.genericDataList = filtered
	return d
}

// paginate 分页,根据Limit和Page的传参,取一定范围内的数据返回，调用这个函数之前要先排序
func (d *dataSelector) paginate() *dataSelector {
	limit := d.dataSelect.paginate.limit
	page := d.dataSelect.paginate.page
	// 验证参数合法，若参数不合法，则返回所有数据
	if limit <= 0 || page <= 0 {
		return d
	}
	// 举例：25个元素的数组，limit是10(每页显示 10 个数据)，page是3，startIndex是20（第三页的第一个元素下标），endIndex是30（实际上endIndex是25）（第三页最后一个元素下标）
	startIndex := limit * (page - 1)
	endIndex := limit * page

	// 处理最后一页，这时候就把endIndex由30改为25了
	if len(d.genericDataList) < endIndex {
		endIndex = len(d.genericDataList)
	}
	d.genericDataList = d.genericDataList[startIndex:endIndex] // 切片分割
	return d
}

// func (d *dataSelector) FromCells(cells []interf.DataCell) []corev1.Pod {
// 	pods := make([]corev1.Pod, len(cells))
// 	for i := range cells {
// 		// cells[i].(podCell)是将DataCell类型转换成podCell
// 		pods[i] = corev1.Pod(cells[i].(kube.PodCell))
// 	}
// 	return pods
// }

func (d *dataSelector) fromCells() []corev1.Pod {
	pods := make([]corev1.Pod, len(d.genericDataList))
	for i := range d.genericDataList {
		pods[i] = corev1.Pod(d.genericDataList[i].(kube.PodCell))
	}
	return pods
}

/*实现自定义的排序方法,需要重写Len,Swap,Less方法
这个排序没有什么具体意义，因为申请时间排序才有一点参考价值，应该在第一次创建的时候写入时间*/

// Len用于获取数组的长度
func (d *dataSelector) Len() int {
	return len(d.genericDataList)
}

// Swap用于数据比较大小后的位置变更
func (d *dataSelector) Swap(i, j int) {
	d.genericDataList[i], d.genericDataList[j] = d.genericDataList[j], d.genericDataList[i]
}

// Less用于比较大小,根据创建时间
func (d *dataSelector) Less(i, j int) bool {
	return d.genericDataList[i].GetCreation().Before(d.genericDataList[j].GetCreation())
}

// 重写以上三个方法,用sort.Sort 方法触发排序
func (d *dataSelector) Sort() *dataSelector {
	sort.Sort(sort.Reverse(d)) // 返序排列
	return d
}

func (d *dataSelector) PodList() []kube.PodList {
	data := make([]kube.PodList, d.Len())
	for k, pod := range d.filter().Sort().paginate().fromCells() { // 过滤对应的 pod名称 排序、分页、类型转化
		data[k].Name = pod.Name
		data[k].Namespace = pod.Namespace
		data[k].Status = string(pod.Status.Phase)
		for _, container := range pod.Spec.Containers {
			// yyyMap :=new()
			// yyy := make([]map[string]string, len(pod.Spec.Containers)) //这样写就多了一层列表，所以不需要了
			containerMsg := make(map[string]string) // 对 []map[string]string 中的 map 进行初始化
			containerMsg["name"] = container.Name
			containerMsg["Image"] = container.Image
			data[k].Containers = append(data[k].Containers, containerMsg)
		}
	}
	return data
}
