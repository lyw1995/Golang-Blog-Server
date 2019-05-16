package repositories

import (
	"fmt"
	"github.com/track/blogserver/pkg/config"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//文件上传
type UploadRepository struct {
	parentPath string
}

func NewUploadRepository() *UploadRepository {
	return &UploadRepository{
		parentPath: config.Config().ImgPath,
	}
}

//创建文件夹
func createDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		return os.MkdirAll(path, os.ModePerm)
	} else {
		return err
	}
}

//创建文件
func createFile(path string, in io.Reader) error {
	out, err := os.Create(path)
	defer out.Close()
	_, err = io.Copy(out, in)
	return err
}

//保存图片
func (ur *UploadRepository) Save(name string, file multipart.File) (string, error) {
	currentTime := time.Now()
	path := fmt.Sprintf("%s/%d/%d", ur.parentPath, currentTime.Year(), currentTime.Month())
	if err := createDir(path); err != nil {
		return "", err
	}

	filename := filepath.Join(path, "/", name)

	if err := createFile(filename, file); err != nil {
		return "", err
	}

	return strings.Replace(filename, ur.parentPath, "image", -1), nil
}
