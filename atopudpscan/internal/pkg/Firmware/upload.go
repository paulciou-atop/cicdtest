package Firmware

import (
	"encoding/json"
	"io"
	"net/http"
	"nms/api/v1/atopudpscan"
	"nms/atopudpscan/internal/pkg/AtopResponse"
	"os"
	"time"
)

func HandleFileUpload(w http.ResponseWriter, r *http.Request, pathParams map[string]string) {
	fwname := time.Now().Format("20060102150405") + ".dld"
	res := AtopResponse.NewResponse(r, w)
	file, _, err := r.FormFile("uploadfile")
	if err != nil {
		f := &atopudpscan.UploadStatus{Code: atopudpscan.UploadStatusCode_Failed, Message: err.Error()}
		j, _ := json.Marshal(f)
		res.SendResponse(j)
		return
	}
	dst, err := os.Create(fwname)
	if err != nil {
		f := &atopudpscan.UploadStatus{Code: atopudpscan.UploadStatusCode_Failed, Message: err.Error()}
		j, _ := json.Marshal(f)
		res.SendResponse(j)
		return
	}

	if _, err := io.Copy(dst, file); err != nil { //copy upload file  to dst
		f := &atopudpscan.UploadStatus{Code: atopudpscan.UploadStatusCode_Failed, Message: err.Error()}
		j, _ := json.Marshal(f)
		res.SendResponse(j)
		return
	}
	defer dst.Close()
	f := &atopudpscan.UploadStatus{Code: atopudpscan.UploadStatusCode_Ok, FileName: fwname}
	j, _ := json.Marshal(f)
	res.SendResponse(j)
}
