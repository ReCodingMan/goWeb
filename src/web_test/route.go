package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// spaHandler 实现了 http.Handler 接口，所以可以用来处理 HTTP 请求
// 其中 staticPath 用于定义前端静态资源目录（包含js、css 文件）
// indexPath 用于定义入口视图模板文件，通常是 index.html
type spaHandler struct {
	staticPath string
	indexPath  string
}

// 处理 SPA 应用请求（主要是首次访问时入口 HTML 文档和相关静态资源文件，暂不涉及 API 接口）
func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 获取 URL 路径的绝对路径
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// 如果获取失败，返回 400 响应
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// 在 URL 路径前加上静态资源根目录
	path = filepath.Join(h.staticPath, path)

	// 检查对应资源文件是否存在
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// 文件不存在返回入口 HTML 文档内容作为响应
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// 如果期间报错，返回 500 响应
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 一切顺利，则使用 http.FileServer 处理静态资源请求
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

func main()  {
	router := mux.NewRouter()

	spa := spaHandler{staticPath: "dist", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler: router,
		Addr:    "127.0.0.1:8000",
		// 最佳实践：为服务器读写设置超时时间
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}