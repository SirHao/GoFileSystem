package handler

import (
	"io"
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
	"time"
	"strconv"
	"encoding/json"
	"../meta"
	"../util"
)

func UploadHandler(w http.ResponseWriter,r *http.Request){
	
	if r.Method == "GET"{   //get request
		data,err := ioutil.ReadFile("./static/view/index.html")
		if err!=nil{
			io.WriteString(w,"[handler.go uploadHandler]internal data error!\n")
			return 
		}
		io.WriteString(w,string(data))

	}else if r.Method == "POST"{   //post request

		file,head,err := r.FormFile("uploadfile")
		if err!=nil {
			fmt.Printf("[handler.go uploadHandler]Failed to get data:%s\n",err.Error())
			return
		}
		defer file.Close()

		fileMeta := meta.FileMeta{
			FileName:head.Filename,
			Location:"/tmp/"+head.Filename,
			UploadAt:time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile,err := os.Create(fileMeta.Location)
		if err!=nil {
			fmt.Printf("[handler.go uploadHandler]Failed to create file:%s\n",err.Error())
			return
		}
		defer newFile.Close()

		fileMeta.FileSize,err = io.Copy(newFile,file)
		if err!=nil {
			fmt.Printf("[handler.go uploadHandler]Failed to save file:%s\n",err.Error())
			return
		}
		newFile.Seek(0,0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		fmt.Printf("【handler.go upload】fileupload with name:%s sha:%s\n",string(fileMeta.FileName),string(fileMeta.FileSha1))
		meta.UpdateFileMeta(fileMeta)

		http.Redirect(w,r,"/file/upload/suc",http.StatusFound)
	}

}

func UploadSucHandler(w http.ResponseWriter,r *http.Request)  {
	io.WriteString(w,"Upload finished!\n")
	
}

//GetFileMetaHandler 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()

	filehash:=r.Form["filehash"][0]
	fMeta :=meta.GetFileMeta(filehash)
	data,err:=json.Marshal(fMeta)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

//FileQueryHandler:获取最近的 'limit' 个metas
func FileQueryHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()

	limitCnt,_:=strconv.Atoi(r.Form.Get("limit"))
	fileMetas :=meta.GetLastFileMetas(limitCnt)
	data,err:=json.Marshal(fileMetas)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func DownloadHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	fsha1:=r.Form.Get("filehash")
	fm:=meta.GetFileMeta(fsha1)

	f,err:=os.Open(fm.Location)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
 		return
	}
	defer f.Close()

	data,err:=ioutil.ReadAll(f)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
 		return
	}

	fmt.Printf("【handler.go download】file download with name:%s  sha1:%s\n",string(fm.FileName),string(fm.FileSha1))
	w.Header().Set("Content-Type","application/octect-stream")
	w.Header().Set("content-disposition","attachment;filename=\""+fm.FileName+"\"")
	w.Write(data)
}


//FileUpdateMetaHandler:file rename
func FileUpdateMetaHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	file_hash :=r.Form.Get("filehash")
	new_fname := r.Form.Get("filename")
	opType:=r.Form.Get("op")

	if opType!="0"{
		w.WriteHeader(http.StatusForbidden)
		return
	}
	fmt.Printf("【handler.go update】 startupdate with newname:%s  filehash:%s\n",string(new_fname),string(file_hash))

	if r.Method !="POST"{
		w.WriteHeader(http.StatusMethodNotAllowed)
	}

	curFileMeta:=meta.GetFileMeta(file_hash)
	curFileMeta.FileName=new_fname
	meta.UpdateFileMeta(curFileMeta)

	data,err:=json.Marshal(curFileMeta)
	if err!=nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

//FileDeleteHandler:删除文件
func FileDeleteHandler(w http.ResponseWriter,r *http.Request){
	r.ParseForm()
	file_hash :=r.Form.Get("filehash")
	//new_fname := r.Form.Get("filename")

	fMeta:=meta.GetFileMeta(file_hash)
	os.Remove(fMeta.Location)
	
	meta.RemoveFileMeta(file_hash)

	w.WriteHeader(http.StatusOK)
}