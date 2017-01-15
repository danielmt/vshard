package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"lab.decoded.io/daniel/vshard"
)

const numServers = 5

var (
	serverDistributionMD5      [numServers]int
	serverDistributionFarmhash [numServers]int
)

// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	lines, err := readLines("keys.txt")
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}
	for _, line := range lines {
		serverDistributionMD5[vshard.ShardedServerStrategyMD5(line, numServers)]++
		serverDistributionFarmhash[vshard.ShardedServerStrategyFarmhash(line, numServers)]++
	}

	fmt.Printf("MD5: %#v\n", serverDistributionMD5)
	fmt.Printf("Farmhash: %#v\n", serverDistributionFarmhash)
}
