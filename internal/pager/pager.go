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
	_, err := p.file.Seek(int64(offset), 0)

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
// Bytes 6-9  : numCols
// Bytes 10+ : column names (null-terminated strings)

func (p *Pager) GetColumns() []string {
	if p.numPages == 0 {
		return nil
	}
	page, err := p.GetPage(0)

	if err != nil || page == nil {
		return nil
	}

	numCols := int(binary.LittleEndian.Uint32(page.Data[6:10]))

	cols := make([]string, numCols)

	accumulator := make([]byte, 0)

	offset := 10

	colIdx := 0

	for i := offset; i < PAGE_SIZE && colIdx < numCols; i++ {
		if page.Data[i] == 0 {
			cols[colIdx] = string(accumulator)
			accumulator = accumulator[:0]
			colIdx++
			continue
		}
		accumulator = append(accumulator, page.Data[i])
	}

	return cols

}

func (p *Pager) WriteColumns(cols []string) error {
	page, _ := p.NewPage() // page 0

	binary.LittleEndian.PutUint32(page.Data[:4], uint32(0x474f4442)) // "GODB"
	binary.LittleEndian.PutUint16(page.Data[4:6], uint16(1))         // version
	binary.LittleEndian.PutUint32(page.Data[6:10], uint32(len(cols)))

	offset := 10

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
