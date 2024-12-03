package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"math"
	"os"
)

func main() {
	// 入力画像を読み込む
	inputFile, err := os.Open("input.jpg")
	if err != nil {
		log.Fatalf("画像の読み込みに失敗しました: %v", err)
	}
	defer inputFile.Close()

	img, _, err := image.Decode(inputFile)
	if err != nil {
		log.Fatalf("画像のデコードに失敗しました: %v", err)
	}

	// 画像の幅と高さを取得
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	aspectRatio := float64(width) / float64(height)

	// 比率を定義
	aspectRatio4_3 := 4.0 / 3.0
	aspectRatio16_9 := 16.0 / 9.0

	var newImg image.Image

	// アスペクト比に基づいて処理を行う
	if aspectRatio < aspectRatio4_3 {
		// 縦長の場合、4:3の比率に合わせる
		newWidth := width
		newHeight := int(math.Round(float64(newWidth) / aspectRatio4_3))

		// 背景画像を作成
		bg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
		gray := color.RGBA{128, 128, 128, 255}
		draw.Draw(bg, bg.Bounds(), &image.Uniform{gray}, image.Point{}, draw.Src)

		// 元の画像をリサイズして背景に配置
		scaleW := float64(newWidth) / float64(width)
		scaleH := float64(newHeight) / float64(height)
		scale := math.Min(scaleW, scaleH)

		resizedWidth := int(float64(width) * scale)
		resizedHeight := int(float64(height) * scale)

		resizedImg := resizeImage(img, resizedWidth, resizedHeight)

		xOffset := (newWidth - resizedWidth) / 2
		yOffset := (newHeight - resizedHeight) / 2

		draw.Draw(bg, image.Rect(xOffset, yOffset, xOffset+resizedWidth, yOffset+resizedHeight), resizedImg, image.Point{}, draw.Over)

		newImg = bg
	} else if aspectRatio > aspectRatio16_9 {
		// 横長の場合、16:9の比率に合わせる
		newHeight := height
		newWidth := int(math.Round(float64(newHeight) * aspectRatio16_9))

		// 背景画像を作成
		bg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
		gray := color.RGBA{128, 128, 128, 255}
		draw.Draw(bg, bg.Bounds(), &image.Uniform{gray}, image.Point{}, draw.Src)

		// 元の画像をリサイズして背景に配置
		scaleW := float64(newWidth) / float64(width)
		scaleH := float64(newHeight) / float64(height)
		scale := math.Min(scaleW, scaleH)

		resizedWidth := int(float64(width) * scale)
		resizedHeight := int(float64(height) * scale)

		resizedImg := resizeImage(img, resizedWidth, resizedHeight)

		xOffset := (newWidth - resizedWidth) / 2
		yOffset := (newHeight - resizedHeight) / 2

		draw.Draw(bg, image.Rect(xOffset, yOffset, xOffset+resizedWidth, yOffset+resizedHeight), resizedImg, image.Point{}, draw.Over)

		newImg = bg
	} else {
		// 加工しない場合
		newImg = img
	}

	// 処理した画像を保存
	outputFile, err := os.Create("output.jpg")
	if err != nil {
		log.Fatalf("画像の保存に失敗しました: %v", err)
	}
	defer outputFile.Close()

	// JPEGでエンコードして保存
	err = jpeg.Encode(outputFile, newImg, &jpeg.Options{Quality: 95})
	if err != nil {
		log.Fatalf("画像のエンコードに失敗しました: %v", err)
	}
}

// 簡単な最近傍法によるリサイズ関数
func resizeImage(img image.Image, newWidth, newHeight int) *image.RGBA {
	resizedImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))
	srcBounds := img.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := int(float64(x) * float64(srcWidth) / float64(newWidth))
			srcY := int(float64(y) * float64(srcHeight) / float64(newHeight))
			color := img.At(srcX, srcY)
			resizedImg.Set(x, y, color)
		}
	}

	return resizedImg
}
func getImage(path string) (img image.Image, close func(), err error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, nil, err
	}
	img, _, err = image.Decode(f)
	if err != nil {
		return nil, nil, err
	}

	return img, func() {
		f.Close()
	}, nil
}
