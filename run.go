package main

import (
	"github.com/gorilla/mux"
	"strings"
	"net/http"
	"github.com/gorilla/csrf"
	_ "github.com/yurencloud/yugo/log"
	log "github.com/sirupsen/logrus"
	"strconv"
	"github.com/yurencloud/yugo/config"
)

func staticServer(router *mux.Router) {
	staticConfig := config.Get("static")
	staticArray := strings.Split(staticConfig, ",")
	// 生成一个或多个静态目录，默认static,可自行修改，或添加，以英文逗号分隔
	for index := range staticArray {
		static := staticArray[index]
		router.PathPrefix("/").Handler(http.StripPrefix("/"+static, http.FileServer(http.Dir(static))))
	}
}

func Run() {
	router := mux.NewRouter()
	InitRouter(router)
	staticServer(router)
	configMap := config.GetConfigMap()
	appName := configMap["app.name"]
	log.Info("app: " + appName + ", started at port " + configMap["port"])
	if configMap["csrf.enabled"] == "true" {
		maxAge, _ := strconv.Atoi(configMap["csrf.max.age"])
		CSRF := csrf.Protect(
			[]byte(configMap["csrf.key"]),
			csrf.RequestHeader(configMap["csrf.request.header"]),
			csrf.FieldName(configMap["csrf.field.name"]),
			csrf.MaxAge(maxAge),
			csrf.Secure(false), // 本地开发要加false,生产环境要加true
		)
		http.ListenAndServe(":"+configMap["port"], CSRF(router))
	}else{
		http.ListenAndServe(":"+configMap["port"], router)
	}
}
