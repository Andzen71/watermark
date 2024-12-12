package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"image"
	"image/png"
)

func main() {

	// Куда будем отправлять запрос
	url := "http://localhost:8080/watermark"

	// Открываем картинку, которую будем отправлять (она вообще есть?)
	file, err := os.Open("perf1.png")
	if err != nil {
		fmt.Println("Error opening perf1.png file", err)
		return
	}
	defer file.Close()

	// Читаем картинку
	fileBytes, err := ioutil.ReadAll(file)
	// image, _, err := r.FormFile("image") Это получает картинку из формы. А Выше из файла
	if err != nil {
		fmt.Println(err)
		return
	}
	payload := bytes.NewReader(fileBytes)

	// Готовим POST запрос
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Println(err)
	}
	req.Header.Add("Content-Type", "image/png")

	// Отправляем подготовленный запрос
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()

	// Читаем тело ответа
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}
	// Тут мы должны получить картинку с сервера в виде image.png

	// Локально сохраняем картинку с водяным знаком
	img, image.Image, filename string {
  file, err := os.Create(filename)
  if err != nil {
		fmt.Println(err)
	}
  defer res.Body.Close()
	}

	// Дополнительно выводим ответ и тело ответа (в формате строки)
	fmt.Println(res)
	fmt.Println(string(body))

}
