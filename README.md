# FastCheck - Website Health Monitor

![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

CLI tool to concurrently check website availability with SSL verification and latency measurement.

## Features

- ✅ Concurrent status checks (goroutines)
- 🌐 HTTP/HTTPS support with custom timeout
- 🔒 SSL certificate validation
- ⏱️ Response latency measurement
- 📊 JSON export for automation
- 🎨 Color-coded terminal output

## Installation

```bash
git clone https://github.com/seuuser/fastcheck
cd fastcheck
go build -o fastcheck
```

## Usage

Basic check:
```bash
./fastcheck --file urls.txt
```

Advanced options:
```bash
./fastcheck --file urls.txt --ssl --timeout 3 --output results.json
```

## Sample Output

## Tech Stack
- Concurrency: Goroutines + WaitGroup
- Networking: net/http, crypto/tls
- CLI: flag package
- Data: JSON encoding

