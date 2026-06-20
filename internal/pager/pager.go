package pager

import "os"

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
