package table

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/mugiwara999/goDB/internal/pager"
)

type Table struct {
	Pager *pager.Pager
	Name  string
	Path  string
	cols  []string
}

var (
	ErrorTakingInput = errors.New("read user input")
)

func Open(name string) (*Table, error) {
	path, err := tablePath(name)
	if err != nil {
		return nil, fmt.Errorf("open table %q: %w", name, err)
	}

	pg, err := pager.OpenPager(path)
	if err != nil {
		return nil, fmt.Errorf("open table %q at %q: %w", name, path, err)
	}

	cols, err := pg.GetColumns()
	if err != nil {
		_ = pg.Close()
		return nil, fmt.Errorf("open table %q at %q: %w", name, path, err)
	}

	return &Table{
		Pager: pg,
		Name:  strings.ToLower(name),
		Path:  path,
		cols:  cols,
	}, nil
}

func Create(name string, cols []string) (*Table, error) {
	if len(cols) == 0 {
		return nil, fmt.Errorf("create table %q: at least one column name is required", name)
	}

	path, err := tablePath(name)
	if err != nil {
		return nil, fmt.Errorf("create table %q: %w", name, err)
	}

	pg, err := pager.CreatePager(path)
	if err != nil {
		return nil, fmt.Errorf("create table %q at %q: %w", name, path, err)
	}

	if err := pg.WriteColumns(cols); err != nil {
		_ = pg.Close()
		return nil, fmt.Errorf("create table %q at %q: %w", name, path, err)
	}

	return &Table{
		Pager: pg,
		Name:  strings.ToLower(name),
		Path:  path,
		cols:  cols,
	}, nil
}

func (t *Table) Close() error {
	if t == nil || t.Pager == nil {
		return nil
	}
	return t.Pager.Close()
}

func (t *Table) GetColumns() []string {
	cols := make([]string, len(t.cols))
	copy(cols, t.cols)
	return cols
}

func tablePath(name string) (string, error) {
	if err := godotenv.Load(); err != nil {
		return "../data/" + strings.ToLower(name) + ".bin", nil
	}

	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "../data"
	}
	return dataDir + "/" + strings.ToLower(name) + ".bin", nil
}
