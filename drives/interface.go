package drives

import (
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path"
	"strings"

	gofiletype "github.com/qwxingzhe/go-file-type"
)

type ObjectStorageDrive interface {
	// PutFile 上传本地文件122
	PutFile(localFilePath string, key string) error
	// PutContent 上传字符串到对象存储
	PutContent(fileInfo FileInfo, key string) error
}

type FileInfo struct {
	Content []byte
	DataLen int64
	Ext     string
}

// GetNetFileInfo 读取网路文件基础信息
func GetNetFileInfo(fileUrl string) FileInfo {
	res, err := http.Get(fileUrl)

	if err != nil {
		panic(err)
	}

	defer func() {
		if ferr := res.Body.Close(); ferr != nil {
			err = ferr
		}
	}()

	if err != nil {
		panic(err)
	}

	dataLen := res.ContentLength
	bytes, _ := ioutil.ReadAll(res.Body)

	// 获取文件后缀
	Ext := gofiletype.GetFileTypeByByte(bytes[:10])

	return FileInfo{
		Content: bytes,
		DataLen: dataLen,
		Ext:     Ext,
	}
}

// GetFileHeaderFileInfo 读取文件流中文件信息
func GetFileHeaderFileInfo(file *multipart.FileHeader) FileInfo {
	fileContent, err := file.Open()
	if err != nil {
		panic(err)
	}

	bytes, err := ioutil.ReadAll(fileContent)
	if err != nil {
		panic(err)
	}

	Ext := gofiletype.GetFileTypeByByte(bytes[:10])

	return FileInfo{
		Content: bytes,
		DataLen: int64(len(bytes)),
		Ext:     Ext,
	}
}

// GetContentInfo 读取字符串基础信息
func GetContentInfo(content string) FileInfo {
	return FileInfo{
		Content: []byte(content),
		DataLen: int64(len(content)),
	}
}

// GetLocalFileInfo 读取本地文件基础信息
func GetLocalFileInfo(localFile string) FileInfo {
	Ext := path.Ext(localFile)
	Ext = strings.Replace(Ext, ".", "", 1)
	return FileInfo{
		Ext: Ext,
	}
}
