package pods

import (
	"ginskeleton/app/global/consts"
	"ginskeleton/app/http/validator/core/data_transfer"
	"ginskeleton/app/utils/response"

	"github.com/gin-gonic/gin"
)

type PodUpdate struct {
	Namespace string `form:"namespace" json:"namespace" binding:"required"`
	PodName   string `form:"pod_name" json:"pod_name" binding:"required"`
	Content   string `form:"content" json:"content" binding:"required"`
	// image
}

func (p PodUpdate) CheckParams(context *gin.Context) {
	// 1.先按照验证器提供的基本语法，基本可以校验90%以上的不合格参数
	if err := context.ShouldBind(&p); err != nil {
		// 将表单参数验证器出现的错误直接交给错误翻译器统一处理即可
		response.ValidatorError(context, err)
		return
	}

	//  该函数主要是将绑定的数据以 键=>值 形式直接传递给下一步（控制器）
	extraAddBindDataContext := data_transfer.DataAddContext(p, consts.ValidatorPrefix, context)
	if extraAddBindDataContext == nil {
		response.ErrorSystem(context, "PodUpdate 表单验证器json化失败", "")
	} else {
		// 验证完成，调用控制器,并将验证器成员(字段)递给控制器，保持上下文数据一致性
		// 调用 controller 层
		//(&api.Pod{}).Update(extraAddBindDataContext)
	}
}
