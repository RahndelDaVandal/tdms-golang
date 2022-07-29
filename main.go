package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"tdms/common"
)

// 07 00 00 00: tdsTypeU32
// 20 00 00 00: tdsTypeString
// 03 00 00 00: tdsTypeI32

type Property struct{
  lenPropName uint32
  Name string
  Dtype uint32
  LenValStr uint32
  Value any
}
type Obj struct{
  lenObjPath uint32
  Name string
  RawDataIndex uint32
  NumOfProperties uint32
  Properties []Property
}
type Segment struct{
  Loc int
  Header common.Header
  NumOfObj uint32
  Objs []Obj
}
type Segments []Segment

func (s *Segments)Show(){
  for _, v := range *s{
    fmt.Printf("Segment Location: %v\n", v.Loc)
    fmt.Printf("# of Objects: %v\n", v.NumOfObj)
    for _, o := range v.Objs{
      fmt.Printf("\tName: %v\n", o.Name)
      fmt.Printf("\t\t# of Properties: %v\n", o.NumOfProperties)
      for _, p := range o.Properties{
        fmt.Printf("\t\t\t%v: %v\n", p.Name, p.Value)
      }
    } 
  }
}

func (r *BytesReader) readSegment()Segment{
  s := Segment{}
  s.Loc = r.loc
  r.readNext(4)
  s.Header.Tag = r.str
  r.readNext(4)
  s.Header.Toc = r.i32
  r.readNext(4)
  s.Header.Ver = r.i32
  r.readNext(8)
  s.Header.Seg = r.i64
  r.readNext(8)
  s.Header.Raw = r.i64
  // r.readNext(4)
  // s.NumOfObj = r.i32
  return s
}

func (r *BytesReader) readObj() Obj{
  o := Obj{}
  r.readNext(4)
  o.lenObjPath = r.i32
  r.readNext(int(o.lenObjPath))
  o.Name = r.str
  r.readNext(4)
  o.RawDataIndex = r.i32
  r.readNext(4)
  o.NumOfProperties = r.i32
  return o
}

func (r *BytesReader) readProperty() Property{
  p := Property{}
  r.readNext(4)
  p.lenPropName = r.i32
  r.readNext(int(p.lenPropName))
  p.Name = r.str
  r.readNext(4)
  p.Dtype = r.i32

  if p.Dtype == 32{
    r.readNext(4)
    p.LenValStr = r.i32
    r.readNext(int(p.LenValStr))
    p.Value = r.str
  } else {
    r.readNext(4)
    p.Value = r.i32
  }
  return p
}

func main(){
  fileName := "./test_files/2020-09-17T22-45-47_.tdms"

  r := BytesReader{}
  r.file = loadFile(fileName)


  h := common.Headers{}
  getSegmentHeaders(fileName, &h)

  for i:=0;i<5;i++{
    fmt.Println(h[i])
  }


  defer r.file.Close()
}

type BytesReader struct{
  file *os.File
  loc int
  bytes []byte
  i32 uint32
  i64 uint64
  str string
}

func (b *BytesReader) readNext(numOfBytes int){
  read := func() int {
    bytesRead, err := b.file.Read(b.bytes)
    if err != nil && err != io.EOF{
      log.Fatal(err)
    }
    return bytesRead
  }

  if numOfBytes == 4{
    b.i64 = 0
    b.bytes = make([]byte, numOfBytes)
    bytesRead := read()
    b.str = string(b.bytes)
    buf := bytes.NewBuffer(b.bytes)
    binary.Read(buf, binary.LittleEndian, &b.i32)
    b.loc += bytesRead
  } else if numOfBytes == 8 {
    b.i32 = 0
    b.bytes = make([]byte, numOfBytes)
    bytesRead := read()
    b.str = string(b.bytes)
    buf := bytes.NewBuffer(b.bytes)
    binary.Read(buf, binary.LittleEndian, &b.i64)
    b.loc += bytesRead
  } else{
    b.bytes = make([]byte, numOfBytes)
    bytesRead := read()
    b.str = string(b.bytes)
    b.loc += bytesRead
  }
}

func (b *BytesReader) printNext(numOfBytes int){
  b.readNext(numOfBytes)
  fmt.Printf("%v|%X|%v|%v|%v|\n",b.loc, b.bytes, b.i32, b.i64, b.str)
}

func loadFile(fileName string) *os.File{
  file, openErr := os.Open(fileName)
  if openErr != nil{
    log.Fatalf("Failed to open %v", fileName)
  }
  return file
}
// ---------------------------------------------------------------------------
type Window struct{
  start int
  end int
}

func (w *Window) Move(size int){
  w.start = w.end
  w.end += size
}

func decodeBytes[T any](data []byte, field *T){
  buffer := bytes.NewBuffer(data)
  binary.Read(buffer, binary.LittleEndian, field)
}

func decodeBytesWithBuf[T any](buffer *bytes.Buffer, field *T){
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

func getSegmentHeaders(fileName string, headers *common.Headers){
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
    if err == io.EOF{break}
    h.Tag = string(data)
    data, _, err = readNextBytes(f, 24)
    if err == io.EOF{break}
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
