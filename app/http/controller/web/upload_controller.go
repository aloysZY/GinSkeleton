package web

import (
	"ginskeleton/app/global/consts"
	"ginskeleton/app/global/variable"
	"ginskeleton/app/service/upload_file"
	"ginskeleton/app/utils/response"

	"github.com/gin-gonic/gin"
)

type Upload struct{}

//	文件上传是一个独立模块，给任何业务返回文件上传后的存储路径即可。
//
// 开始上传
func (u *Upload) StartUpload(ctx *gin.Context) {
	savePath := variable.BasePath + variable.ConfigYml.GetString("FileUploadSetting.UploadFileSavePath")
	if r, finnalSavePath := upload_file.Upload(ctx, savePath); r == true {
		response.Success(ctx, consts.CurdStatusOkMsg, finnalSavePath)
	} else {
		response.Fail(ctx, consts.FilesUploadFailCode, consts.FilesUploadFailMsg, "")
	}
}
