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
	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}

func WatermarkHandler(w http.ResponseWriter, r *http.Request) {

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
