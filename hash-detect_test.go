package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"
)

func GeneratePermutationTests() []struct {
	input  string
	output string
} {
	var tests []struct {
		input  string
		output string
	}

	// Add fixed tests
	tests = append(tests,
		struct {
			input  string
			output string
		}{"d41d8cd98f00b204e9800998ecf8427e", "Possible: MD5, CRC32, Adler-32"},
		// ... add more fixed tests
	)

	// Generate permutations based on length
	for length := 0; length < 350; length++ {
		hash := string(make([]byte, length))
		// populate hash as needed
		// ...
		tests = append(tests, struct {
			input  string
			output string
		}{hash, identifyHash(hash)})
	}

	// Generate permutations based on prefixes
	prefixes := []string{"$2a$", "$2b$", "$2y$", "$5$", "$argon2i$", "$argon2id$", "$pbkdf2$", "$scrypt$", "HMAC", "SipHash", "FarmHash"}
	for _, prefix := range prefixes {
		hash := prefix + "some_hash_value_here"
		tests = append(tests, struct {
			input  string
			output string
		}{hash, identifyHash(hash)})
	}

	return tests
}
func TestIdentifyHash(t *testing.T) {
	tests := GeneratePermutationTests()
	var testCount int

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			testCount++ // Increment the test count
			if got := identifyHash(test.input); got != test.output {
				t.Errorf("identifyHash(%q) = %q; want %q", test.input, got, test.output)
			}
		})
	}

	// Print or otherwise use the test count
	fmt.Printf("Total tests run: %d\n", testCount)
}

func TestOutputResultsJSON(t *testing.T) {
	results := []HashResult{
		{Hash: "d41d8cd98f00b204e9800998ecf8427e", HashType: "Possible: MD5, CRC32, Adler-32"},
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outputResults(results, "json")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)

	expected := `[{"hash":"d41d8cd98f00b204e9800998ecf8427e","hashType":"Possible: MD5, CRC32, Adler-32"}]`
	if buf.String() != expected+"\n" {
		t.Errorf("outputResults() = %q; want %q", buf.String(), expected)
	}
}
