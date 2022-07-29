package common

import "fmt"


var ToC_Flag = map[uint32]string{
  (1 << 2): "kTocNewObjList",       // 2
  (1 << 3): "kTocRawData",          // 4
  (1 << 5): "kTocInterleavedData",  // 8
  (1 << 6): "kTocBigEndian",        // 32
  (1 << 7): "kTocDAQmxRawData",     // 64
  (1 << 1): "kTocMetaData",         // 128
}

// TODO - Set up structs to house segment info
// Something like this maybe?
// TDMS_File{[]Segment{Header, Meta, Raw}}
type ObjProperty struct{
  LenName uint32
  Name string
  Dtype uint32
  value any // will have to make a generic or interface

}

type ObjProperties []ObjProperty

type SegObj struct{
  LenObjPath uint32
  ObjName string
  RawIndex uint32
  NumOfProperties uint32
  Properties []ObjProperties

}

type SegObjs []SegObj

type Segment struct{
  numOfObj uint32
  Objs SegObjs
}

type Headers []Header

type Header struct{
  Loc int64
  Tag string
  Toc uint32
  Ver uint32
  Seg uint64
  Raw uint64
}

func (lI *Header) Show(){
  fmt.Printf("Loc: %v\n", lI.Loc)
  fmt.Printf("Tag: %s\n", lI.Tag)
  fmt.Printf("Toc: %v\n", lI.Toc)
  fmt.Printf("Ver: %v\n", lI.Ver)
  fmt.Printf("Seg: %v\n", lI.Seg)
  fmt.Printf("raw: %v\n", lI.Raw)
}

