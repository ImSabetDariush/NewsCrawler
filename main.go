package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func createTextSourceFolder() {
	folderPath := "text_source"
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.Mkdir(folderPath, 0755)
		if err != nil {
			log.Fatal("Error creating folder: ", err)
		}
	}
}

func generateFileName(url string) string {
	cleanedURL := strings.Replace(url, "https://", "", -1)
	cleanedURL = strings.Replace(cleanedURL, "http://", "", -1)
	cleanedURL = strings.Replace(cleanedURL, "/", "_", -1)
	cleanedURL = strings.Replace(cleanedURL, ":", "_", -1)
	return filepath.Join("text_source", cleanedURL+".txt")
}

func extractAndSaveText(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching the website: ", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("Error: Status code %d %s\n", resp.StatusCode, resp.Status)
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Println("Error parsing HTML: ", err)
		return
	}

	var textContent string
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		textContent += s.Text() + "\n"
	})

	fileName := generateFileName(url)

	file, err := os.Create(fileName)
	if err != nil {
		log.Println("Error creating file: ", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(textContent)
	if err != nil {
		log.Println("Error writing to file: ", err)
		return
	}

	fmt.Printf("Text successfully extracted and saved to %s.\n", fileName)
}

func main() {
	createTextSourceFolder()

	for {
		var url string
		fmt.Print("Please enter the website URL (or type 'exit' to quit): ")
		fmt.Scan(&url)

		if strings.ToLower(url) == "exit" {
			fmt.Println("Exiting the program. Goodbye!")
			break
		}

		extractAndSaveText(url)
	}
}
