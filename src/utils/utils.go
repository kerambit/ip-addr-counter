package utils

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"
)

type Boundaries struct {
	Start int64
	End   int64
}

func SplitFileIntoChunks(filePath string, parts int) ([]Boundaries, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get file size
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := fileInfo.Size()

	chunkSize := fileSize / int64(parts)
	boundaries := make([]Boundaries, 0, parts)

	var start int64 = 0

	for i := 0; i < parts; i++ {
		var end int64
		if i == parts-1 {
			end = fileSize
		} else {
			end = start + chunkSize
		}

		if end < fileSize {
			// Seek to the end position
			_, err := file.Seek(end, 0)
			if err != nil {
				return nil, err
			}

			// Read bytes until we find a newline or reach end of file
			buffer := make([]byte, 1)
			for {
				// Get current position
				currentPos, err := file.Seek(0, 1)
				if err != nil {
					return nil, err
				}

				if currentPos >= fileSize {
					break
				}

				n, err := file.Read(buffer)
				if err != nil || n == 0 {
					break
				}

				if buffer[0] == '\n' {
					end++
					break
				}
				end++
			}
		}

		boundaries = append(boundaries, Boundaries{Start: start, End: end})
		start = end
	}

	return boundaries, nil
}

func IpToUint32(ipStr string) (uint32, error) {
	ip := net.ParseIP(ipStr).To4()
	if ip == nil {
		return 0, fmt.Errorf("invalid IP: %s", ipStr)
	}
	return binary.BigEndian.Uint32(ip), nil
}

func ParseCmd() (string, int, bool) {
	var (
		filePath      string
		numWorkers    int
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

func CountUniqueIPs(bitmask []byte) uint32 {
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

func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", m.Alloc/1024/1024)
	fmt.Printf("\tTotalAlloc = %v MiB", m.TotalAlloc/1024/1024)
	fmt.Printf("\tSys = %v MiB", m.Sys/1024/1024)
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}
