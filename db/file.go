package db

import (
	"database/sql"
	"fmt"
	"go-filestore-server/db/mysql"
)

// 保存文件信息
func OnFileUploadFinished(filehash string,filename string,filesize int64,
	fileaddr string)bool{
	stmt, err := mysql.DBConn().Prepare("insert ignore into tbl_file (`file_sha1`,`file_name`,`file_size`," +
		"`file_addr`,`status`) values (?,?,?,?,1)")
	if err != nil{
		fmt.Printf("Failed to Prepare OnFileUploadFinished,err:%s",err.Error())
		return false
	}
	defer stmt.Close()
	result, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil{
		fmt.Printf("Failed to Exec OnFileUploadFinished,err:%s",err.Error())
		return false
	}
	affected, err := result.RowsAffected()
	if err != nil{
		fmt.Printf("Failed to RowsAffected OnFileUploadFinished,err:%s",err.Error())
		return false
	}
	if affected <= 0{
		fmt.Printf("File with hash:%s has been uploaded before", filehash)
		return false
	}
	return true
}

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// 从数据库中获取文件元信息
func GetFileMeta(filehash string)(*TableFile,error)  {
	stmt, err := mysql.DBConn().Prepare("select file_sha1,file_addr,file_name,file_size from tbl_file " +
		"where file_sha1=? and status=1 limit 1")
	if err != nil{
		fmt.Printf("Failed to Prepare GetFileMeta,err:%s \n",err.Error())
		return nil,err
	}
	defer stmt.Close()
	tfile := TableFile{}
	err = stmt.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileAddr,
		&tfile.FileName, &tfile.FileSize)
	if err != nil{
		fmt.Printf("Failed to QueryRow GetFileMeta,err:%s \n",err.Error())
		return nil,err
	}
	return &tfile,nil
}