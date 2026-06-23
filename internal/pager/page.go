package pager

import (
	"encoding/binary"
	"fmt"
	"os"
)

const PAGE_SIZE = 4096

type Page struct {
	ID   int
	Data [PAGE_SIZE]byte
}

var (
	ErrorWritingPage    = fmt.Errorf("could not write full page")
	ErrorReadingPage    = fmt.Errorf("could not read full page")
	ErrorNotEnoughSpace = fmt.Errorf("space isn't enough to fit this data")
)

// Bytes 0-1:   numSlots  (how many rows on this page)
// Bytes 2-3:   freeStart (where free space starts, after slots)
// Bytes 4-5:   freeEnd   (where free space ends, before row data)
// Bytes 6+:    slot array - each slot is 5 bytes
// 2 bytes since they are used to represent offset, 1 byte can represent only till 256
// each slot is 4 bytes. 2 bytes for offset, 2 bytes for len of row
// and last slot is a flag to represent a deleted slot
// since after deletion moving all the slots is expensive
// so my approach is to flag it and have a function run every 24 hrs
// which rebuilds the page with only the active slots

func (p *Page) Init() {
	binary.LittleEndian.PutUint16(p.Data[0:2], 0)         // numSlots = 0
	binary.LittleEndian.PutUint16(p.Data[2:4], 6)         // freeStart = 6 (after header)
	binary.LittleEndian.PutUint16(p.Data[4:6], PAGE_SIZE) // freeEnd = 4096
}

func WritePage(file *os.File, page *Page) error {
	_, err := file.Seek(int64(page.ID)*PAGE_SIZE, 0)

	if err != nil {
		return err
	}

	n, err := file.Write(page.Data[:])

	if err != nil {
		return err
	}

	if n < PAGE_SIZE {
		return ErrorWritingPage
	}

	return nil
}

func ReadPage(file *os.File, pageID int) ([]byte, error) {

	_, err := file.Seek(int64(pageID)*PAGE_SIZE, 0)

	if err != nil {
		return nil, ErrorReadingPage
	}

	data := make([]byte, PAGE_SIZE)

	n, err := file.Read(data)

	if n != PAGE_SIZE {
		return nil, fmt.Errorf("%w : Read only %d bytes", ErrorReadingPage, n)
	}

	if err != nil {
		return nil, fmt.Errorf("%w : %w", ErrorReadingPage, err)
	}
	return data, nil
}

func (p *Page) AddRow(data []byte) error {

	// numSlots := int(p.Data[0]) | int(p.Data[1])<<8
	numSlots := int(binary.LittleEndian.Uint16(p.Data[:2]))
	// freeStart := int(p.Data[2]) | int(p.Data[3])<<8
	freeStart := int(binary.LittleEndian.Uint16(p.Data[2:4]))
	// freeEnd := int(p.Data[4]) | int(p.Data[5])<<8
	freeEnd := int(binary.LittleEndian.Uint16(p.Data[4:6]))

	if freeEnd-freeStart < len(data)+5 {
		return ErrorNotEnoughSpace
	}

	start := freeEnd - len(data)

	n := copy(p.Data[start:freeEnd], data)

	if n != len(data) {
		return fmt.Errorf("%w : Read only %d bytes", ErrorWritingPage, n)
	}

	numSlots++ // increase num of rows

	binary.LittleEndian.PutUint16(p.Data[:2], uint16(numSlots))

	binary.LittleEndian.PutUint16(p.Data[freeStart:freeStart+2], uint16(start)) // add offset to slots
	binary.LittleEndian.PutUint16(p.Data[freeStart+2:freeStart+4], uint16(len(data)))
	p.Data[freeStart+4] = 0

	freeStart += 5

	binary.LittleEndian.PutUint16(p.Data[2:4], uint16(freeStart))

	freeEnd = start - 1

	binary.LittleEndian.PutUint16(p.Data[4:6], uint16(freeEnd))

	return nil

}

func (p *Page) GetRow(slotID int) []byte {
	slotPos := 6 + slotID*5
	offset := int(binary.LittleEndian.Uint16(p.Data[slotPos : slotPos+2]))
	size := int(binary.LittleEndian.Uint16(p.Data[slotPos+2 : slotPos+4]))

	data := p.Data[offset : offset+size]

	return data

}

func (p *Page) DeleteRow(slotID int) {
	offsetPos := 6 + slotID*5
	p.Data[offsetPos+4] = 1
}

func (p *Page) CanFit(dataLen int) bool {

	freeStart := int(binary.LittleEndian.Uint16(p.Data[2:4]))
	freeEnd := int(binary.LittleEndian.Uint16(p.Data[4:6]))

	return (freeEnd - freeStart) >= dataLen+5

}

func (p *Page) isDeleted(slotId int) bool {
	offset := 6 + slotId*5

	x := int(p.Data[offset+4])

	return x == 1

}

func (p *Page) Overwrite(slotId int, data []byte) error {
	offset := 6 + slotId*5
	size := int(binary.LittleEndian.Uint16(p.Data[offset+2 : offset+4]))
	if len(data) > size {
		return fmt.Errorf("data too large for slot")
	}
	start := int(binary.LittleEndian.Uint16(p.Data[offset : offset+2]))
	copy(p.Data[start:start+len(data)], data)
	return nil
}
