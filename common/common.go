package common

import "fmt"

type Property struct{
  LenPropName uint32
  Name string
  Dtype uint32
  LenValStr uint32
  Value any
}
type Obj struct{
  LenObjPath uint32
  Name string
  RawDataIndex uint32
  NumOfProperties uint32
  Properties []Property
}
type Segment struct{
  Loc int
  Header Header
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

var ToC_Flag = map[uint32]string{
  (1 << 2): "kTocNewObjList",       // 2
  (1 << 3): "kTocRawData",          // 4
  (1 << 5): "kTocInterleavedData",  // 8
  (1 << 6): "kTocBigEndian",        // 32
  (1 << 7): "kTocDAQmxRawData",     // 64
  (1 << 1): "kTocMetaData",         // 128
}
