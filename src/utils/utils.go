package utils

import (
	"encoding/binary"
	"fmt"
	"net"
	"flag"
	"runtime"
)

func IpToUint32(ipStr string) (uint32, error) {
	ip := net.ParseIP(ipStr).To4()
	if ip == nil {
		return 0, fmt.Errorf("Invalid IP: %s", ipStr)
	}
	return binary.BigEndian.Uint32(ip), nil
}

func ParseCmd() (string, int, bool) {
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
