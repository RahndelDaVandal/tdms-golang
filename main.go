package main

import (
	// "bufio"
	// "encoding/binary"
	"fmt"
	"log"
	"os"
	"reflect"
)

var ToC_Flag = map[int]string{
  (1 << 1): "kTocMetaData",
  (1 << 3): "kTocRawData",
  (1 << 7): "kTocDAQmxRawData",
  (1 << 5): "kTocInterleavedData",
  (1 << 6): "kTocBigEndian",
  (1 << 2): "kTocNewObjList",
}

type LeadIn struct{
  tag string
  ToC string
  version uint32
  metaDataLoc uint64
  rawDataLoc uint64
}

func main() {
  fileName := "./test_files/2020-09-17T22-45-47_.tdms"
  fmt.Println(reflect.TypeOf(readLoc(fileName, int64(4))))
}

func loadFile(fileName string) *os.File{
  file, openErr := os.Open(fileName)
  if openErr != nil{
    log.Fatalf("Failed to open %v", fileName)
  }
  return file
}

func readLoc(fileName string, loc int64) []byte {
  f := loadFile(fileName)
  f.Seek(loc, 0)
  buf := make([]byte, 8)
  rLen, err := f.ReadAt(buf, int64(loc))
    if err != nil{
      log.Fatalf("READ ERROR: %v", err)
    }
  return buf[:rLen]
}

func readLeadIn(fileName string) {
  f := loadFile(fileName)
  bufSize := 4
  buf := make([]byte, bufSize)
  for i := 0; i < 3; i++{
    readTotal, err := f.Read(buf)
    if err != nil{
      log.Fatalf("READ ERROR: %v", err)
    }
    fmt.Printf("%x  %v\n", buf[:readTotal], string(buf[:readTotal]))
  }
  bufSize = 8
  buf = make([]byte, bufSize)
  for i := 0; i < 2; i++{
    readTotal, err := f.Read(buf)
    if err != nil{
      log.Fatalf("READ ERROR: %v", err)
    }
    fmt.Printf("%x  %v\n", buf[:readTotal], string(buf[:readTotal]))
  }
  defer f.Close()

}

func readByChunk(fileName string, chunkSize int, numOfChunks int){
  file, openErr := os.Open(fileName)
  if openErr != nil{
    log.Fatalf("Failed to open %v", fileName)
  }
  
  defer file.Close()
  buf := make([]byte, chunkSize)
  for i := 0; i < numOfChunks; i++{
    readTotal, readErr := file.Read(buf)
    if readErr != nil{
      log.Fatalf("READ ERROR: %v", readErr)
    }
    fmt.Printf("%x  %v\n", buf[:readTotal], string(buf[:readTotal]))
  }
}

// func readByScanner(fileName string){
//   file, openErr := os.Open(fileName)
//   if openErr != nil{
//     log.Fatalf("Failed to open %v", fileName)
//   }
//   defer file.Close()
//
//   scanner := bufio.NewScanner(file)
//   
//   
// }
