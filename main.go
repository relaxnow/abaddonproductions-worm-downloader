package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"crypto/tls"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://abaddonproductions.org/worm/"

	// Skip SSL certificate verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Send an HTTP GET request to the URL
	response, err := client.Get(url)
	if err != nil {
		fmt.Println("Failed to retrieve the webpage:", err)
		return
	}
	defer response.Body.Close()

	// Check if the request was successful
	if response.StatusCode != 200 {
		fmt.Println("Failed to retrieve the webpage. Status code:", response.StatusCode)
		return
	}

	// Parse the HTML content using goquery
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		fmt.Println("Failed to parse HTML:", err)
		return
	}

	// Find all elements using the CSS selector
	doc.Find("article").Each(func(i int, s *goquery.Selection) {
		headerLink := s.Find("header h1 a")
		hrefValue, exists := headerLink.Attr("href")
		if exists {
			fmt.Println("Value of the href attribute:", hrefValue)
		} else {
			fmt.Println("Selected element does not have an href attribute.")
		}

		downloadLink := s.Find("a[title=Download]")
		downloadHref, _ := downloadLink.Attr("href")
		fmt.Println("Downloads: " + downloadHref)

		re := regexp.MustCompile(`https:\/\/abaddonproductions.org\/(\d{4})\/(\d{2})\/(\d{2})\/(.*?)\/`)
		mp3Filename := re.ReplaceAllString(hrefValue, "$1-$2-$3-$4.mp3")
		fmt.Println("New Filename: " + mp3Filename)

		download(downloadHref, "./Wildbow/Worm/", mp3Filename)
	})
}

func download(fileURL string, downloadDir string, filename string) {
	// Create the download directory if it doesn't exist
	err := os.MkdirAll(downloadDir, 0755)
	if err != nil {
		fmt.Println("Failed to create download directory:", err)
		return
	}

	// Create the file in the download directory
	filePath := filepath.Join(downloadDir, filename)

	if _, err := os.Stat(filePath); err == nil {
		fmt.Println("File " + filePath + " already exists, skipping...")
		return
	}

	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Failed to create file:", err)
		return
	}
	defer file.Close()

	// Skip SSL certificate verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Send an HTTP GET request to the URL
	response, err := client.Get(fileURL)
	if err != nil {
		fmt.Println("Failed to download the file:", err)
		return
	}
	defer response.Body.Close()

	// Check if the request was successful
	if response.StatusCode != http.StatusOK {
		fmt.Println("Failed to download the file. Status code:", response.StatusCode)
		return
	}

	// Copy the file contents from the response body to the created file
	_, err = io.Copy(file, response.Body)
	if err != nil {
		fmt.Println("Failed to save file:", err)
		return
	}

	fmt.Println("File downloaded and saved to:", filePath)
}
