/*
 * Copyright (c) 2020. Jinlong Chen.
 */

package http

import (
	"github.com/brickman-source/golang-utilities/http"
	"github.com/brickman-source/golang-utilities/log"
	"github.com/labstack/echo/v4"
	"strings"
)

func (httpH *Handler) AuditTextHandler(ctx echo.Context) (err error) {
	auditResult, err := httpH.ctx.Baidu.AuditText(
		string(http.GetRequestBody(ctx.Request())),
		httpH.ctx.Config.GetString("baidu.audit.apiKey"),
		httpH.ctx.Config.GetString("baidu.audit.secretKey"),
	)
	if err != nil {
		log.Errorf( "audit image err: %v", err)
		return ctx.JSON(200, &AuditResult{
			ErrCode:    "internal_server_error",
			ErrMessage: "文本处理异常",
			Success:    false,
		})
	}
	if auditResult.ErrorCode != 0 || auditResult.ErrorMsg != "" {
		log.Infof( "recognize picture err: %v %v", auditResult.ErrorCode, auditResult.ErrorMsg)
		return ctx.JSON(200, &AuditResult{
			ErrCode:    "internal_server_error",
			ErrMessage: "文本处理异常",
			Success:    false,
		})
	}
	if auditResult.ConclusionType != 1 {
		log.Infof( "recognize picture result: %v %v", auditResult.ConclusionType, auditResult.Conclusion)
		conclusionMessage := make([]string, 0)
		for _, datum := range auditResult.Data {
			conclusionMessage = append(conclusionMessage, datum.Msg)
		}
		return ctx.JSON(200, &AuditResult{
			ErrCode:    "text_error",
			ErrMessage: strings.Join(conclusionMessage, ","),
			Success:    false,
		})
	}
	return ctx.JSON(200, &AuditResult{
		ErrCode:    "",
		ErrMessage: "",
		Success:    true,
	})
}
