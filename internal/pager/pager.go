package pager

import (
	"encoding/binary"
	"os"
	// "strings"
)

type Pager struct {
	file     *os.File
	numPages int
}

// page 0 for meta data

func NewPager(path string) (*Pager, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	numPages := int(info.Size()) / PAGE_SIZE
	return &Pager{file: file, numPages: numPages}, nil
}

func (p *Pager) GetPage(pageID int) (*Page, error) {
	offset := pageID * PAGE_SIZE
	_, err := p.file.Seek(0, 0)

	if err != nil {
		return nil, ErrorReadingPage
	}
	_, err = p.file.Seek(int64(offset), 0)

	if err != nil {
		return nil, ErrorReadingPage
	}

	data := [PAGE_SIZE]byte{}

	_, err = p.file.Read(data[:])

	if err != nil {
		return nil, ErrorReadingPage
	}

	page := &Page{
		ID:   pageID,
		Data: data,
	}

	return page, nil

}

func (p *Pager) NewPage() (*Page, error) {
	data := [PAGE_SIZE]byte{}
	page := &Page{
		ID:   p.numPages,
		Data: data,
	}

	page.Init()
	p.numPages++

	return page, nil
}

func (p *Pager) Flush(page *Page) error {

	return WritePage(p.file, page)

}

// meta data
// Bytes 0-3   : GODB
// Bytes 4-5   : version
// Bytes 6-9  : numPages
// Bytes 10+ : column names (null-terminated strings)

func (p *Pager) GetColumns() []string {
	page, err := p.GetPage(0)

	if err != nil || page == nil {
		return nil
	}

	info, err := p.file.Stat()

	if err != nil {
		return nil
	}

	numCols := int(info.Size() / PAGE_SIZE)

	cols := make([]string, numCols)

	accumulator := make([]byte, 0)

	colIdx := 0

	for i := 12; i < PAGE_SIZE && colIdx < numCols; i++ {
		if page.Data[i] == 0 {
			cols[colIdx] = string(accumulator)
			accumulator = accumulator[:0]
			colIdx++
		}
		accumulator = append(accumulator, page.Data[i])
	}

	return cols

}

func (p *Pager) WriteColumns(cols []string) error {
	page, _ := p.NewPage() // page 0

	binary.LittleEndian.PutUint32(page.Data[:4], uint32(0x474f4442)) // "GODB"
	binary.LittleEndian.PutUint16(page.Data[4:6], uint16(1))         // version
	binary.LittleEndian.PutUint16(page.Data[6:8], uint16(len(cols)))
	binary.LittleEndian.PutUint32(page.Data[8:12], uint32(p.numPages))

	offset := 12

	for _, col := range cols {
		copy(page.Data[offset:], []byte(col))
		offset += len(col)
		page.Data[offset] = 0
		offset++
	}

	return p.Flush(page)
}

func (p *Pager) Close() error {
	return p.file.Close()
}
