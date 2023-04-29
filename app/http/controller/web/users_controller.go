package web

import (
	"time"

	"ginskeleton/app/global/consts"
	"ginskeleton/app/global/variable"
	"ginskeleton/app/model/web/user"
	"ginskeleton/app/service/users/curd"
	userstoken "ginskeleton/app/service/users/token"
	"ginskeleton/app/utils/response"

	"github.com/gin-gonic/gin"
)

type Users struct{}

// 1.用户注册
func (u *Users) Register(ctx *gin.Context) {
	//  由于本项目骨架已经将表单验证器的字段(成员)绑定在上下文，因此可以按照 GetString()、context.GetBool()、GetFloat64（）等快捷获取需要的数据类型，注意：相关键名规则：  前缀+验证器结构体中的 json 标签
	// 注意：在 ginskeleton 中获取表单参数验证器中的数字键（字段）,请统一使用 GetFloat64(),其它获取数字键（字段）的函数无效，例如：GetInt()、GetInt64()等(数字类型在 gin 默认格式化的时候就是 float64)
	// 当然也可以通过gin框架的上下文原始方法获取，例如： context.PostForm("user_name") 获取，这样获取的数据格式为文本，需要自己继续转换
	userName := ctx.GetString(consts.ValidatorPrefix + "user_name")
	pass := ctx.GetString(consts.ValidatorPrefix + "pass")
	userIp := ctx.ClientIP()
	if curd.CreateUserCurdFactory(ctx.Request.Context()).Register(userName, pass, userIp) {
		response.Success(ctx, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(ctx, consts.CurdRegisterFailCode, consts.CurdRegisterFailMsg, "")
	}
}

// 2.用户登录
func (u *Users) Login(ctx *gin.Context) {
	userName := ctx.GetString(consts.ValidatorPrefix + "user_name")
	pass := ctx.GetString(consts.ValidatorPrefix + "pass")
	phone := ctx.GetString(consts.ValidatorPrefix + "phone")

	userModelFact := user.CreateUserFactory(ctx.Request.Context(), "")
	userModel := userModelFact.Login(ctx.Request.Context(), userName, pass)

	if userModel != nil {
		userTokenFactory := userstoken.CreateUserFactory()
		if userToken, err := userTokenFactory.GenerateToken(userModel.Id, userModel.UserName, userModel.Phone, variable.ConfigYml.GetInt64("Token.JwtTokenCreatedExpireAt")); err == nil {
			if userTokenFactory.RecordLoginToken(ctx.Request.Context(), userToken, ctx.ClientIP()) {
				data := gin.H{
					"userId":     userModel.Id,
					"user_name":  userName,
					"realName":   userModel.RealName,
					"phone":      phone,
					"token":      userToken,
					"updated_at": time.Now().Format(variable.DateFormat),
				}
				response.Success(ctx, consts.CurdStatusOkMsg, data)
				go userModel.UpdateUserloginInfo(ctx.Request.Context(), ctx.ClientIP(), userModel.Id)
				return
			}
		}
	}
	response.Fail(ctx, consts.CurdLoginFailCode, consts.CurdLoginFailMsg, "")
}

// 刷新用户token
func (u *Users) RefreshToken(ctx *gin.Context) {
	oldToken := ctx.GetString(consts.ValidatorPrefix + "token")
	if newToken, ok := userstoken.CreateUserFactory().RefreshToken(ctx.Request.Context(), oldToken, ctx.ClientIP()); ok {
		res := gin.H{
			"token": newToken,
		}
		response.Success(ctx, consts.CurdStatusOkMsg, res)
	} else {
		response.Fail(ctx, consts.CurdRefreshTokenFailCode, consts.CurdRefreshTokenFailMsg, "")
	}
}

// 后面是 curd 部分，自带版本中为了降低初学者学习难度，使用了最简单的方式操作 增、删、改、查
// 在开发企业实际项目中，建议使用我们提供的一整套 curd 快速操作模式
// 参考地址：https://gitee.com/daitougege/GinSkeleton/blob/master/docs/concise.md
// 您也可以参考 Admin 项目地址：https://gitee.com/daitougege/gin-skeleton-admin-backend/ 中， app/model/  提供的示例语法

// 3.用户查询（show）
func (u *Users) Show(ctx *gin.Context) {
	userName := ctx.GetString(consts.ValidatorPrefix + "user_name")
	page := ctx.GetFloat64(consts.ValidatorPrefix + "page")
	limit := ctx.GetFloat64(consts.ValidatorPrefix + "limit")
	limitStart := (page - 1) * limit
	counts, showlist := user.CreateUserFactory(ctx.Request.Context(), "").Show(userName, int(limitStart), int(limit))
	if counts > 0 && showlist != nil {
		response.Success(ctx, consts.CurdStatusOkMsg, gin.H{"counts": counts, "list": showlist})
	} else {
		response.Fail(ctx, consts.CurdSelectFailCode, consts.CurdSelectFailMsg, "")
	}
}

// 4.用户新增(store)
func (u *Users) Store(ctx *gin.Context) {
	userName := ctx.GetString(consts.ValidatorPrefix + "user_name")
	pass := ctx.GetString(consts.ValidatorPrefix + "pass")
	realName := ctx.GetString(consts.ValidatorPrefix + "real_name")
	phone := ctx.GetString(consts.ValidatorPrefix + "phone")
	remark := ctx.GetString(consts.ValidatorPrefix + "remark")

	if curd.CreateUserCurdFactory(ctx.Request.Context()).Store(userName, pass, realName, phone, remark) {
		response.Success(ctx, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(ctx, consts.CurdCreatFailCode, consts.CurdCreatFailMsg, "")
	}
}

// 5.用户更新(update)
func (u *Users) Update(ctx *gin.Context) {
	// 表单参数验证中的int、int16、int32 、int64、float32、float64等数字键（字段），请统一使用 GetFloat64() 获取，其他函数无效
	userId := ctx.GetFloat64(consts.ValidatorPrefix + "id")
	userName := ctx.GetString(consts.ValidatorPrefix + "user_name")
	pass := ctx.GetString(consts.ValidatorPrefix + "pass")
	realName := ctx.GetString(consts.ValidatorPrefix + "real_name")
	phone := ctx.GetString(consts.ValidatorPrefix + "phone")
	remark := ctx.GetString(consts.ValidatorPrefix + "remark")
	userIp := ctx.ClientIP()

	// 检查正在修改的用户名是否被其他人使用
	if user.CreateUserFactory(ctx.Request.Context(), "").UpdateDataCheckUserNameIsUsed(int(userId), userName) > 0 {
		response.Fail(ctx, consts.CurdUpdateFailCode, consts.CurdUpdateFailMsg+", "+userName+" 已经被其他人使用", "")
		return
	}

	// 注意：这里没有实现更加精细的权限控制逻辑，例如：超级管理管理员可以更新全部用户数据，普通用户只能修改自己的数据。目前只是验证了token有效、合法之后就可以进行后续操作
	// 实际使用请根据真是业务实现权限控制逻辑、再进行数据库操作
	if curd.CreateUserCurdFactory(ctx.Request.Context()).Update(ctx.Request.Context(), int(userId), userName, pass, realName, phone, remark, userIp) {
		response.Success(ctx, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(ctx, consts.CurdUpdateFailCode, consts.CurdUpdateFailMsg, "")
	}

}

// 6.删除记录
func (u *Users) Destroy(ctx *gin.Context) {
	// 表单参数验证中的int、int16、int32 、int64、float32、float64等数字键（字段），请统一使用 GetFloat64() 获取，其他函数无效
	userId := ctx.GetFloat64(consts.ValidatorPrefix + "id")
	if user.CreateUserFactory(ctx.Request.Context(), "").Destroy(ctx.Request.Context(), int(userId)) {
		response.Success(ctx, consts.CurdStatusOkMsg, "")
	} else {
		response.Fail(ctx, consts.CurdDeleteFailCode, consts.CurdDeleteFailMsg, "")
	}
}
