/*
 * Copyright (c) 2019. 陈金龙.
 */

package context

import (
	"github.com/brickman-source/golang-utilities/baidu"
	"github.com/brickman-source/golang-utilities/cache"
	"github.com/brickman-source/golang-utilities/config"
	"github.com/jmoiron/sqlx"
)

type Context struct {
	Baidu   *baidu.Baidu
	Config   *config.Config
	Database *sqlx.DB
	Cache    cache.Cache
}
