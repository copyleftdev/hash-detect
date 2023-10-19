# Hash-Detect

Hash-Detect is a command-line tool for identifying hash types based on the input string. It checks hashes against various characteristics such as length, prefix, and suffix to make an educated guess on what the possible hash type could be.

## Table of Contents

1. [Installation](#installation)
2. [Usage](#usage)
3. [Supported Formats](#supported-formats)
4. [Limitations and Future Work](#limitations-and-future-work)
5. [Contribution](#contribution)

## Installation

Clone the repository to your local machine and build the project.

```bash
git clone https://github.com/copyleftdev/hash-detect.git
cd hash-detect
go build
```
## Usage

### Single Hash Detection

To identify a single hash, simply run the tool with the hash as an argument.

```bash
hash-detect <hash>
```

### Multiple Hashes from File

To identify multiple hashes listed in a text file, use the `-f` flag followed by the filename.

```bash
hash-detect -f <filename>
```

### Output Format

By default, the tool outputs in plain text. You can specify different formats like JSON, XML, or CSV by using the `-o` flag.

```bash
hash-detect -f <filename> -o <format>
```

Supported formats are: `json`, `xml`, `text`, `csv`

## Supported Formats

Currently, the tool can detect the following hash types based on their length, prefix, or suffix:

- MD5, CRC32, Adler-32
- SHA-1, RIPEMD
- SHA-224
- SHA-256, BLAKE2, MurmurHash, CityHash, xxHash
- SHA-384
- SHA-512, Whirlpool
- RIPEMD-160, RIPEMD-320
- And more...

## Limitations and Future Work

While the tool aims to be as accurate as possible, it's worth noting that:

- It makes an educated guess and should not be considered 100% reliable.
- Length-based detection can be ambiguous for certain hash lengths that multiple algorithms share.
- The pattern-based matching is still being improved for better accuracy.

## Contribution

Feel free to contribute to this project by submitting pull requests or opening issues.

