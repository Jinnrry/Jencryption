package core

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// 使用md5算法，取最后32位转为uint64
func md5hash(str []byte) uint64 {
	h := md5.New()
	h.Write(str)
	sum := h.Sum(nil)
	return bytesToUint64(sum[12:])
}

func uint64ToBytes(num uint64) []byte {
	return []byte(strconv.FormatUint(num, 10))
	//var buf = make([]byte, 8)
	//binary.BigEndian.PutUint64(buf, num)
	//return buf
}

func bytesToUint64(b []byte) uint64 {
	return uint64(binary.BigEndian.Uint32(b))
}

func point2location(point, width, height int) (x int, y int) {
	return point % (width), point / (width)
}

func OpenImg(src string) (image.Image, error) {
	ff, err := ioutil.ReadFile(src) //读取文件
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	bbb := bytes.NewBuffer(ff)
	m, _, err := image.Decode(bbb)
	return m, err
}

func shuffleColor(data [256]uint8, passwordHash uint64) [256]uint8 {
	var i uint64 = 0
	for ; i < 256; i++ {
		tmp := data[i]
		idx := md5hash(uint64ToBytes(passwordHash+i)) % 256
		data[i] = data[idx]
		data[idx] = tmp
	}
	return data
}

func shufflePoint(data []int, passwordHash uint64) []int {
	var i uint64 = 0
	var lens = uint64(len(data))
	for ; i < lens; i++ {
		tmp := data[i]
		idx := md5hash(uint64ToBytes(passwordHash+i)) % lens
		data[i] = data[idx]
		data[idx] = tmp
	}
	return data
}

// 把一个数组的值和索引对调
func evertArray(data []int) []int {
	var lens = len(data)
	ret := make([]int, lens)
	for oldIdx, newIdx := range data {
		ret[newIdx] = oldIdx
	}
	return ret
}

func everyUint8Array(data [256]uint8) [256]uint8 {
	ret := [256]uint8{}
	for oldIdx, newIdx := range data {
		ret[newIdx] = uint8(oldIdx)
	}
	return ret
}

func getNewPoint(pointsOffset []int, index int, imgWidth, imgHeight int) (x, y int) {
	newIdx := pointsOffset[index]
	return point2location(newIdx, imgWidth, imgHeight)
}

// 解密
func Decrypt(img image.Image, password string) *image.NRGBA {
	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()
	imgSize := imgHeight * imgWidth

	newRgba := image.NewNRGBA(bounds) //new 一个新的图片
	_ = newRgba
	// 密码hash
	hash := md5hash([]byte(password))

	var pointOffset []int = make([]int, imgSize)
	for i := 0; i < imgSize; i++ {
		pointOffset[i] = i
	}

	var colorOffset [256]uint8
	for i := 0; i < 256; i++ {
		colorOffset[i] = uint8(i)
	}

	rOffset := shuffleColor(colorOffset, hash)
	gOffset := shuffleColor(rOffset, hash)
	bOffset := shuffleColor(gOffset, hash)

	rOffset = everyUint8Array(rOffset)
	gOffset = everyUint8Array(gOffset)
	bOffset = everyUint8Array(bOffset)

	pointOffset = shufflePoint(pointOffset, hash)
	pointOffset = evertArray(pointOffset)

	for i := 0; i < imgSize; i++ {
		oldX, oldY := point2location(i, imgWidth, imgHeight)
		newX, newY := getNewPoint(pointOffset, i, imgWidth, imgHeight)

		c := img.At(oldX, oldY)

		var (
			newR uint8
			newG uint8
			newB uint8
			newA uint8
		)

		switch c.(type) {
		case color.NRGBA:
			newR = c.(color.NRGBA).R
			newG = c.(color.NRGBA).G
			newB = c.(color.NRGBA).B
			newA = c.(color.NRGBA).A
		default:
			r, g, b, a := c.RGBA()
			newR = uint8(r >> 8)
			newG = uint8(g >> 8)
			newB = uint8(b >> 8)
			newA = uint8(a >> 8)
		}
		// js canvas中没法获取未预乘的RGBA值，因此针对携带alpha 通道值的点不做颜色偏移
		if newA != 255 {
			newRgba.SetNRGBA(newX, newY, color.NRGBA{newR, newG, newB, newA})
		} else {
			newRgba.SetNRGBA(newX, newY, color.NRGBA{rOffset[newR], gOffset[newG], bOffset[newB], newA})
		}

	}

	return newRgba
}

// 加密
func Encrypt(img image.Image, password string) *image.NRGBA {
	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()
	imgSize := imgHeight * imgWidth

	newRgba := image.NewNRGBA(bounds) //new 一个新的图片

	// 密码hash
	hash := md5hash([]byte(password))

	var pointOffset []int = make([]int, imgSize)
	for i := 0; i < imgSize; i++ {
		pointOffset[i] = i
	}

	var colorOffset [256]uint8
	for i := 0; i < 256; i++ {
		colorOffset[i] = uint8(i)
	}

	rOffset := shuffleColor(colorOffset, hash)
	gOffset := shuffleColor(rOffset, hash)
	bOffset := shuffleColor(gOffset, hash)

	pointOffset = shufflePoint(pointOffset, hash)
	for i := 0; i < imgSize; i++ {
		newX, newY := getNewPoint(pointOffset, i, imgWidth, imgHeight)
		oldX, oldY := point2location(i, imgWidth, imgHeight)
		c := img.At(oldX, oldY)

		var (
			newR uint8
			newG uint8
			newB uint8
			newA uint8
		)

		switch c.(type) {
		case color.NRGBA:
			newR = c.(color.NRGBA).R
			newG = c.(color.NRGBA).G
			newB = c.(color.NRGBA).B
			newA = c.(color.NRGBA).A
		default:
			r, g, b, a := c.RGBA()
			newR = uint8(r >> 8)
			newG = uint8(g >> 8)
			newB = uint8(b >> 8)
			newA = uint8(a >> 8)
		}

		// js canvas中没法获取未预乘的RGBA值，因此针对携带alpha 通道值的点不做颜色偏移
		if newA != 255 {
			newRgba.SetNRGBA(newX, newY, color.NRGBA{newR, newG, newB, newA})
		} else {
			newRgba.SetNRGBA(newX, newY, color.NRGBA{rOffset[newR], gOffset[newG], bOffset[newB], newA})
		}

	}

	return newRgba

}

//保存图片
func SaveImg(inputName string, file *os.File, rgba *image.NRGBA) {
	if strings.HasSuffix(inputName, "jpg") || strings.HasSuffix(inputName, "jpeg") {
		//jpeg.Encode(file, rgba, &jpeg.Options{Quality: 100})  // bug : jpeg编码后颜色值会变化
		png.Encode(file, rgba)
	} else if strings.HasSuffix(inputName, "png") {
		png.Encode(file, rgba)
	} else if strings.HasSuffix(inputName, "gif") {
		//gif.Encode(file, rgba, nil)
		fmt.Errorf("不支持的图片格式")
	} else {
		fmt.Errorf("不支持的图片格式")
	}
}
