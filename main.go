package main

import(
	"net/http"
	"fmt"
	"./handler"
)
func main()  {
	http.HandleFunc("/file/upload",handler.UploadHandler)
	http.HandleFunc("/file/upload/suc",handler.UploadSucHandler)
	http.HandleFunc("/file/meta",handler.GetFileMetaHandler)
//  http.HandleFunc("/file/",handler.)
	http.HandleFunc("/file/query",handler.FileQueryHandler)
	http.HandleFunc("/file/download",handler.DownloadHandler)
	http.HandleFunc("/file/update",handler.FileUpdateMetaHandler)
	http.HandleFunc("/file/delete",handler.FileDeleteHandler)
	err:=http.ListenAndServe(":8888",nil)
	if err!=nil{
		fmt.Printf("Fail to start:%s",err.Error())
	}
	
}