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
	group := make([]byte, 0)
	searchingByte := data[0]
	count := 0
	groupCount := 0
	for _, current := range data {
		if searchingByte == current {
			count++
		} else {
			if count < 4 {
				if groupCount+count >= 127 {
					result = append(result, sa.createControlByte(groupCount, false))
					groupCount = 0
					result = append(result, group...)
					group = make([]byte, 0)
				}
				groupCount += count
				group = append(group, searchingByte)
				count = 1
				searchingByte = current
			} else {
				if len(group) != 0 {
					result = append(result, sa.createControlByte(groupCount, false))
					groupCount = 0
					result = append(result, group...)
					group = make([]byte, 0)
				}
				result = append(result, sa.createControlByte(count, true), searchingByte)
				count = 1
				searchingByte = current
			}
		}
	}
	if len(group) != 0 {
		if count < 4 {
			groupCount += count
			group = append(group, searchingByte)
			count = 0
		}
		result = append(result, sa.createControlByte(groupCount, false))
		result = append(result, group...)
	}
	if count != 0 {
		result = append(result, sa.createControlByte(count, true), searchingByte)
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
