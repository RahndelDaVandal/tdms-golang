package reader

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"tdms/common"
)

type BytesReader struct {
	File  *os.File
	Loc   int
	Bytes []byte
	I32   uint32
	I64   uint64
	Str   string
}

func (b *BytesReader) LoadFile(fileName string) {
	var err error
	b.File, err = os.Open(fileName)
	if err != nil {
		log.Fatalf("Failed to open %v", fileName)
	}
}

func (b *BytesReader) Move(newLoc int64) {
	ret, err := b.File.Seek(newLoc, 1)
	if err != nil && err != io.EOF {
		log.Fatal(err)
	}
	b.Loc = int(ret)
}

func (b *BytesReader) ReadNext(numOfBytes int) error {
	read := func() (int, error) {
		bytesRead, err := b.File.Read(b.Bytes)
		if err != nil && err != io.EOF {
			log.Fatal(err)
		}
		return bytesRead, err
	}

	if numOfBytes == 4 {
		b.I64 = 0
		b.Bytes = make([]byte, numOfBytes)
		bytesRead, err := read()
		b.Str = string(b.Bytes)
		buf := bytes.NewBuffer(b.Bytes)
		binary.Read(buf, binary.LittleEndian, &b.I32)
		b.Loc += bytesRead
		return err
	} else if numOfBytes == 8 {
		b.I32 = 0
		b.Bytes = make([]byte, numOfBytes)
		bytesRead, err := read()
		b.Str = string(b.Bytes)
		buf := bytes.NewBuffer(b.Bytes)
		binary.Read(buf, binary.LittleEndian, &b.I64)
		b.Loc += bytesRead
		return err
	} else {
		b.Bytes = make([]byte, numOfBytes)
		bytesRead, err := read()
		b.Str = string(b.Bytes)
		b.Loc += bytesRead
		return err
	}
}

func (b *BytesReader) PrintNext(numOfBytes int) error {
	err := b.ReadNext(numOfBytes)
	fmt.Printf("%v|%X|%v|%v|%v|\n", b.Loc, b.Bytes, b.I32, b.I64, b.Str)
	return err
}

func (r *BytesReader) ReadSegment() (common.Segment, error) {
	s := common.Segment{}
	s.Loc = r.Loc
	err := r.ReadNext(4)
	s.Header.Tag = r.Str
	err = r.ReadNext(4)
	s.Header.Toc = r.I32
	err = r.ReadNext(4)
	s.Header.Ver = r.I32
	err = r.ReadNext(8)
	s.Header.Seg = r.I64
	err = r.ReadNext(8)
	s.Header.Raw = r.I64
	// r.readNext(4)
	// s.NumOfObj = r.i32
	return s, err
}

func (r *BytesReader) ReadObj() common.Obj {
	o := common.Obj{}
	r.ReadNext(4)
	o.LenObjPath = r.I32
	r.ReadNext(int(o.LenObjPath))
	o.Name = r.Str
	r.ReadNext(4)
	o.RawDataIndex = r.I32
	r.ReadNext(4)
	o.NumOfProperties = r.I32
	return o
}

func (r *BytesReader) ReadProperty() common.Property {
	p := common.Property{}
	r.ReadNext(4)
	p.LenPropName = r.I32
	r.ReadNext(int(p.LenPropName))
	p.Name = r.Str
	r.ReadNext(4)
	p.Dtype = r.I32

	if p.Dtype == 32 {
		r.ReadNext(4)
		p.LenValStr = r.I32
		r.ReadNext(int(p.LenValStr))
		p.Value = r.Str
	} else {
		r.ReadNext(4)
		p.Value = r.I32
	}
	return p
}
