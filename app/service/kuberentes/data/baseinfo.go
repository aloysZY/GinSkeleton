package data

import (
	service_data_interf "ginskeleton/app/service/kuberentes/data/data_interf"
	service_pod "ginskeleton/app/service/kuberentes/data/pod"
	"sort"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

// 实例化dataselector结构体,组装数据
func CreateDataFactory(dataCell []service_data_interf.DataCell, filterName string, limit, page int) *data {
	return &data{
		genericDataList: dataCell,
		dataQuery: &dataQuery{
			filter: &filterQuery{name: filterName},
			paginate: &paginateQuery{
				limit: limit,
				page:  page,
			},
		},
	}
}

type data struct {
	genericDataList []service_data_interf.DataCell //这个不能使用切片指针类型，不然后面使用不了方法
	dataQuery       *dataQuery
}

// dataSelectQuery 定义过滤和分页的结构体,过滤:name 分页:Limit和Page
type dataQuery struct {
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
func (d *data) Filter() *data {
	if d.dataQuery.filter.name == "" {
		return d
	}
	// 匹配 pod 名称，将符合名称的进行返回
	var filtered []service_data_interf.DataCell
	for _, value := range d.genericDataList {
		// 定义是否匹配标签变量,默认是匹配的
		// matches := true
		objName := value.GetName()
		// 如果 pod 名称中包含了要查找的名字，就添加到列表
		if strings.Contains(objName, d.dataQuery.filter.name) {
			filtered = append(filtered, value)
		}
	}
	d.genericDataList = filtered
	return d
}

// paginate 分页,根据Limit和Page的传参,取一定范围内的数据返回，调用这个函数之前要先排序
func (d *data) Paginate() *data {
	limit := d.dataQuery.paginate.limit
	page := d.dataQuery.paginate.page
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

//实现自定义的排序方法,需要重写Len,Swap,Less方法

// Len用于获取数组的长度
func (d *data) Len() int {
	return len(d.genericDataList)
}

// Swap用于数据比较大小后的位置变更
func (d *data) Swap(i, j int) {
	d.genericDataList[i], d.genericDataList[j] = d.genericDataList[j], d.genericDataList[i]
}

// Less用于比较大小,根据创建时间
func (d *data) Less(i, j int) bool {
	return d.genericDataList[i].GetCreation().Before(d.genericDataList[j].GetCreation())
}

// 重写以上三个方法,用sort.Sort 方法触发Len、swap、less
// 有一个问题，就是怎么确认是按照时间排序的？ less 根据 less的实现进行排序的（less 实现的是时间）
func (d *data) Sort() *data {
	sort.Sort(sort.Reverse(d)) // 反序排列
	return d
}

// 每种自定义都要实现一个将自己类型转换为data_interf.DataCell类型的方法
func (d *data) FromPod() []*corev1.Pod {
	pods := make([]*corev1.Pod, len(d.genericDataList)) //列表指针 make 的是列表，还需要额外初始化指针
	for i := range d.genericDataList {
		pods[i] = new(corev1.Pod)
		// cells[i].(podCell)是将DataCell类型转换成podCell
		*pods[i] = corev1.Pod(d.genericDataList[i].(service_pod.PodCell))
	}
	return pods
}
