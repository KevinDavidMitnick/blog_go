package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/toolkits/file"
	"github.com/ulricqin/blog/config"
	"github.com/ulricqin/blog/http/render"
	"github.com/ulricqin/blog/utils"
)

// 默认将文件存放在static/upload目录，按照年份做一下简单目录划分，个人blog，每年产生的图片不会特别多
func Upload(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	year := now.Format("2006")
	ns := now.UnixNano()

	formFile, header, err := r.FormFile("image")
	if err != nil {
		render.Error(w, err)
		return
	}
	defer formFile.Close()

	basename := fmt.Sprintf("/%s/%d%s", year, ns, file.Ext(header.Filename))
	toFilePath := "static/" + config.G.Upload + basename

	err = file.EnsureDir(file.Dir(toFilePath))
	if err != nil {
		render.Error(w, err)
		return
	}

	toFile, err := os.OpenFile(toFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		render.Error(w, err)
		return
	}
	defer toFile.Close()

	_, err = io.Copy(toFile, formFile)
	if err != nil {
		render.Error(w, err)
		return
	}

	imgUrl := "/" + config.G.Upload + basename
	if config.G.Qiniu != nil && config.G.Qiniu.Enabled {
		addr, err := utils.UploadQiniu(toFilePath, config.G.Upload+basename)
		if err != nil {
			render.Error(w, err)
			return
		}
		imgUrl = addr
	}

	render.Data(w, imgUrl, nil)
}
