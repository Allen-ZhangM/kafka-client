package kago

import (
	"io/ioutil"
	"os"
	"strings"
)

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func ListDir(dirPth string, suffix string) (files []string, err error) {
	files = make([]string, 0, 10)

	err = PathExistsOrCreate(dirPth)
	if err != nil {
		return nil, err
	}

	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}

	//PthSep := string(os.PathSeparator)
	suffix = strings.ToUpper(suffix) //忽略后缀匹配的大小写

	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			continue
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
			files = append(files, fi.Name())
		}
	}

	return files, nil
}

// 判断文件夹是否存在,不存在则创建
func PathExistsOrCreate(path string) error {
	var exist bool
	_, err := os.Stat(path)
	if err == nil {
		exist = true
	}
	if os.IsNotExist(err) {
		exist = false
	}

	if !exist {
		// 创建文件夹
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}

func Max(a, b int64) int64 {
	if a > b {
		return a
	} else {
		return b
	}
}
