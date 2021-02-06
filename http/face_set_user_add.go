package http

import (
	"github.com/brickman-source/golang-utilities/baidu"
	"github.com/brickman-source/golang-utilities/log"
	"github.com/brickman-source/golang-utilities/rand"
	"github.com/labstack/echo/v4"
)

type FaceSetUserAddResult struct {
	ErrCode    string                        `json:"err_code"`
	ErrMessage string                        `json:"err_message"`
	Biz        *baidu.FaceSetUserAddResponse `json:"biz"`
}

func (httpH *Handler) FaceSetUserAddHandler(ctx echo.Context) (err error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		log.Errorf("parse form: %v", err)
		return err
	}

	imageData,err := getFirstFileFromMultipartForm(form)
	if err != nil {
		log.Errorf("parse form: %v", err)
		return err
	}

	fullName := "unknown"
	for k, values := range form.Value {
		if k == "full_name" && len(values) > 0 {
			fullName = values[0]
		}
	}
	bdResult, err := httpH.ctx.Baidu.FaceSetUserAdd(
		httpH.ctx.Config.GetString("baidu.face.userGroupId"),
		"qc_"+rand.GetShortTimestampRandString(),
		fullName,
		imageData,
		httpH.ctx.Config.GetString("baidu.face.apiKey"),
		httpH.ctx.Config.GetString("baidu.face.secretKey"),
	)
	if err != nil {
		log.Errorf("face set add user err: %v", err)
		return ctx.JSON(200, &AuditResult{
			ErrCode:    "internal_server_error",
			ErrMessage: "处理异常",
		})
	}
	if bdResult.ErrorCode != 0 || (bdResult.ErrorMsg != "SUCCESS" && bdResult.ErrorMsg != "") {
		log.Infof("face set add user err: %v %v", bdResult.ErrorCode, bdResult.ErrorMsg)
		return ctx.JSON(200, &FaceSetUserAddResult{
			ErrCode:    "internal_server_error",
			ErrMessage: "处理异常",
		})
	}
	return ctx.JSON(200, &FaceSetUserAddResult{
		ErrCode:    "",
		ErrMessage: "",
		Biz:        bdResult,
	})
}
