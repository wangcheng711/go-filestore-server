package handler

import (
	"encoding/json"
	"fmt"
	"go-filestore-server/meta"
	"go-filestore-server/util"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

// 处理文件上传
func UploadHandler(w http.ResponseWriter,r *http.Request)  {
	if r.Method == http.MethodGet{
		// 返回上传页面
		bytes, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil{
			io.WriteString(w,"server err:"+err.Error())
		}
		io.WriteString(w,string(bytes))
	}else if r.Method == http.MethodPost{
		// 接收文件流及存储到本地目录
		file, header, err := r.FormFile("file")
		if err != nil{
			fmt.Printf("Failed to get data,err:%s \n",err.Error())
			return
		}
		defer file.Close()
		filemeta := meta.FileMeta{
			FileName:header.Filename,
			Location:"./tmp/" + header.Filename,
			UploadAt:time.Now().Format("2006-01-02 15:04:05"),
		}
		newFile, err := os.Create(filemeta.Location)
		if err != nil{
			fmt.Printf("Failed new File,err:%s \n",err.Error())
			return
		}
		defer newFile.Close()

		if filemeta.FileSize, err = io.Copy(newFile, file);err != nil{
			fmt.Printf("Failed copy File,err:%s \n",err.Error())
			return
		}
		newFile.Seek(0,0)
		filemeta.FileSha1 = util.FileSha1(newFile)
		//meta.UploadFileMeta(filemeta)
		meta.UploadFileMetaDB(filemeta)
		// 重定向
		http.Redirect(w,r,"/file/upload/suc",http.StatusFound)
	}
	
}
// 上传成功
func UploadSucHandler(w http.ResponseWriter,r *http.Request)  {
	io.WriteString(w,"Upload Success")
}

// 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter,r *http.Request)  {
	r.ParseForm()
	fileHash := r.Form["filehash"][0]
	//fileMeta := meta.GetFileMeta(fileHash)
	fileMeta,err := meta.GetFileMetaDB(fileHash)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(fileMeta)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
	return

}

// 查询批量的文件信息
func FileQueryHandler(w http.ResponseWriter,r *http.Request)  {
	r.ParseForm()
	limit, _ := strconv.Atoi(r.Form.Get("limit"))
	fileMetas := meta.GetLastFileMetas(limit)
	data, err := json.Marshal(fileMetas)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
	return
}

// 根据hash下载文件
func DownloadHandler(w http.ResponseWriter,r *http.Request)  {
	r.ParseForm()
	fileHash := r.Form.Get("filehash")
	fileMeta := meta.GetFileMeta(fileHash)
	file, err := os.Open(fileMeta.Location)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := ioutil.ReadAll(file)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("content-disposition", "attachment; filename=\""+fileMeta.FileName+"\"")
	w.Write(data)
}

// 修改文件元信息(重命名)
func FileMetaUpateHandler(w http.ResponseWriter,r *http.Request)  {
	r.ParseForm()
	opType := r.Form.Get("op")
	filehash := r.Form.Get("filehash")
	fileName := r.Form.Get("filename")

	if opType != "0"{
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != http.MethodPost{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	fileMeta := meta.GetFileMeta(filehash)
	fileMeta.FileName = fileName
	meta.UploadFileMeta(fileMeta)

	data, err := json.Marshal(fileMeta)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

// 删除文件及元信息
func FileDeleteHandler(w http.ResponseWriter,r *http.Request)  {
	r.ParseForm()
	filehash := r.Form.Get("filehash")

	// 删除文件
	fileMeta := meta.GetFileMeta(filehash)
	os.Remove(fileMeta.Location)

	// 删除文件元信息
	meta.RemoveFileMeta(filehash)


	w.WriteHeader(http.StatusOK)
}