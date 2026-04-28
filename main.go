package main

import (
	"bufio"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"io"
	"os"
	"path/filepath"
)

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
func (sa *SimpleArchiver) countRepeating(data []byte, i int) int {
	b := data[i]
	count := 0
	for i < len(data) && data[i] == b && i < 127 {
		count++
		i++
	}
	return count
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

func (sa *SimpleArchiver) compress(data []byte) []byte {
	if len(data) == 0 {
		return sa.compressEmpty(data)
	}
	result := make([]byte, 0)

	for i := 0; i < len(data); {
		run := sa.countRepeating(data, i)
		if run >= 3 {
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
			run = sa.countRepeating(data, i)
			if run >= 3 {
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

func (sa *SimpleArchiver) decompress(data []byte) []byte {
	if len(data) == 0 {
		return []byte{}
	}
	result := make([]byte, 0)
	i := 0
	var control byte
	for i < len(data) {
		control = data[i]
		length := int(control & 127)
		if control&128 != 0 {
			i++
			for k := 0; k < length; k++ {
				result = append(result, data[i])
			}
			i++
		} else {
			i++
			result = append(result, data[i:i+length]...)
			i += length
		}
	}
	return result
}

func (sa *SimpleArchiver) CompressFile(inputPath, outputPath string) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("Error opening input file: %v", err)
	}
	defer inputFile.Close()
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("Error creating output file: %v", err)
	}
	defer outputFile.Close()
	writer := bufio.NewWriter(outputFile)
	defer writer.Flush()
	inputFileName := filepath.Base(inputPath)
	err = writer.WriteByte(byte(len(inputFileName)))
	if err != nil {
		return fmt.Errorf("Error writing input file: %v", err)
	}
	_, err = writer.Write([]byte(inputFileName))
	if err != nil {
		return fmt.Errorf("Error writing input file: %v", err)
	}
	reader := bufio.NewReader(inputFile)
	for {
		dataLen, readErr := reader.Read(sa.buffer)
		if readErr != nil && readErr != io.EOF {
			return fmt.Errorf("Error reading input file: %v", readErr)
		}
		compressedData := sa.compress(sa.buffer[:dataLen])
		blockSize := uint16(len(compressedData))
		err = writer.WriteByte(byte(blockSize >> 8))
		if err != nil {
			return fmt.Errorf("Error writing block size: %v", err)
		}
		err = writer.WriteByte(byte(blockSize))
		if err != nil {
			return fmt.Errorf("Error writing block size: %v", err)
		}
		_, err = writer.Write(compressedData)
		if err != nil {
			return fmt.Errorf("Error writing compressed data: %v", err)
		}
		if readErr == io.EOF {
			break
		}
	}
	return nil
}

func (sa *SimpleArchiver) DecompressFile(inputPath, outputDir string) error {
	archive, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("Error opening input file: %v", err)
	}
	defer archive.Close()
	reader := bufio.NewReader(archive)
	length, err := reader.ReadByte()
	if err != nil {
		return fmt.Errorf("Error reading input file: %v", err)
	}
	buffer := make([]byte, length)
	_, err = reader.Read(buffer)
	if err != nil {
		return fmt.Errorf("Error reading input file: %v", err)
	}
	outputFile, err := os.Create(filepath.Join(outputDir, string(buffer)))
	if err != nil {
		return fmt.Errorf("Error creating output file: %v", err)
	}
	defer outputFile.Close()
	for {
		blockSizeBuf := make([]byte, 2)
		_, err := reader.Read(blockSizeBuf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("Error reading input file: %v", err)
		}
		blockSize := uint16(blockSizeBuf[0])<<8 + uint16(blockSizeBuf[1])
		data := make([]byte, blockSize)
		_, readErr := reader.Read(data)
		if readErr != nil && readErr != io.EOF {
			return fmt.Errorf("Error reading input file: %v", readErr)
		}
		decompressedData := sa.decompress(data[:blockSize])
		_, err = outputFile.Write(decompressedData)
		if err != nil {
			return fmt.Errorf("Error writing decompressed data: %v", err)
		}
		if readErr == io.EOF {
			break
		}
	}
	return nil
}

func main() {
	tea.NewProgram(initialModel()).Run()
}
