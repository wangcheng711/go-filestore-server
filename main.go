package main

import (
	"go-filestore-server/handler"
	"net/http"
)

func main()  {
	http.HandleFunc("/file/upload",handler.UploadHandler)
	http.HandleFunc("/file/upload/suc",handler.UploadSucHandler)
	http.HandleFunc("/file/meta",handler.GetFileMetaHandler)
	http.HandleFunc("/file/query",handler.FileQueryHandler)
	http.HandleFunc("/file/download",handler.DownloadHandler)
	http.HandleFunc("/file/update",handler.FileMetaUpateHandler)
	http.HandleFunc("/file/delete",handler.FileDeleteHandler)
	if err := http.ListenAndServe(":8080", nil);err != nil{
		panic(err)
	}
}
