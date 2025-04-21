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
)

const totalIPs = 1 << 32

func main() {
    filePath, numWorkers, allowParallel := parseCmd()

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
        printMemUsage()
        return
    }

    fmt.Println("Parallel processing is not implemented yet.", numWorkers)
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

		ip, err := ipToUint32(ipStr)
		if err != nil {
			fmt.Println("Error transforming to Uint32:", err)
			continue
		}

        if bitmask[ip/8]&(1<<(ip%8)) == 0 {
            bitmask[ip/8] |= 1 << (ip % 8)
            iterator++
        }
	}

	if err := scanner.Err(); err != nil {
		return iterator, err
	}

	return iterator, nil
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

func printMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	fmt.Printf("\tSys = %v MiB", m.Sys/1024/1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}
