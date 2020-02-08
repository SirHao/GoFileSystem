package meta

import(
	//"sort"
)
//filemeta:文件元信息结构
type FileMeta struct{
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init(){
	fileMetas=make(map[string]FileMeta)
}

//updatefilemeta:新增/更新文件元信息
func UpdateFileMeta(fmeta FileMeta){
	fileMetas[fmeta.FileSha1]=fmeta
}

//GetFileMeta:根据sha1获取filemeta
func GetFileMeta(fileSha1 string)  FileMeta{
	return fileMetas[fileSha1]
}

func GetLastFileMetas(count int)[]FileMeta{
	fMetaArray:=make([]FileMeta,len(fileMetas))
	for _,v:=range fileMetas{
		fMetaArray=append(fMetaArray,v)
	}
	//sort.Sort(ByUploadTime(fMetaArray))
	return fMetaArray[0:count]
}

//RemoveFileMeta:移除meta   todo:线程同步mutex
func RemoveFileMeta(fileSha1 string){
	delete(fileMetas,fileSha1)
}