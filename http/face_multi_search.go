package http

import (
	"github.com/brickman-source/golang-utilities/baidu"
	"github.com/brickman-source/golang-utilities/json"
	"github.com/brickman-source/golang-utilities/log"
	"github.com/labstack/echo/v4"
)

type FaceMultiSearchResult struct {
	ErrCode    string                         `json:"err_code"`
	ErrMessage string                         `json:"err_message"`
	Biz        *baidu.FaceMultiSearchResponse `json:"biz"`
}

func (httpH *Handler) FaceMultiSearchHandler(ctx echo.Context) (err error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		log.Errorf("parse form: %v", err)
		return err
	}

	imageData, err := getFirstFileFromMultipartForm(form)
	if err != nil {
		log.Errorf("parse form: %v", err)
		return err
	}
	qualityControl := "LOW"
	livenessControl := "NONE"
	for k, values := range form.Value {
		if k == "face_quality" && len(values) > 0 {
			qualityControl = values[0]
		} else if k == "face_liveness" && len(values) > 0 {
			livenessControl = values[0]
		}
	}

	bdResult, err := httpH.ctx.Baidu.FaceMultiSearch(
		[]string{httpH.ctx.Config.GetString("baidu.face.userGroupId")},
		imageData,
		baidu.FaceControlLevel(qualityControl),
		baidu.FaceControlLevel(livenessControl),
		httpH.ctx.Config.GetString("baidu.face.apiKey"),
		httpH.ctx.Config.GetString("baidu.face.secretKey"),
	)
	if err != nil {
		log.Errorf("face search err: %v", err)
		return ctx.JSON(200, &FaceMultiSearchResult{
			ErrCode:    "internal_server_error",
			ErrMessage: "处理异常",
		})
	}
	log.Infof("search result: %v", string(json.ShouldMarshal(bdResult)))
	if bdResult.ErrorCode == 222207 {
		// 未找到匹配的用户
		return ctx.JSON(200, &FaceMultiSearchResult{
			ErrCode:    "",
			ErrMessage: "",
			Biz:        bdResult,
		})
	}
	if bdResult.ErrorCode != 0 || (bdResult.ErrorMsg != "SUCCESS" && bdResult.ErrorMsg != "") {
		log.Infof("face search err: %v %v", bdResult.ErrorCode, bdResult.ErrorMsg)
		return ctx.JSON(200, &FaceMultiSearchResult{
			ErrCode:    "internal_server_error",
			ErrMessage: "处理异常",
		})
	}
	return ctx.JSON(200, &FaceMultiSearchResult{
		ErrCode:    "",
		ErrMessage: "",
		Biz:        bdResult,
	})
}
