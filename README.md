
### IpV4 Address Counter

App calculates the number of unique IPv4 addresses in a given file and print the result.

### Features:

- Processing large files
- Ready-to-work with simple CLI
- Supports goroutines for faster processing

### Usage

Build project into binary:

```bash
make build-mac
```

Supports two modes of operation:

1. Without goroutines

```bash
./bin/counter -path="path_to_file"
```

2. With goroutines

```bash
./bin/counter -path="path_to_file" -allowParallel=true
```

#### Note: by default, the app will use the first mode. To use the second mode, you need to specify the `-allowParallel` flag.

### Benchmarks

Tested on large [file](https://ecwid-vgv-storage.s3.eu-central-1.amazonaws.com/ip_addresses.zip).

Without goroutines:
```bash
Unique ips:  1000000000
Total time: 8m39.046350542s
Alloc = 577 MiB TotalAlloc = 122578 MiB Sys = 1075 MiB  NumGC = 241
```

With goroutines (5 workers):
```bash
Unique ips:  1000000000
Total time: 4m34.130946167s
Alloc = 556 MiB TotalAlloc = 122582 MiB Sys = 1079 MiB  NumGC = 241
```
