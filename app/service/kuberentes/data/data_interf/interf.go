package data_interf

import "time"

// dataCell 接口,用于各种资源List的类型转换,转换后可以使用dataselector的排序,过滤,分页方法
type DataCell interface {
	GetCreation() time.Time
	GetName() string
}
