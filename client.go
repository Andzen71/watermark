package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {

	url := "http://localhost:8081/watermark"

	file, err := os.Open("perf1.png")
	if err != nil {
		fmt.Println("Error opening perf1.png file", err)
		return
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	// image, _, err := r.FormFile("image") Это получает картинку из формы. А Выше из файла
	if err != nil {
		fmt.Println(er)
		return
	}

	payload := bytes.NewReader(fileBytes)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Content-Type", "text/csv")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(res)
	fmt.Println(string(body))

}
