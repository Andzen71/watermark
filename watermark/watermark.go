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
	// Читаем полученное изображение
	img, _, err := image.Decode(imgReader)
	if err != nil {
		fmt.Println("cant't decode image", err)
		return nil, err
	}

	// Создаем водяной знак (здесь предполагается, что watermark.png находится в watermark каталоге относительно main файла)

	// Открываем файл с водяным знаком (существует ли файл вообще)
	watermarkFile, err := os.Open("watermark/watermark.png")
	if err != nil {
		return nil, err
	}
	defer watermarkFile.Close()

	// Читаем файл с водяным знаком (Этот файл целый и не пустой?)
	watermarkBytes, err := ioutil.ReadAll(watermarkFile)
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return nil, err
	}
	reader := bytes.NewReader(watermarkBytes)

	// Декодируем файл с водяным знаком (Этот файл - картинка?)
	watermarkImg, _, err := image.Decode(reader)
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
