package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

var wg sync.WaitGroup

func main() {

	lines, _ := readLines("/Users/andreasbrommund/Desktop/l.txt")
	numLines := len(lines)
	threads := 100
	span := numLines / threads
	for i := 0; i < threads; i++ {
		wg.Add(1)
		log.Println("Start go", i)
		go work(lines[(i * span) : (i+1)*span])
	}
	wg.Add(1)
	work(lines[threads*span : numLines])
	wg.Wait()
	fmt.Println("DONE")
}

func work(lines []string) {
	base := "http://35.190.155.168/3534986ffd/"
	for _, l := range lines {
		fetch(base + l)
	}
	wg.Done()
}

func fetch(url string) {
	var err error
	var resp *http.Response

	for i := 0; i < 10; i++ {
		resp, err = http.Get(url)
		if err != nil {
			log.Println("Could not get page: "+url, err)
		} else {
			defer resp.Body.Close()
			break
		}
	}

	if err != nil {
		log.Println("Faild to get page", err)
	}

	if resp.StatusCode != http.StatusNotFound {
		fmt.Println(url)
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
