package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/jasonlvhit/gocron"
)

type Camera struct {
	Camera string
	URL    string
}

func task() {
	fmt.Println("test")
}

func main() {

	// downloads := make(chan int)
	var wg sync.WaitGroup

	csvFile, _ := os.Open("cameras.csv")
	reader := csv.NewReader(bufio.NewReader(csvFile))

	var cameras []Camera

	gocron.Every(1).Second().Do(task)

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		camera := Camera{line[0], strings.Trim(line[6], " ")}
		cameras = append(cameras, camera)
	}

	for cameranum := range cameras {
		wg.Add(1)
		go downloadImg(cameras[cameranum], &wg)
	}

	wg.Wait()
}

func downloadImg(camera Camera, w *sync.WaitGroup) {

	response, err := http.Get(camera.URL)
	if err != nil {
		w.Done()
		return
	}

	defer response.Body.Close()

	if response.Header.Get("Content-Type") != "image/jpeg" {
		fmt.Printf("Not found: %s\n", camera.Camera)
	}

	file, err := os.Create(fmt.Sprintf("./cameras/%s.jpg", camera.Camera))
	if err != nil {
		w.Done()
		return
	}

	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println(err)
	}
	file.Close()
	w.Done()
}
