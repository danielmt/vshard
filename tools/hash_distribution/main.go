package main

import (
	"bufio"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"lab.decoded.io/daniel/vshard"
)

const numServers = 7

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
func GenerateKeyHash(key string) string {
	hash := md5.Sum([]byte(key))
	return hex.EncodeToString(hash[:])
}

func GetDist(data [numServers]int, evenDistribution int) {
	for n, result := range data {
		diff := 0
		plus := true

		if result > evenDistribution {
			diff = result - evenDistribution
		} else {
			diff = evenDistribution - result
			plus = false
		}

		percent := diff / evenDistribution * 100

		if !plus {
			diff = -diff
		}

		fmt.Printf("%d: %d (%d / %d%%)\n", n, result, diff, percent)
	}
}

func main() {
	lines, err := readLines("keys.txt")
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	numberOfLines := 0

	for _, line := range lines {
		numberOfLines++
		hexHash := GenerateKeyHash(line)
		serverDistributionMD5[vshard.ShardedServerStrategyMD5(hexHash, numServers)]++
		serverDistributionFarmhash[vshard.ShardedServerStrategyFarmhash(hexHash, numServers)]++
	}

	evenDistribution := numberOfLines / numServers

	fmt.Printf("results:\n\n")
	fmt.Printf("* MD5: %#v\n", serverDistributionMD5)

	GetDist(serverDistributionMD5, evenDistribution)

	fmt.Printf("\n* Farmhash: %#v\n", serverDistributionFarmhash)

	GetDist(serverDistributionFarmhash, evenDistribution)

	fmt.Printf("\neven distribution: %d\n\n", evenDistribution)

}
