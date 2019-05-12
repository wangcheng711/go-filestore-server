package meta

import (
	"fmt"
	"go-filestore-server/db"
	"sort"
)

// 文件元信息结构
type FileMeta struct {
	// 文件哈希码
	FileSha1 string
	// 文件名称
	FileName string
	// 文件大小
	FileSize int64
	// 文件路径
	Location string
	// 更新时间
	UploadAt string
}

var fileMetas map[string]FileMeta

func init()  {
	fileMetas = make(map[string]FileMeta)
}

// 新增/更新文件元信息
func UploadFileMeta(meta FileMeta)  {
	fmt.Printf("name:%s,fileHash:%s \n",meta.FileName,meta.FileSha1)
	fileMetas[meta.FileSha1]=meta
}
// 新增/更新文件元信息到数据库中
func UploadFileMetaDB(fileMeta FileMeta)bool  {
	return db.OnFileUploadFinished(fileMeta.FileSha1,fileMeta.FileName,
		fileMeta.FileSize,fileMeta.Location)
}

// 从数据库中获取文件元信息
func GetFileMetaDB(filehash string)(FileMeta,error)  {
	fileMeta := FileMeta{}
	tfile, err := db.GetFileMeta(filehash)
	if err != nil{
		return fileMeta,err
	}
	fileMeta.FileSize=tfile.FileSize.Int64
	fileMeta.FileName=tfile.FileName.String
	fileMeta.Location=tfile.FileAddr.String
	fileMeta.FileSha1 = tfile.FileHash
	return fileMeta,nil
}

// 通过filesha1 获取文件元信息
func GetFileMeta(fileSha1 string)FileMeta  {
	return fileMetas[fileSha1]
}

// 获取批量文件信息列表
func GetLastFileMetas(count int)[]FileMeta  {
	fileMetaArray := make([]FileMeta,len(fileMetas))
	for _,value:=range fileMetas{
		fileMetaArray=append(fileMetaArray,value)
	}
	sort.Sort(ByUploadTime(fileMetaArray))
	return fileMetaArray[0:count]
}

// 删除文件元信息
func RemoveFileMeta(filehash string)  {
	delete(fileMetas,filehash)
}
