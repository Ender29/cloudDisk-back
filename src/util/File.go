package util

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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
	return int(d2.Sub(d1).Hours()/24) + 1
}

// GetRandStr 随机字符串
func GetRandStr(lens int) string {
	char := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	arr := strings.Split(char, "")
	length := len(arr)
	ran := rand.New(rand.NewSource(time.Now().Unix()))
	randStr := ""
	for i := 0; i < lens; i++ {
		randStr = randStr + arr[ran.Intn(length)]
	}
	return randStr
}

// Copy 复制文件
func Copy(src, dir string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()
	dest, err := os.Create(dir)
	if err != nil {
		return err
	}
	defer dest.Close()
	_, err = io.Copy(dest, source)
	return err
}

// Zip 压缩
func Zip(src_dir string, archive *zip.Writer) {
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

//func Zip(srcDir string, zw *zip.Writer) {
//
//	// 遍历路径信息
//	filepath.Walk(srcDir, func(path string, info os.FileInfo, _ error) error {
//
//		// 如果是源路径，提前进行下一个遍历
//		if path == srcDir {
//			return nil
//		}
//
//		// 获取：文件头信息
//		header, _ := zip.FileInfoHeader(info)
//		header.Name = strings.TrimPrefix(path,  filepath.Dir(srcDir) +`/`)
//
//		// 判断：文件是不是文件夹
//		if info.IsDir() {
//			header.Name += `/`
//		} else {
//			// 设置：zip的文件压缩算法
//			header.Method = zip.Deflate
//		}
//
//		// 创建：压缩包头部信息
//		writer, _ := zw.CreateHeader(header)
//		if !info.IsDir() {
//			file, _ := os.Open(path)
//			defer file.Close()
//			io.Copy(writer, file)
//		}
//		return nil
//	})
//}

// IsExist 判断文件是否存在
func IsExist(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// ListDir 返回文件列表
func ListDir(pathname string, s []string) ([]string, error) {
	rd, err := ioutil.ReadDir(pathname)
	if err != nil {
		fmt.Println("read dir fail:", err)
		return s, err
	}

	for _, fi := range rd {
		if !fi.IsDir() {
			fullName := pathname + "/" + fi.Name()
			s = append(s, fullName)
		}
	}
	return s, nil
}
