/**
 * Copyright 2014 @ z3q.net.
 * name :
 * author : jarryliu
 * date : 2013-12-16 21:47
 * description : 静态资源包,主要用于提功静态服务器的功能
 * history :
 */

package pub

import (
	"github.com/jsix/gof"
	"go2o/core/variable"
	"log"
	"net/http"
	"strings"
)

// 静态文件
type StaticHandler struct {
}

func (s *StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.ServeFile(w, r, "./public/static"+r.URL.Path)
}

// 图片处理
type ImageFileHandler struct {
	app       gof.App
	upSaveDir string
}

func NewImageFileHandler(a gof.App) *ImageFileHandler {
	return &ImageFileHandler{
		app: a,
	}
}

func (i *ImageFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	//if strings.HasPrefix(path, "/res/") {
	//	http.ServeFile(w, r, "public/static"+path)
	//} else {
	if len(i.upSaveDir) == 0 {
		i.upSaveDir = i.app.Config().GetString(variable.UploadSaveDir)
	}
	http.ServeFile(w, r, i.upSaveDir+path)

	//}
}

func Listen(ch chan bool, app gof.App, addr string) {
	log.Println("** [ Go2o][ Web][ Booted] - Pub server running on", addr)
	h := &pubHandler{}
	s := &StaticHandler{}
	i := NewImageFileHandler(app)
	if err := http.ListenAndServe(addr, h.set(s, i)); err != nil {
		log.Println("** [ Go2o][ Web][ Exit] -", err.Error())
	}
}

type pubHandler struct {
	staticServe *StaticHandler
	imgServe    *ImageFileHandler
}

func (this *pubHandler) set(s *StaticHandler, i *ImageFileHandler) http.Handler {
	this.staticServe = s
	this.imgServe = i
	return this
}

func (this *pubHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	subName := r.Host[:strings.Index(r.Host, ".")+1]
	switch subName {
	case variable.DOMAIN_PREFIX_STATIC:
		this.staticServe.ServeHTTP(w, r)
	case variable.DOMAIN_PREFIX_IMAGE:
		this.imgServe.ServeHTTP(w, r)
	default:
		http.Error(w, "no such file", 404)
	}
}
