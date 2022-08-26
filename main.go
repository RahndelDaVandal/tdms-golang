package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"tdms/common"
	"tdms/reader"
)

// KNOWN NI DTYPES
// 07 00 00 00: tdsTypeU32
// 20 00 00 00: tdsTypeString
// 03 00 00 00: tdsTypeI32

func main() {
	fileName := "./test_files/2020-09-17T22-45-47_.tdms"

	r := reader.BytesReader{}
	r.LoadFile(fileName)

	S := common.Segments{}
	s := common.Segment{}
	var err error

	for err == nil {
		s, err = r.ReadSegment()
		S = append(S, s)
		r.Move(int64(s.Header.Seg))
	}

	fmt.Println("Number of Segments: ", len(S))

	for i := 0; i < 11; i++ {
		fmt.Println(S[i])
	}

	defer r.File.Close()
}

// ---------------------------------------------------------------------------

func loadFile(fileName string) *os.File {
	file, openErr := os.Open(fileName)
	if openErr != nil {
		log.Fatalf("Failed to open %v", fileName)
	}
	return file
}

type Window struct {
	start int
	end   int
}

func (w *Window) Move(size int) {
	w.start = w.end
	w.end += size
}

func decodeBytes[T any](data []byte, field *T) {
	buffer := bytes.NewBuffer(data)
	binary.Read(buffer, binary.LittleEndian, field)
}

func decodeBytesWithBuf[T any](buffer *bytes.Buffer, field *T) {
	binary.Read(buffer, binary.LittleEndian, field)
}

func readNextBytes(file *os.File, number int) ([]byte, int, error) {
	bytes := make([]byte, number)

	bytesRead, err := file.Read(bytes)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}

	return bytes, bytesRead, err
}

func getSegmentHeaders(fileName string, headers *common.Headers) {
	f := loadFile(fileName)
	fileStats, _ := f.Stat()
	fileSize := fileStats.Size()
	fmt.Printf("fileSize: %v\n", fileSize)

	h := common.Header{}
	var loc int64
	var data []byte
	var err error

	// for loc < fileSize{
	for {
		loc, _ = f.Seek(int64(h.Seg), 1)
		h.Loc = loc
		data, _, err = readNextBytes(f, 4)
		if err == io.EOF {
			break
		}
		h.Tag = string(data)
		data, _, err = readNextBytes(f, 24)
		if err == io.EOF {
			break
		}
		buffer := bytes.NewBuffer(data)
		decodeBytesWithBuf(buffer, &h.Toc)
		decodeBytesWithBuf(buffer, &h.Ver)
		decodeBytesWithBuf(buffer, &h.Seg)
		decodeBytesWithBuf(buffer, &h.Raw)
		*headers = append(*headers, h)
	}
	fmt.Printf("Num of Segment Headers: %v\n", len(*headers))

	defer f.Close()
}
