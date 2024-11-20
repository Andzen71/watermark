package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"example.com/watermark"
)

func main() {

	http.HandleFunc("/watermark", WatermarkHandler)
	fmt.Println("Server is running on :8081")
	http.ListenAndServe(":8081", nil)
}

func WatermarkHandler(w http.ResponseWriter, r *http.Request) {

	// // Дополнительный функционал чтения картинки из файла
	// file, err := os.Open("perf1.png")
	// if err != nil {
	// 	fmt.Println("Error opening perf1.png file", err)
	// 	return
	// }
	// defer file.Close()

	headers := r.Header

	// Выводим все заголовки запроса
	for key, values := range headers {
		for _, value := range values {
			fmt.Printf("%s: %s\n", key, value)
		}
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	// Получаем тело запроса (изображение)
	imageData, err := ioutil.ReadAll(r.Body)
	// image, _, err := r.FormFile("image") Это получает картинку из формы. А Выше из файла
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// str1 := string(imageData)
	// fmt.Println(str1)

	// for _, b := range imageData {
	// 	fmt.Printf("%02X ", b)
	// }

	// decodedData, err := url.QueryUnescape(str1)
	// if err != nil {
	// 	fmt.Println("Ошибка раскодирования данных:", err)
	// 	return
	// }

	// imageData1, err := ioutil.ReadAll(decodedData)
	// // image, _, err := r.FormFile("image") Это получает картинку из формы. А Выше из файла
	// if err != nil {
	// 	http.Error(w, "Failed to read request body", http.StatusBadRequest)
	// 	return
	// }
	// defer image.Close() Это при получении картинки из формы

	// Создаем Reader из полученных данных. Это только если получаем картинку из тела запроса
	// imageData64, err := base64.StdEncoding.Decode(imageData)
	// if err != nil {
	// 	fmt.Println("Ошибка декодирования base64:", err)
	// 	return
	// }
	// imageReader := bytes.NewReader(imageData1)

	// if reader == imageReader {
	// 	fmt.Println("yes")
	// } else {
	// 	fmt.Println("no")
	// }

	// contentType := http.DetectContentType(imageData)
	// fmt.Println(contentType)
	// // fmt.Println(string(imageData))

	// // раскодируем данные
	// decodedData, err := url.QueryUnescape(string(imageData))
	// if err != nil {
	// 	http.Error(w, "Failed to decode image data", http.StatusBadRequest)
	// 	return
	// }

	// img, _, err := image.Decode(bytes.NewReader(imageData))
	// if err != nil {
	// 	fmt.Println("Failed to decode image:", err)
	// 	return
	// }

	// fmt.Println(img)

	// Применяем водяной знак
	watermarkedImage, err := watermark.ApplyWatermark(bytes.NewReader(imageData))
	if err != nil {
		http.Error(w, "Failed to apply watermark", http.StatusInternalServerError)
		return
	}

	// Отправляем водяной знак в ответе
	w.Header().Set("Content-Type", "image/png")
	if _, err := io.Copy(w, watermarkedImage); err != nil {
		http.Error(w, "Failed to send watermarked image", http.StatusInternalServerError)
		return
	}
	fmt.Println("finished request")
}
