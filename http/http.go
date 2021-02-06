/*
 * Copyright (c) 2019. 陈金龙.
 */

package http

import (
	goContext "context"
	"encoding/gob"
	"github.com/jinlongchen/baidu/context"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/brickman-source/golang-utilities/log"
	"github.com/brickman-source/golang-utilities/version"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
	"golang.org/x/oauth2"
)

type Handler struct {
	echoEngine  *echo.Echo
	ctx         *context.Context
	templates   *template.Template
	//accessToken string
}

func NewHttpHandler(ctx *context.Context) *Handler {
	return &Handler{
		ctx: ctx,
	}
}

func (httpH *Handler) GetVersion() string {
	return "0.1"
}

func (httpH *Handler) initRouter() *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	httpH.corsReq(e)

	if httpH.ctx.Config.GetBool("http.secure.csrf") {
		e.Use(middleware.CSRF())
	}
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("520c9f0e66a958ab6cee38253dbf7e20adfdf624"))))

	//httpH.templates = template.Must(template.ParseGlob(httpH.ctx.Config.GetString("http.template.path") + "/*.html"))

	e.POST("/audit/image", httpH.AuditImageHandler)
	e.POST("/audit/text", httpH.AuditTextHandler)
	e.POST("/face-set/add", httpH.FaceSetUserAddHandler)
	e.POST("/face/search", httpH.FaceSearchHandler)
	e.POST("/face/multi-search", httpH.FaceMultiSearchHandler)

	e.GET("/version", func(ctx echo.Context) error {
		ctx.String(200, version.String(httpH.ctx.Config.GetString("application.name"), httpH.GetVersion()))
		return nil
	})

	httpH.echoEngine = e
	return e
}
func (httpH *Handler) corsReq(e *echo.Echo) {
	allowOrigins := httpH.ctx.Config.GetStringSlice("http.header.allowOrigins")
	allowCredentials := httpH.ctx.Config.GetBool("http.header.allowCredentials")
	exposeHeaders := httpH.ctx.Config.GetStringSlice("http.header.exposeHeaders")
	allowHeaders := httpH.ctx.Config.GetStringSlice("http.header.allowHeaders")
	allowMethods := httpH.ctx.Config.GetStringSlice("http.header.allowMethods")

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowCredentials: allowCredentials,
		ExposeHeaders:    exposeHeaders, //[]string{"X-Total-Count"},
		AllowHeaders:     allowHeaders,  //[]string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAccessControlExposeHeaders, echo.HeaderXCSRFToken, "X-Total-Count"},
		AllowMethods:     allowMethods,  //[]string{echo.POST, echo.GET, echo.OPTIONS},
	}))
}
func (httpH *Handler) initTemplates() {
	templ := template.New("")
	templatePath := httpH.ctx.Config.GetString("http.template.path")
	if templatePath == "" {
		return
	}
	err := filepath.Walk(templatePath,
		func(path string, info os.FileInfo, err error) error {
			if strings.Contains(path, ".html") {
				_, err = templ.ParseFiles(path)
				if err != nil {
					log.Errorf( "parse html file err:%s", err.Error())
				}
			}
			return err
		})
	if err != nil {
		panic(err)
	}
	httpH.templates = templ
}

func (httpH *Handler) Serve() {
	gob.Register(oauth2.Token{})

	httpH.initRouter()
	//httpH.initTemplates()

	err := httpH.echoEngine.Start(httpH.ctx.Config.GetString("network.http.listenOn"))
	if err != nil {
		log.Infof( "%s", err.Error())
	}
}
func (httpH *Handler) Exit() {
	_ = httpH.echoEngine.Shutdown(goContext.Background())
}
