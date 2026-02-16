# go-concurrent-file-downloader
A Concurrent File Downloader in Go

A high-performance concurrent file downloader built in Go that supports:
Single file downloads
Sequential multi-file downloads
Controlled concurrent downloads
Structured result reporting (size, duration, errors)
This project demonstrates practical use of Go’s concurrency primitives for real-world I/O-bound workloads.

🧠 What This Project Demonstrates
- Goroutines for parallel execution
- Buffered channels as a semaphore (concurrency limiting)
- sync.WaitGroup for coordination
- Channel-based result aggregation
- Proper error handling and cleanup
- HTTP streaming with io.Copy
- Controlled Concurrency

To prevent resource exhaustion:
- limiter := make(chan struct{}, maxConcurrent)

Each download acquires a slot before starting and releases it when done.
This ensures the system never spawns unlimited goroutines.

🏗 Real-World Applications

This architecture mirrors systems such as:
- Web crawlers

- Media download managers

- Data ingestion pipelines

- Backup/sync services

0 CDN replication tools

▶ Example
urls := []string{"file1.jpg", "file2.jpg", "file3.jpg"}
ConcurrentFileDownloader(urls, "./downloads", 3)

🔧 Tech Used

Go • Goroutines • Channels • WaitGroup • HTTP • File I/O
