package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"
	"os"
)

const(one = uint32(1))

var ToC_Flag = map[uint32]string{
  (1 << 2): "kTocNewObjList",       // 2
  (1 << 3): "kTocRawData",          // 4
  (1 << 5): "kTocInterleavedData",  // 8
  (1 << 6): "kTocBigEndian",        // 32
  (1 << 7): "kTocDAQmxRawData",     // 64
  (1 << 1): "kTocMetaData",         // 128
}

type Header struct{
  tag string
  toc uint32
  ver uint32
  seg uint64
  raw uint64
}

func (lI *Header) show(){
  fmt.Printf("tag: %s\n", lI.tag)
  fmt.Printf("toc: %v\n", lI.toc)
  fmt.Printf("ver: %v\n", lI.ver)
  fmt.Printf("seg: %v\n", lI.seg)
  fmt.Printf("raw: %v\n", lI.raw)
}

func main(){
  fileName := "./test_files/2020-09-17T22-45-47_.tdms"
  f := loadFile(fileName)

  H := Header{}
  var loc int64

  fmt.Println("Starting Loc: ", loc)
  br := readSegmentHeader(f, &H)
  H.show()
  fmt.Println("Header End: ", int64(br)+loc)
  loc, _ = f.Seek(int64(H.seg), 1)
  fmt.Println("Starting Loc: ", loc)
  br = readSegmentHeader(f, &H)
  H.show()
  fmt.Println("Header End: ", int64(br)+loc)
  loc, _ = f.Seek(int64(H.seg), 1)
  fmt.Println("Starting Loc: ", loc)
  br = readSegmentHeader(f, &H)
  H.show()
  fmt.Println("Header End: ", int64(br)+loc)
  loc, _ = f.Seek(int64(H.seg), 1)
  fmt.Println("Starting Loc: ", loc)
  br = readSegmentHeader(f, &H)
  H.show()
  fmt.Println("Header End: ", int64(br)+loc)
  loc, _ = f.Seek(int64(H.seg), 1)

  buf, _ := readNextBytes(f, 512)
  fmt.Printf("\n[%v]", string(buf))
  
  defer f.Close()
}

func showMultipleHeaders(f *os.File, headerNum int){
  header := Header{}
  var loc int64
  var bytesRead int
  var end int64

  for i := 0; i < headerNum; i++{
    loc, _ = f.Seek(int64(header.seg), 1)
    fmt.Println("-------------------------------")
    fmt.Println("Segment Start: ", loc)
    fmt.Println("-------------------------------")
    bytesRead = readSegmentHeader(f, &header)
    header.show()
    end = loc + int64(bytesRead)
    fmt.Println("-------------------------------")
    fmt.Println("Header End: ", end)
    fmt.Println("-------------------------------")
  }
}

func readSegmentHeader(f *os.File, h *Header) int {
  data, bytesRead := readNextBytes(f, 4)
  h.tag = string(data)
  data, bytesRead = readNextBytes(f, 24)
  buffer := bytes.NewBuffer(data)
  decodeBytes(buffer, &h.toc)
  decodeBytes(buffer, &h.ver)
  decodeBytes(buffer, &h.seg)
  decodeBytes(buffer, &h.raw)
  return bytesRead
}

func decodeBytes[T any](buffer *bytes.Buffer, field *T){
  binary.Read(buffer, binary.LittleEndian, field)
}

func readNextBytes(file *os.File, number int) ([]byte, int) {
	bytes := make([]byte, number)

	bytesRead, err := file.Read(bytes)
	if err != nil {
		log.Fatal(err)
	}

	return bytes, bytesRead
}

func loadFile(fileName string) *os.File{
  file, openErr := os.Open(fileName)
  if openErr != nil{
    log.Fatalf("Failed to open %v", fileName)
  }
  return file
}

// func readLeadIn(f *os.File, lI *LeadIn){
//   buf := make([]byte, 28)
//
//   _, err := f.Read(buf)
//   if err != nil{
//     fmt.Println(err)
//   }
//
//   br := bytes.NewBuffer(buf)
//   
//   fmt.Printf("%v\n", string(buf[:4]))
//   binary.Read(br, binary.LittleEndian, &lI.tag)
//   binary.Read(br, binary.LittleEndian, &lI.toc)
//   binary.Read(br, binary.LittleEndian, &lI.ver)
//   binary.Read(br, binary.LittleEndian, &lI.seg)
//   binary.Read(br, binary.LittleEndian, &lI.raw)
//
// }

