package watermark

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"io"
	"io/ioutil"
	"os"
)

// https://go.dev/src/image/decode_example_test.go

// ApplyWatermark применяет водяной знак к изображению и возвращает новое изображение
func ApplyWatermark(imgReader io.Reader) (io.Reader, error) {
	// Читаем изображение
	// fmt.Println("В http:", imgReader)
	fmt.Println("yes")
	fmt.Println(imgReader)
	img, imageType, err := image.Decode(imgReader)
	// fmt.Println(img)
	fmt.Println(imageType)
	if err != nil {
		fmt.Println("cant't decode image", err)
		return nil, err
	}

	fileInfo, err := os.Stat("watermark/watermark.png")
	if err != nil {
		fmt.Println("Ошибка при получении информации о файле:", err)
		return nil, err
	}

	fmt.Println("yes1")
	fmt.Println(fileInfo)

	// Создаем водяной знак (здесь предполагается, что watermark.png находится в том же каталоге, что и main файл)
	watermarkFile, err := os.Open("watermark/watermark.png")
	if err != nil {
		return nil, err
	}
	defer watermarkFile.Close()

	watermarkBytes, err := ioutil.ReadAll(watermarkFile)
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return nil, err
	}

	reader := bytes.NewReader(watermarkBytes)

	watermarkImg, _, err := image.Decode(reader)
	fmt.Println(watermarkImg)
	// fmt.Println(watermarkImgType)
	if err != nil {
		fmt.Println("can't decode watermark image")
		return nil, err
	}

	// Накладываем водяной знак на изображение
	resultImg := image.NewRGBA(img.Bounds())
	draw.Draw(resultImg, img.Bounds(), img, image.Point{}, draw.Over)
	draw.Draw(resultImg, watermarkImg.Bounds(), watermarkImg, image.Point{}, draw.Over)

	// Кодируем изображение в формат PNG
	var resultBuf bytes.Buffer
	err = png.Encode(&resultBuf, resultImg)
	if err != nil {
		return nil, err
	}

	return &resultBuf, nil
}
