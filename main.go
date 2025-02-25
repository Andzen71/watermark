package main

import (
	"fmt"
	"image"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
)

func main() {
	http.HandleFunc("/nikita", nikitaHandler)

	http.HandleFunc("/watermark", WatermarkHandler)
	fmt.Println("Server is running on :8080")
	http.ListenAndServe(":8080", nil)
}
func nikitaHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello ebat\n")
}
func WatermarkHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed", http.StatusMethodNotAllowed)
		return
	}

	// Открываем обрабатываемое изображение
	postImage, err := imaging.Open("perf1.png")
	if err != nil {
		log.Fatalf("Failed to open image: %v", err)
	}

	// Открываем PNG с прозрачностью
	watermark, err := imaging.Open("watermark.png")
	if err != nil {
		log.Fatalf("Failed to open watermark: %v", err)
	}

	// Параметры сетки
	cols := 3
	rows := 3
	margin := 20
	spacing := 10

	// Рассчитываем доступную область для сетки
	piWidth := postImage.Bounds().Dx()
	piHeight := postImage.Bounds().Dy()

	// Вычисляем максимальный размер для одного элемента
	cellWidth := (piWidth - 2*margin - (cols-1)*spacing) / cols
	cellHeight := (piHeight - 2*margin - (rows-1)*spacing) / rows

	// Масштабируем PNG с сохранением пропорций
	scaledWatermark := imaging.Fit(
		watermark,
		cellWidth,
		cellHeight,
		imaging.Lanczos,
	)

	// Рассчитываем смещение для центрирования
	offsetX := (cellWidth - scaledWatermark.Bounds().Dx()) / 2
	offsetY := (cellHeight - scaledWatermark.Bounds().Dy()) / 2

	result := scaledWatermark
	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			// Позиция ячейки
			posX := margin + col*(cellWidth+spacing) + offsetX
			posY := margin + row*(cellHeight+spacing) + offsetY

			// Наложение масштабированного изображения
			result = imaging.Overlay(
				result,
				scaledWatermark,
				image.Pt(posX, posY),
				0.5,
			)
		}
	}

	// Сохраняем результат
	err = imaging.Save(result, "watermarkedImage.png")
	if err != nil {
		log.Fatalf("Failed to save result: %v", err)
	}

	// Открываем PNG-изображение
	img, err := imaging.Open("watermarkedImage.png")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}

	// Конвертируем в JPG с качеством 95%
	err = imaging.Save(img, "watermarkedImage.jpg", imaging.JPEGQuality(95))
	if err != nil {
		log.Fatalf("Failed to save jpg: %v", err)
	}

	// Открыть исходный файл (PNG/JPG)
	srcFile, err := os.Open("watermarkedImage.jpg")
	if err != nil {
		log.Fatal("Failed to open result:", err)
	}
	defer srcFile.Close()

	// Создать целевой файл
	dstFile, err := os.Create("watermarkedImageCopy.jpg")
	if err != nil {
		log.Fatal("Failed to create a file:", err)
	}
	defer dstFile.Close()

	// Копирование через io.Copy
	bytesCopied, err := io.Copy(dstFile, srcFile)
	if err != nil {
		log.Fatal("Failed to copy:", err)
		log.Printf("Copied %d byte", bytesCopied)

		// Отправляем  обработанное изображение в ответе
		w.Header().Set("Content-Type", "image/png")
		if _, err := io.Copy(w, dstFile); err != nil {
			http.Error(w, "Failed to send result", http.StatusInternalServerError)
			return
		}
		fmt.Println("finished request")
	}
}
