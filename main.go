package main

import "fmt"

type SimpleArchiver struct {
	inputPath  string
	outputPath string
	buffer     []byte
}

func NewArchiver(inputPath string) *SimpleArchiver {
	buffer := make([]byte, 1024*8)
	return &SimpleArchiver{inputPath: inputPath, buffer: buffer}
}

func (sa *SimpleArchiver) compressEmpty(data []byte) []byte {
	if len(data) == 0 {
		return []byte{}
	}
	return data
}
func (sa *SimpleArchiver) countRepeating(data []byte) []byte {
	if len(data) == 0 {
		return []byte{}
	}
	result := make([]byte, 0)
	countRun := func(data []byte, i int) int {
		b := data[i]
		count := 0
		for i < len(data) && data[i] == b {
			count++
			i++
		}
		return count
	}

	for i := 0; i < len(data); {
		run := countRun(data, i)
		if run >= 4 {
			for run > 127 {
				result = append(result, sa.createControlByte(127, true), data[i])
				run -= 127
			}
			result = append(result, sa.createControlByte(run, true), data[i])
			i += run
			continue
		}

		start := i
		length := 0

		for i < len(data) && length < 127 {
			run = countRun(data, i)
			if run >= 4 {
				break
			}
			i++
			length++
		}
		result = append(result, sa.createControlByte(length, false))
		result = append(result, data[start:start+length]...)
	}
	return result
}

func (sa *SimpleArchiver) createControlByte(count int, isCompressed bool) byte {
	if count > 127 {
		count = 127
	}
	if isCompressed {
		return byte(count + 128)
	} else {
		return byte(count)
	}
}

func main() {
	archiver := NewArchiver("test.txt")
	fmt.Println(archiver.countRepeating([]byte("AAAAABCD")))
}
