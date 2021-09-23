package core

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
	"testing"
)

func TestCreateTestImage(t *testing.T) {
	width := 2
	height := 4
	skip := 255 / (width * height)
	if skip == 0 {
		skip = 1
	}
	newRgba := image.NewNRGBA(image.Rect(0, 0, width, height))
	for i := 0; i < width*height; i++ {
		newRgba.SetNRGBA(i%width, i/width, color.NRGBA{R: uint8(i*skip) % 255, G: uint8(i*skip) % 255, B: uint8(i*skip) % 255, A: 255})
	}
	f, _ := os.Create("./test.png")
	defer f.Close()
	SaveImg("test.png", f, newRgba)
}

func TestEncryption(t *testing.T) {
	img, err := OpenImg("./1.jpg")
	if err != nil {
		t.Error(err)
	}

	//起始点： 467 643 0 0 0
	imgData := Encrypt(img, "jinnrry")
	fmt.Println(imgData.At(0, 0))
	_ = imgData
	f, _ := os.Create("./encryption.png")
	defer f.Close()
	SaveImg("encryption.png", f, imgData)
}

func TestImgSave(t *testing.T) {
	newRgba := image.NewNRGBA(image.Rect(0, 0, 1, 1)) //new image
	newRgba.SetNRGBA(0, 0, color.NRGBA{R: 55, G: 23, B: 14, A: 122})
	f, _ := os.Create("./save.png")
	defer f.Close()
	// save image
	png.Encode(f, newRgba)

	ff, _ := ioutil.ReadFile("./save.png") //read image
	bbb := bytes.NewBuffer(ff)
	m, _, _ := image.Decode(bbb)


	fmt.Println(m.At(0,0).(color.NRGBA).R)

	fmt.Printf("Before: %T: %[1]v\n", newRgba.At(0, 0))
	fmt.Printf("After:  %T: %[1]v\n", m.At(0, 0))



}

func TestDecrypt(t *testing.T) {
	img, err := OpenImg("./encryption.png")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(img.At(0, 0))
	//起始点： 467 643 0 0 0
	imgData := Decrypt(img, "123")

	f, _ := os.Create("./decryption.png")
	defer f.Close()
	SaveImg("decryption.png", f, imgData)
}

func Test_hash(t *testing.T) {
	res := md5hash([]byte("2"))
	fmt.Println(res)
	if res != 3423897132 {
		t.Error("hash 计算错误")
	}

	res = md5hash(uint64ToBytes(1))
	fmt.Println(res)
}

func Test_bytesToInt(t *testing.T) {
	res := bytesToUint64([]byte{229, 129, 73, 157, 123, 112, 22, 104})
	fmt.Println(res)

}

func Test_uint64ToBytes(t *testing.T) {
	res := uint64ToBytes(16537580247410808424)
	fmt.Println(res)
}

func TestEncryptDecrypt(t *testing.T) {
	img, err := OpenImg("./test.png")
	if err != nil {
		t.Error(err)
	}

	//起始点： 467 643 0 0 0
	Encrypt(img, "123")

	img, err = OpenImg("./encryption.png")
	if err != nil {
		t.Error(err)
	}
	//起始点： 467 643 0 0 0
	Decrypt(img, "123")
}

func Test_shuffleColor(t *testing.T) {
	var colorOffset [256]uint8
	for i := 0; i < 256; i++ {
		colorOffset[i] = uint8(i)
	}
	hash := md5hash([]byte("123"))
	fmt.Println(hash)
	rOffset := shuffleColor(colorOffset, hash)
	fmt.Println(rOffset)
}

func Test_shufflePoint(t *testing.T) {
	var colorOffset []int = make([]int, 256)
	for i := 0; i < 256; i++ {
		colorOffset[i] = int(i)
	}
	hash := md5hash([]byte("123"))
	rOffset := shufflePoint(colorOffset, hash)
	fmt.Println(rOffset)
	evert := evertArray(rOffset)
	fmt.Println(evert)
}
