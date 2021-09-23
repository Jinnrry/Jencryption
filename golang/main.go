package main

import (
	"Jencryption/core"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func main() {
	process := ""
	command := "all"
	passowrd := ""
	// 获取命令行参数
	if len(os.Args) == 3 {
		process = os.Args[1]
		passowrd = os.Args[2]
	} else if len(os.Args) == 4 {
		process = os.Args[1]
		passowrd = os.Args[3]
		command = os.Args[2]
	} else {
		showHelp()
	}

	doEncryption := false
	if process == "encrypt" {
		doEncryption = true
	} else if process == "decrypt" {
		doEncryption = false
	} else {
		showHelp()
	}

	if command == "all" {
		processAllImg(passowrd, "./", doEncryption)
	} else if isDir(command) {
		processAllImg(passowrd, command, doEncryption)
	} else {
		processOneImg(passowrd, "", command, doEncryption)
	}

}

func isDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()

}

func showHelp() {
	fmt.Println("命令错误！")
	fmt.Println("使用说明:")
	fmt.Println("加解密单个文件:  Jencryption [encrypt/decrypt] [文件名] [密码]")
	fmt.Println("加解密当前路径下的全部文件:  Jencryption [encrypt/decrypt] [密码]")
	os.Exit(-1)
}

func processOneImg(pwd, path, fileName string, encryption bool) {
	dir := "./encryption"
	if !encryption {
		dir = "./decryption"
	}

	err := os.Mkdir(dir, os.ModePerm)
	if err != nil {
		if !strings.Contains(err.Error(), "file exists") {
			panic(err)
		}
	}

	img, err := core.OpenImg(path + fileName)
	if err != nil {
		panic(err)
	}
	var imgData *image.NRGBA
	if encryption {
		imgData = core.Encrypt(img, pwd)
	} else {
		imgData = core.Decrypt(img, pwd)
	}

	f, err := os.Create(dir + "/" + fileName)
	if err != nil {
		panic(err)
	}
	core.SaveImg(dir+"/"+fileName, f, imgData)
	f.Close()
	if encryption {
		fmt.Printf("加密")
	} else {
		fmt.Printf("解密")
	}
	fmt.Println(path+fileName, "到", dir+"/"+fileName)

}

func processAllImg(pwd string, dir string, encryption bool) {
	files, _ := ioutil.ReadDir(dir)
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	for _, fi := range files {
		if !fi.IsDir() {
			fileExt := path.Ext(fi.Name())
			if fileExt == ".png" || fileExt == ".jpg" || fileExt == ".jpeg" {
				processOneImg(pwd, dir, fi.Name(), encryption)
			}

		}
	}

}
