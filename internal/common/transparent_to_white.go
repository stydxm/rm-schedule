package common

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
)

// ConvertTransparentToWhite 将输入图像的透明背景转为纯白色，返回处理后的图像字节数组
// 支持PNG和JPEG格式
func ConvertTransparentToWhite(input []byte) ([]byte, error) {
	// 从内存中读取图像数据
	reader := bytes.NewReader(input)

	// 解码图像
	img, format, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("解码图像失败: %v", err)
	}

	// 处理透明通道，将透明部分替换为白色
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// 创建一个新的RGBA图像，背景为白色
	rgba := image.NewRGBA(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 获取像素的RGBA值
			r, g, b, a := img.At(x, y).RGBA()

			// 如果像素完全透明，直接设为白色
			if a == 0 {
				rgba.SetRGBA(x, y, color.RGBA{R: 255, G: 255, B: 255, A: 255})
				continue
			}

			// 将16位颜色值转换为8位
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			// 计算与白色背景混合后的颜色
			alpha := float64(a8) / 255.0
			r8 = uint8(float64(r8)*alpha + 255*(1-alpha))
			g8 = uint8(float64(g8)*alpha + 255*(1-alpha))
			b8 = uint8(float64(b8)*alpha + 255*(1-alpha))

			rgba.SetRGBA(x, y, color.RGBA{R: r8, G: g8, B: b8, A: 255})
		}
	}

	// 创建内存缓冲区存储处理后的图像数据
	var buffer bytes.Buffer

	// 根据原始格式选择编码方式
	switch format {
	case "png":
		err = png.Encode(&buffer, rgba)
	case "jpeg", "jpg":
		err = jpeg.Encode(&buffer, rgba, nil)
	default:
		return nil, fmt.Errorf("不支持的图像格式: %s", format)
	}

	// 返回处理后的图像数据
	return buffer.Bytes(), err
}

// PNGToJPG 将PNG字节数据转换为JPG字节数据，透明部分将被替换为白色
// 输入: PNG图片的字节数组，JPG质量(1-100)
// 输出: JPG图片的字节数组和可能的错误
func PNGToJPG(pngData []byte, quality int) ([]byte, error) {
	// 检查质量参数是否有效
	if quality < 1 || quality > 100 {
		return nil, fmt.Errorf("质量必须在1-100之间")
	}

	// 从内存中读取PNG数据
	reader := bytes.NewReader(pngData)

	// 解码PNG
	img, _, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("解码PNG失败: %v", err)
	}

	// 处理透明通道，将透明部分替换为白色
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	// 创建一个新的RGBA图像，背景为白色
	rgba := image.NewRGBA(bounds)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// 获取像素的RGBA值
			r, g, b, a := img.At(x, y).RGBA()

			// 如果像素完全透明，直接设为白色
			if a == 0 {
				rgba.SetRGBA(x, y, color.RGBA{R: 255, G: 255, B: 255, A: 255})
				continue
			}

			// 处理半透明像素，与白色背景混合
			// 将16位颜色值转换为8位
			r8 := uint8(r >> 8)
			g8 := uint8(g >> 8)
			b8 := uint8(b >> 8)
			a8 := uint8(a >> 8)

			// 计算与白色背景混合后的颜色
			// 公式: 结果色 = 前景色 * 透明度 + 背景色 * (1 - 透明度)
			alpha := float64(a8) / 255.0
			r8 = uint8(float64(r8)*alpha + 255*(1-alpha))
			g8 = uint8(float64(g8)*alpha + 255*(1-alpha))
			b8 = uint8(float64(b8)*alpha + 255*(1-alpha))

			rgba.SetRGBA(x, y, color.RGBA{R: r8, G: g8, B: b8, A: 255})
		}
	}

	// 创建内存缓冲区存储JPG数据
	var jpgBuffer bytes.Buffer

	// 设置JPG编码选项
	opts := &jpeg.Options{
		Quality: quality,
	}

	// 将处理后的图像编码为JPG并写入缓冲区
	err = jpeg.Encode(&jpgBuffer, rgba, opts)
	if err != nil {
		return nil, fmt.Errorf("编码JPG失败: %v", err)
	}

	// 返回JPG字节数据
	return jpgBuffer.Bytes(), nil
}
