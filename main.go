package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Lets define the result
type Result struct {
	URL      string
	Size     int64
	Error    error
	Filename string
	Duration time.Duration
}

// Lets start with a normal downloader
func FileDownloader(url string, destDir string) error {
	fileName := filepath.Base(url)
	filePath := filepath.Join(fileName)

	output, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer output.Close()

	timeStart := time.Now()
	fmt.Println("....Downloading", url)

	response, err := http.Get(url)
	if err != nil {
		os.Remove(fileName)
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		os.Remove(fileName)
		log.Fatal(err)
	}

	_, err = io.Copy(output, response.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Download took %s to complete", time.Since(timeStart))

	return nil
}

// A Multiple file downloader
func MultipleFileDownloader(urls []string, destDir string) error {
	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		return err
	}
	timeStart := time.Now()

	for _, url := range urls {
		err := FileDownloader(url, destDir)
		if err != nil {
			return err
		}
		fmt.Println("Error Downloading.....", url, destDir)

	}
	fmt.Printf("Downloading took %s", time.Since(timeStart))
	return nil
}

// A concurrent downloader
func ConcurrentFileDownloader(urls []string, destDir string, maxConcurrent int) error {
	err := os.MkdirAll(destDir, 0755)
	if err != nil {
		return err
	}

	// create a channel for result
	resultChannel := make(chan Result)
	// Use waitgroup to cordinate workers
	var wg sync.WaitGroup
	// Rate limiter
	limiter := make(chan struct{}, maxConcurrent)
	timeStart := time.Now()
	for _, url := range urls {
		// Add a waitgroup
		wg.Add(1)
		go func(url string) {
			// Done with waitgroup
			defer wg.Done()

			limiter <- struct{}{}
			defer func() { <-limiter }()

			

			fileName := filepath.Base(url)
			filePath := filepath.Join(fileName)

			output, err := os.Create(filePath)
			if err != nil {
				resultChannel <- Result{URL: url, Error: err}
			}
			defer output.Close()

			response, err := http.Get(url)
			if err != nil {
				os.Remove(url)
				resultChannel <- Result{URL: url, Error: err}
			}
			defer response.Body.Close()

			if response.StatusCode != http.StatusOK {
				os.Remove(url)
				resultChannel <- Result{URL: url, Error: fmt.Errorf("Bad Request %s\n", response.Status)}
				return
			}

			size, err := io.Copy(output, response.Body)
			if err != nil {
				resultChannel <- Result{URL: url, Error: err}
				return
			}

			timeSince := time.Since(timeStart)

			resultChannel <- Result{URL: url, Filename: fileName, Size: size, Duration: timeSince, Error: nil}
		}(url)
	}

	// Close the channel
	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	var totalSize int64
	var errors []error

	for result := range resultChannel {
		if result.Error != nil {
			fmt.Printf("Error Downloading %s: %s\n", result.URL, result.Error.Error())
			errors = append(errors, result.Error)
		} else {
			totalSize += result.Size
			fmt.Printf("Downloaded %s (%d bytes) in %s\n", result.Filename, result.Size, result.Duration)
		}
	}

	startedSince := time.Since(timeStart)
	fmt.Printf("All download completed in %s, Total: %d bytes\n", startedSince, totalSize)
	if len(errors) > 0 {
		return  fmt.Errorf("Error Downloading: %+v", errors)
	}
	return nil
}
func main() {
	urls := []string{
		"https://cdn-hednb.nitrocdn.com/yYolPsxkHfeoqRKSyqGQlaFpLZMHhVYI/assets/images/optimized/rev-23c0ecc/hpe-photos.s3.us-east-2.amazonaws.com/wp-content/uploads/2025/02/Hennessey-2024-Ford-Mustang-GT-for-Sale-Black-0554-2.jpg",
		"https://cdn-hednb.nitrocdn.com/yYolPsxkHfeoqRKSyqGQlaFpLZMHhVYI/assets/images/optimized/rev-23c0ecc/hpe-photos.s3.us-east-2.amazonaws.com/wp-content/uploads/2025/02/Hennessey-2024-Ford-Mustang-GT-for-Sale-Black-0554-4.jpg", 
		"https://cdn-hednb.nitrocdn.com/yYolPsxkHfeoqRKSyqGQlaFpLZMHhVYI/assets/images/optimized/rev-23c0ecc/hpe-photos.s3.us-east-2.amazonaws.com/wp-content/uploads/20230307160059/Hennessey-H800-Mustang-GT-2.jpg",
		"https://cdn-hednb.nitrocdn.com/yYolPsxkHfeoqRKSyqGQlaFpLZMHhVYI/assets/images/optimized/rev-23c0ecc/hpe-photos.s3.us-east-2.amazonaws.com/wp-content/uploads/2025/02/Hennessey-2024-Ford-Mustang-GT-for-Sale-Black-0554-2.jpg",
		"https://cdn-hednb.nitrocdn.com/yYolPsxkHfeoqRKSyqGQlaFpLZMHhVYI/assets/images/optimized/rev-23c0ecc/hpe-photos.s3.us-east-2.amazonaws.com/wp-content/uploads/2025/02/Hennessey-2024-Ford-Mustang-GT-for-Sale-Black-0554-4.jpg",
		"https://cdn-hednb.nitrocdn.com/yYolPsxkHfeoqRKSyqGQlaFpLZMHhVYI/assets/images/optimized/rev-23c0ecc/hpe-photos.s3.us-east-2.amazonaws.com/wp-content/uploads/20230307160059/Hennessey-H800-Mustang-GT-2.jpg",
		}

	destDir := "./"

	err := ConcurrentFileDownloader(urls, destDir, 6)
	if err != nil {
		log.Fatal(err)
	}
}
