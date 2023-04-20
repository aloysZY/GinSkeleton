package register_validator

import (
	"ginskeleton/app/core/container"
	"ginskeleton/app/global/consts"
	"ginskeleton/app/http/validator/api/pods"
)

// 各个业务模块验证器必须进行注册（初始化），程序启动时会自动加载到容器
func ApiRegisterValidator() {
	// 创建容器
	containers := container.CreateContainersFactory()

	//  key 按照前缀+模块+验证动作 格式，将各个模块验证注册在容器
	var key string

	// pod list
	key = consts.ValidatorPrefix + "PodList"
	containers.Set(key, pods.PodList{})

	key = consts.ValidatorPrefix + "Detail"
	containers.Set(key, pods.PodDetail{})

	key = consts.ValidatorPrefix + "Delete"
	containers.Set(key, pods.PodDelete{})

	key = consts.ValidatorPrefix + "Update"
	containers.Set(key, pods.PodUpdate{})

	key = consts.ValidatorPrefix + "Create"
	containers.Set(key, pods.PodCreate{})
}
