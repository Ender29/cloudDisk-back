package main

import (
	"archive/zip"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func RandString() string {
	char := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	arr := strings.Split(char, "")
	length := len(arr)
	ran := rand.New(rand.NewSource(time.Now().Unix()))
	randStr := ""
	for i := 0; i < 15; i++ {
		randStr = randStr + arr[ran.Intn(length)]
	}

	return randStr
}

// GetDays
// date2 大于 date1
// -1 失败
func GetDays(format, date1, date2 string) int {
	d1, err := time.ParseInLocation(format, date1, time.Local)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	d2, err := time.ParseInLocation(format, date2, time.Local)
	if err != nil {
		return -1
	}
	return int(d2.Sub(d1).Hours() / 24)
}

// 打包成zip文件
func Zip(src_dir string, zip_file_name string) {

	// 预防：旧文件无法覆盖
	os.RemoveAll(zip_file_name)

	// 创建：zip文件
	zipfile, _ := os.Create(zip_file_name)
	defer zipfile.Close()

	// 打开：zip文件
	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	// 遍历路径信息
	filepath.Walk(src_dir, func(path string, info os.FileInfo, _ error) error {

		// 如果是源路径，提前进行下一个遍历
		if path == src_dir {
			return nil
		}

		// 获取：文件头信息
		header, _ := zip.FileInfoHeader(info)
		header.Name = strings.TrimPrefix(path, src_dir+`/`)

		// 判断：文件是不是文件夹
		if info.IsDir() {
			header.Name += `/`
		} else {
			// 设置：zip的文件压缩算法
			header.Method = zip.Deflate
		}

		// 创建：压缩包头部信息
		writer, _ := archive.CreateHeader(header)
		if !info.IsDir() {
			file, _ := os.Open(path)
			defer file.Close()
			io.Copy(writer, file)
		}
		return nil
	})
}

func main() {

	decoded, err := base64.StdEncoding.DecodeString("\u0003JC^�Z��\"�oBW��V�vCd��\u001C��ˆ�Ȼ�[�KPrq(�YNRA�.w�5�]SH~���\"c�@�0�\u001A\f�\u0017���Ҟ�� U\u001C����^�;۬\n�źO�\u0011�\\�MP�6�6jD��@�{�37\u0018�X\u0011/�z��!�>\u000F\\���\u0012�@�l�x��N)\u0012]��%��������5'T��\u0013d\u0010�|��\u001Dt�.M[\u001D[rH9\u0013�\u001A�\u001B��\u0014�Xj �Άz��d�L�������|���(\u0012�s�\u0017��")
	decodestr := string(decoded)
	fmt.Println(decodestr, err)
}
