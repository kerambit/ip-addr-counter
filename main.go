package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"runtime"
	"time"
	"flag"
	"io"
	"strings"
	"sync"
)

const totalIPs = 1 << 32

func main() {
    filePath, numWorkers, allowParallel := parseCmd()

    defer printMemUsage()
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

	uniqueIPs := countUniqueIPs(bitmask)
    fmt.Println("Unique ips: ", uniqueIPs)
    fmt.Println("Total time:", time.Since(start))
}

func ipToUint32(ipStr string) (uint32, error) {
	ip := net.ParseIP(ipStr).To4()
	if ip == nil {
		return 0, fmt.Errorf("Invalid IP: %s", ipStr)
	}
	return binary.BigEndian.Uint32(ip), nil
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

		ipUint32, err := ipToUint32(ipStr)
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

		ipUint32, err := ipToUint32(ipStr)
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

func parseCmd() (string, int, bool) {
    var(
        filePath string
        numWorkers int
        allowParallel bool
    )

    flag.StringVar(&filePath, "path", "", "Path to the file")
    flag.IntVar(&numWorkers, "numWorkers", 2, "Max workers for parallel processing")
    flag.BoolVar(&allowParallel, "allowParallel", false, "Allow parallel processing")
    flag.Parse()

	fmt.Println("Path file is:", filePath)
	fmt.Println("numWorkers is:", numWorkers)
	fmt.Println("allowParallelMode is:", allowParallel)

	return filePath, numWorkers, allowParallel
}

func countUniqueIPs(bitmask []byte) uint32 {
	var count uint32 = 0
	for _, b := range bitmask {
		count += uint32(bitsOn(b))
	}
	return count
}

func bitsOn(b byte) int {
	var count int = 0
	for b > 0 {
		count += int(b & 1)
		b >>= 1
	}
	return count
}

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	fmt.Printf("\tSys = %v MiB", m.Sys/1024/1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}
