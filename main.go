package main

import (
	"bufio"
	"fmt"
	"io"
	"ipcounter/utils"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

const totalIPs = 1 << 32

func main() {
	filePath, numWorkers, allowParallel := utils.ParseCmd()

	defer utils.PrintMemUsage()
	start := time.Now()

	if !allowParallel {
		uniqueIPs, err := processFileInSingle(filePath)
		if err != nil {
			fmt.Println("Error: ", err)
			fmt.Println("Total time:", time.Since(start))
			return
		}

		fmt.Println("Unique ips: ", uniqueIPs)
		fmt.Println("Total time:", time.Since(start))
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error open the file:", err)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		fmt.Println("Failed to get file info:", err)
		return
	}

	fileSize := info.Size()
	maxWorkers := runtime.NumCPU()
	if numWorkers > maxWorkers {
		fmt.Println("Number of workers exceeds the maximum available CPU cores. Setting to max workers.")
		numWorkers = maxWorkers
	}

	chunkSize := fileSize / int64(numWorkers)

	bitmask := make([]byte, totalIPs/8)
	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		startOffset := int64(i) * chunkSize
		endOffset := startOffset + chunkSize
		if i == numWorkers-1 {
			endOffset = fileSize
		}

		wg.Add(1)
		go processFilePart(file, startOffset, endOffset, bitmask, &wg)
	}

	wg.Wait()

	uniqueIPs := utils.CountUniqueIPs(bitmask)
	fmt.Println("Unique ips: ", uniqueIPs)
	fmt.Println("Total time:", time.Since(start))
}

func processFileInSingle(inputFile string) (uint32, error) {
	bitmask := make([]byte, totalIPs/8)

	var iterator uint32 = 0

	file, err := os.Open(inputFile)
	if err != nil {
		return iterator, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ipStr := scanner.Text()
		if len(ipStr) == 0 {
			continue
		}

		ipUint32, err := utils.IpToUint32(ipStr)
		if err != nil {
			fmt.Printf("Error transforming to Uint32: %s\n", ipStr)
			continue
		}

		index := ipUint32 / 8
		mask := byte(1 << (ipUint32 % 8))

		if bitmask[index]&mask == 0 {
			bitmask[index] |= mask
			iterator++
		}
	}

	if err := scanner.Err(); err != nil {
		return iterator, err
	}

	return iterator, nil
}

func processFilePart(file *os.File, start int64, end int64, bitmask []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	reader := io.NewSectionReader(file, start, end-start)
	bufReader := bufio.NewReader(reader)

	// move the reader to the start of the line
	if start != 0 {
		_, err := bufReader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Println("Error searching a new line:", err)
			return
		}
	}

	for {
		line, err := bufReader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading the line:", err)
			break
		}

		ipStr := strings.TrimSpace(line)
		if ipStr == "" {
			continue
		}

		ipUint32, err := utils.IpToUint32(ipStr)
		if err != nil {
			fmt.Printf("Error transforming to Uint32: %s\n", ipStr)
			continue
		}

		index := ipUint32 / 8
		mask := byte(1 << (ipUint32 % 8))

		if bitmask[index]&mask == 0 {
			bitmask[index] |= mask
		}
	}
}
