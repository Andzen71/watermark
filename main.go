package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"example.com/watermark"
)

func main() {

	available_commands := []string{"server", "watermark", "--help", "-h"}
	use_available_command := false

	if len(os.Args) < 2 {
		fmt.Println("Please, use available commands", available_commands)
		return
	}

	for _, available_command := range available_commands {
		if os.Args[1] == available_command {
			use_available_command = true
			break
		}
	}

	if !use_available_command {
		fmt.Println("No command", os.Args[1], ". Please, use available commands", available_commands)
	}

	if os.Args[1] == "server" {
		http.HandleFunc("/watermark", WatermarkHandler)
		fmt.Println("Server is running on :8081")
		http.ListenAndServe(":8081", nil)
	}

	if os.Args[1] == "watermark" {
		if len(os.Args) < 3 {
			fmt.Println("Please, specify path to image")
			return
		}

		if os.Args[2] == "--help" || os.Args[2] == "-h" {
			// this is the problem of os.Args. It is hard to do that order of args is not important
			fmt.Println("path to image that should be watermarked")
			fmt.Println("path to watermarked image")
			return
		}

		var input_path string = os.Args[2]
		var output_path string = "output_image.png"

		if len(os.Args) >= 4 {
			output_path = os.Args[3]
		}

		WatermarkCli(input_path, output_path)
	}

	if os.Args[1] == "--help" || os.Args[1] == "-h" {
		fmt.Println("server       --use http server to add watermark,")
		fmt.Println("watermark    --use cli command to add watermark")
	}
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

func WatermarkCli(file_path string, output_path string) {
	// Открываем переданный файл (существует ли файл вообще)
	watermarkFile, err := os.Open(file_path)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer watermarkFile.Close()

	// Читаем переданный файл (Этот файл целый и не пустой?)
	watermarkBytes, err := ioutil.ReadAll(watermarkFile)
	if err != nil {
		fmt.Println(err)
		return
	}

	watermarkedImage, err := watermark.ApplyWatermark(bytes.NewReader(watermarkBytes))
	if err != nil {
		fmt.Println("Failed to apply watermark", err)
		return
	}

	// Декодируем изображение из io.Reader
	img, _, err := image.Decode(watermarkedImage)
	if err != nil {
		fmt.Println("Ошибка при декодировании изображения:", err)
		return
	}

	// Открываем файл для записи
	outFile, err := os.Create(output_path)
	if err != nil {
		fmt.Println("Ошибка при создании файла:", err)
		return
	}
	defer outFile.Close() // Закрытие файла после завершения работы

	// Кодируем изображение в формат PNG и записываем его в файл
	err = png.Encode(outFile, img)
	if err != nil {
		fmt.Println("Ошибка при записи изображения:", err)
		return
	}

	fmt.Println("Изображение успешно записано в файл:", output_path)
}
