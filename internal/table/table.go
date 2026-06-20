package table

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mugiwara999/goDB/internal/pager"
)

type Table struct {
	Pager *pager.Pager
	cols  []string
}

var (
	ErrorCreatingFile   = fmt.Errorf("Error creating file")
	ErrorWritingToFile  = fmt.Errorf("Error writing to file")
	ErrorReadingFile    = fmt.Errorf("Error reading file")
	ErrorTakingInput    = fmt.Errorf("Error taking input")
	ErrorInvalidInput   = fmt.Errorf("Invalid input")
	ErrorColumnNotFound = fmt.Errorf("Column not found")
)

func Open(name string) (*Table, error) {

	// TODO : env load is duplicated
	err := godotenv.Load()
	var DataDir string

	if err != nil {
		DataDir = "../data"
	} else {
		DataDir = os.Getenv("DATA_DIR")
	}

	path := DataDir + "/" + strings.ToLower(name) + ".bin"

	pager, err := pager.NewPager(path)

	if err != nil {
		return nil, ErrorCreatingFile
	}

	cols := pager.GetColumns()

	table := &Table{
		Pager: pager,
		cols:  cols,
	}

	return table, nil

}
func Create(name string, cols []string) (*Table, error) {

	err := godotenv.Load()
	var DataDir string

	if err != nil {
		DataDir = "../data"
	} else {
		DataDir = os.Getenv("DATA_DIR")
	}

	path := DataDir + "/" + strings.ToLower(name) + ".bin"

	pager, err := pager.NewPager(path)

	if err != nil {
		return nil, ErrorCreatingFile
	}

	err = pager.WriteColumns(cols)

	if err != nil {
		return nil, fmt.Errorf("%w : %w", ErrorWritingToFile, err)
	}

	table := &Table{
		Pager: pager,
		cols:  cols,
	}
	return table, nil
}

func (t *Table) Close() error {

	return t.Pager.Close()

}
