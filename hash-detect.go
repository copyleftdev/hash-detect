package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

type HashType struct {
	Length  int
	Prefix  string
	Suffix  string
	Message string
}

type HashResult struct {
	Hash     string `json:"hash" xml:"hash"`
	HashType string `json:"hashType" xml:"hashType"`
}

func identifyHashByLength(hash string, hashTypes map[int]string) string {
	length := len(hash)
	if msg, exists := hashTypes[length]; exists {
		return msg
	}
	return ""
}

func identifyHashByPattern(hash string, patternTypes []HashType) string {
	for _, pt := range patternTypes {
		if strings.HasPrefix(hash, pt.Prefix) && (pt.Suffix == "" || strings.HasSuffix(hash, pt.Suffix)) {
			if pt.Length == 0 || pt.Length == len(hash) {
				return pt.Message
			}
		}
	}
	return ""
}

func identifyHash(hash string) string {
	hashTypes := map[int]string{
		32:  "Possible: MD5, CRC32, Adler-32",
		40:  "Possible: SHA-1, RIPEMD",
		56:  "Possible: SHA-224",
		64:  "Possible: SHA-256, BLAKE2, MurmurHash, CityHash, xxHash",
		96:  "Possible: SHA-384",
		128: "Possible: SHA-512, Whirlpool",
		160: "Possible: RIPEMD-160",
		320: "Possible: RIPEMD-320",
	}

	patternTypes := []HashType{
		{34, "$1$", "", "Possible: MD5 Crypt"},
		{60, "$2a$", "", "Possible: bcrypt (Blowfish Crypt)"},
		{60, "$2b$", "", "Possible: bcrypt (Blowfish Crypt)"},
		{60, "$2y$", "", "Possible: bcrypt (Blowfish Crypt)"},
		{62, "$5$", "", "Possible: SHA-256 Crypt"},
		{96, "$argon2i$", "", "Possible: Argon2i"},
		{96, "$argon2id$", "", "Possible: Argon2id"},
		{0, "$pbkdf2$", "", "Possible: PBKDF2"},
		{0, "$scrypt$", "", "Possible: scrypt"},
		{0, "HMAC", "", "Possible: HMAC"},
		{0, "SipHash", "", "Possible: SipHash"},
		{0, "FarmHash", "", "Possible: FarmHash"},
	}

	lengthMatch := identifyHashByLength(hash, hashTypes)
	if lengthMatch != "" {
		return lengthMatch
	}

	patternMatch := identifyHashByPattern(hash, patternTypes)
	if patternMatch != "" {
		return patternMatch
	}

	return "Unknown hash type"
}

func outputResults(results []HashResult, format string) {
	switch format {
	case "json":
		jsonData, err := json.Marshal(results)
		if err != nil {
			fmt.Println("Error generating JSON:", err)
			return
		}
		fmt.Println(string(jsonData))
	case "xml":
		xmlData, err := xml.MarshalIndent(results, "", "  ")
		if err != nil {
			fmt.Println("Error generating XML:", err)
			return
		}
		fmt.Printf("%s%s\n", xml.Header, xmlData)
	case "text":
		for _, result := range results {
			fmt.Printf("%s: %s\n", result.Hash, result.HashType)
		}
	case "csv":
		writer := csv.NewWriter(os.Stdout)
		_ = writer.Write([]string{"Hash", "HashType"}) // Write header
		for _, result := range results {
			_ = writer.Write([]string{result.Hash, result.HashType})
		}
		writer.Flush()
	default:
		fmt.Println("Unsupported format. Use 'json', 'xml', 'text', or 'csv'.")
	}
}

func processHashes(scanner *bufio.Scanner, results *[]HashResult) {
	for scanner.Scan() {
		hashStr := scanner.Text()
		hashType := identifyHash(hashStr)
		*results = append(*results, HashResult{Hash: hashStr, HashType: hashType})
	}
}
func readHashesFromFile(filename string, results *[]HashResult) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	processHashes(scanner, results)
	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  hash-detect <hash>                : Detect a single hash")
		fmt.Println("  hash-detect -f <filename>         : Detect hashes from a text file")
		fmt.Println("  hash-detect -f <filename> -o <format> : Output in a specific format ('json', 'xml', 'text', 'csv')")
		os.Exit(1)
	}

	var results []HashResult // Slice to collect hash identification data
	format := "text"         // Default output format

	for i, arg := range os.Args[1:] {
		switch arg {
		case "-f":
			if len(os.Args) <= i+2 {
				fmt.Println("Please specify a filename.")
				os.Exit(1)
			}

			if err := readHashesFromFile(os.Args[i+2], &results); err != nil {
				fmt.Printf("Could not open file: %s\n", err)
				os.Exit(1)
			}

		case "-o":
			if len(os.Args) <= i+2 {
				fmt.Println("Please specify an output format.")
				os.Exit(1)
			}
			format = os.Args[i+2]

		default:
			hashType := identifyHash(arg)
			results = append(results, HashResult{Hash: arg, HashType: hashType})
		}
	}

	outputResults(results, format)
}
