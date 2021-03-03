package http

import (
	"github.com/brickman-source/golang-utilities/baidu"
	"github.com/brickman-source/golang-utilities/log"
	"github.com/labstack/echo/v4"
)

type PersonVerifyResult struct {
	ErrCode    string                      `json:"err_code"`
	ErrMessage string                      `json:"err_message"`
	Biz        *baidu.PersonVerifyResponse `json:"biz"`
}

func (httpH *Handler) PersonVerifyHandler(ctx echo.Context) (err error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		log.Errorf("parse form: %v", err)
		return err
	}

	//fullName, _ := getValueFromMultipartForm(form, "fullname")
	//idCardNo, _ := getValueFromMultipartForm(form, "id_card_no")

	idCardImageFrontData, err := getFileFromMultipartForm(form, "id_card_image_front")
	if err != nil {
		log.Errorf("parse form: %v", err)
		return ctx.JSON(200, &AuditResult{
			ErrCode:    "internal_server_error",
			ErrMessage: "处理异常",
		})
	}
	idCardRecognizeResp, err := httpH.ctx.Baidu.IdCardRecognize(
		idCardImageFrontData, baidu.IdCardSide_Front, false,
		httpH.ctx.Config.GetString("baidu.ocr.apiKey"),
		httpH.ctx.Config.GetString("baidu.ocr.secretKey"),
	)
	if err != nil {
		log.Errorf("recognize idcard err: %v", err)
		return ctx.JSON(200, &AuditResult{
			ErrCode:    "internal_server_error",
			ErrMessage: "处理异常",
		})
	}

	currentFaceImageData, err := getFileFromMultipartForm(form, "face_image")
	if err != nil {
		log.Errorf("parse form: %v", err)
		return err
	}

	qualityControl := "NORMAL"
	livenessControl := "NONE"
	for k, values := range form.Value {
		if k == "face_quality" && len(values) > 0 {
			qualityControl = values[0]
		} else if k == "face_liveness" && len(values) > 0 {
			livenessControl = values[0]
		}
	}

	bdResult, err := httpH.ctx.Baidu.PersonVerify(
		currentFaceImageData,
		idCardRecognizeResp.WordsResult.Name.Words,
		idCardRecognizeResp.WordsResult.IdCard.Words,
		baidu.FaceControlLevel(qualityControl),
		baidu.FaceControlLevel(livenessControl),
		httpH.ctx.Config.GetString("baidu.face.apiKey"),
		httpH.ctx.Config.GetString("baidu.face.secretKey"),
	)
	if err != nil {
		log.Errorf("recognize idcard err: %v", err)
		return ctx.JSON(200, &PersonVerifyResult{
			ErrCode:    "internal_server_error",
			ErrMessage: "处理异常",
		})
	}
	if bdResult.ErrorCode != 0 || (bdResult.ErrorMsg != "SUCCESS" && bdResult.ErrorMsg != "") {
		log.Infof("recognize idcard err: %v %v", bdResult.ErrorCode, bdResult.ErrorMsg)
		return ctx.JSON(200, &PersonVerifyResult{
			ErrCode:    "internal_server_error",
			ErrMessage: "处理异常",
		})
	}
	return ctx.JSON(200, &PersonVerifyResult{
		ErrCode:    "",
		ErrMessage: "",
		Biz:        bdResult,
	})
}
