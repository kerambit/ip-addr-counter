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
	filePath, allowParallel := utils.ParseCmd()

	numWorkers := 1

	if allowParallel {
		maxWorkers := runtime.NumCPU()
		numWorkers = maxWorkers / 2
	}

	defer utils.PrintMemUsage()
	start := time.Now()

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error open the file:", err)
		return
	}
	defer file.Close()

	_, err = file.Stat()
	if err != nil {
		fmt.Println("Failed to get file info:", err)
		return
	}

	chunks, _ := utils.SplitFileIntoChunks(filePath, numWorkers)

	bitmask := make([]byte, totalIPs/8)
	var wg sync.WaitGroup

	for _, boundary := range chunks {
		wg.Add(1)
		go processFilePart(file, boundary.Start, boundary.End, bitmask, &wg)
	}

	wg.Wait()

	uniqueIPs := utils.CountUniqueIPs(bitmask)
	fmt.Println("Unique ips: ", uniqueIPs)
	fmt.Println("Total time:", time.Since(start))
}

func processFilePart(file *os.File, start int64, end int64, bitmask []byte, wg *sync.WaitGroup) {
	defer wg.Done()

	reader := io.NewSectionReader(file, start, end-start)
	bufReader := bufio.NewReader(reader)

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
