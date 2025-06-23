package main

import (
	"bufio"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	colorGreen  = "\033[32m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorReset  = "\033[0m"
)

type Result struct {
	URL      string `json:"url"`
	Status   string `json:"status"`
	Code     int    `json:"code,omitempty"`
	Latency  string `json:"latency,omitempty"`
	SSLValid bool   `json:"ssl_valid,omitempty"`
}

func checkSSL(url string) bool {
	conn, err := tls.Dial("tcp", url+":443", &tls.Config{
		ServerName:         url,   // Verifica se o certificado pertence a este domínio
		InsecureSkipVerify: false, // (padrão) → ativa verificação
		// InsecureSkipVerify: true, // (dev/test) → desativa verificação
	})
	if err != nil {
		return false
	}
	defer conn.Close()

	cert := conn.ConnectionState().PeerCertificates[0]
	err = cert.VerifyHostname(url)
	return err == nil
}

func checkURL(url string, timeout int, checkSSLFlag bool, wg *sync.WaitGroup, results chan<- Result) {
	defer wg.Done()

	start := time.Now()
	client := http.Client{
		Timeout: time.Duration(timeout) * time.Second,
	}

	resp, err := client.Get(url)
	result := Result{URL: url}

	if err != nil {
		result.Status = "DOWN"
	} else {
		defer resp.Body.Close()
		result.Status = "UP"
		result.Code = resp.StatusCode
		result.Latency = time.Since(start).Round(time.Millisecond).String()

		if checkSSLFlag {
			result.SSLValid = checkSSL(url)
		}
	}

	results <- result
}

func main() {
	filePath := flag.String("file", "urls.txt", "Arquivo com URLs")
	outputFile := flag.String("output", "", "Salvar resultados em JSON")
	timeout := flag.Int("timeout", 5, "Timeout em segundos")
	checkSSL := flag.Bool("ssl", false, "Verificar validade SSL")
	flag.Parse()

	file, err := os.Open(*filePath)
	if err != nil {
		fmt.Printf("%sErro ao abrir arquivo:%s %v\n", colorRed, colorReset, err)
		os.Exit(1)
	}
	defer file.Close()

	var wg sync.WaitGroup
	results := make(chan Result)
	var allResults []Result

	go func() {
		for res := range results {
			allResults = append(allResults, res)
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()
		if url == "" {
			continue
		}
		wg.Add(1)
		go checkURL(url, *timeout, *checkSSL, &wg, results)
	}

	wg.Wait()
	close(results)

	for _, res := range allResults {
		color := colorGreen
		if res.Status == "DOWN" {
			color = colorRed
		}

		sslInfo := ""
		if *checkSSL {
			sslStatus := "✓"
			if !res.SSLValid {
				sslStatus = "✗"
			}
			sslInfo = fmt.Sprintf(" SSL: %s%s%s", colorYellow, sslStatus, colorReset)
		}

		fmt.Printf("%s[%s]%s %s (Latency: %s)%s\n",
			color, res.Status, colorReset,
			res.URL, res.Latency, sslInfo)
	}

	if *outputFile != "" {
		jsonData, _ := json.MarshalIndent(allResults, "", "  ")
		err := os.WriteFile(*outputFile, jsonData, 0644)
		if err != nil {
			fmt.Printf("%sErro ao salvar arquivo:%s %v\n", colorRed, colorReset, err)
		} else {
			fmt.Printf("\nResultados salvos em: %s\n", *outputFile)
		}
	}
}
