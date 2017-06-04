package utils

import (
	"strings"

	"github.com/qiniu/api.v7/kodo"
	"github.com/ulricqin/blog/config"
	"qiniupkg.com/api.v7/kodocli"
)

//构造返回值字段
type PutRet struct {
	Hash string `json:"hash"`
	Key  string `json:"key"`
}

func UploadQiniu(localFile, destName string) (string, error) {
	c := kodo.New(0, nil)

	policy := &kodo.PutPolicy{
		Scope:      config.G.Qiniu.Bucket + ":" + destName,
		Expires:    3600,
		InsertOnly: 1,
	}

	token := c.MakeUptoken(policy)
	zone := 0
	uploader := kodocli.NewUploader(zone, nil)

	var ret PutRet
	err := uploader.PutFile(nil, &ret, token, destName, localFile, nil)
	addr := ""
	if err == nil {
		domain := config.G.Qiniu.Domain
		if !strings.HasPrefix(domain, "http") {
			domain = "http://" + domain
		}

		addr = strings.TrimRight(domain, "/") + "/" + ret.Key
	}
	return addr, err
}
