package table

import (
	// "bufio"
	"fmt"
	// "strings"

	"github.com/mugiwara999/goDB/internal/pager"
)

type ColEq struct {
	ColIdx int
	Value  string
}

type UpdateValue struct {
	ColIdx int
	Value  string
}

func SerializeRow(values []string) []byte {
	buf := make([]byte, 0)
	for _, v := range values {
		buf = append(buf, []byte(v)...)
		buf = append(buf, 0) // null terminator
	}
	return buf
}

func DeserializeRow(data []byte) []string {
	values := make([]string, 0)
	accumulator := make([]byte, 0)
	for _, b := range data {
		if b == 0 {
			values = append(values, string(accumulator))
			accumulator = accumulator[:0]
			continue
		}
		accumulator = append(accumulator, b)
	}
	return values
}

func (t *Table) GetColumns() []string {
	return t.cols
}

func (t *Table) Insert(values []string) error {

	var lastPage *pager.Page
	var err error

	if t.Pager.GetNumPages() <= 1 {
		lastPage, err = t.Pager.NewPage()

		if err != nil {
			return fmt.Errorf("%w : %w", pager.ErrorReadingPage, err)
		}
	} else {

		lastPage, err = t.Pager.GetPage(t.Pager.GetNumPages() - 1)

		if err != nil {
			return fmt.Errorf("%w : %w", pager.ErrorReadingPage, err)
		}

		if !lastPage.CanFit(len(SerializeRow(values))) {
			lastPage, _ = t.Pager.NewPage()
		}
	}

	err = lastPage.AddRow(SerializeRow(values))

	if err != nil {
		return fmt.Errorf("%w : %w", pager.ErrorWritingPage, err)
	}

	err = t.Pager.Flush(lastPage)

	return err
}

func (t *Table) Select(columns []string, colEquals []ColEq) ([][]string, error) {

	result := make([][]string, 0)

	if len(columns) > 0 && columns[0] == "*" {
		columns = t.GetColumns()
	}
	colMap := make(map[string]int)
	colSet := make(map[string]struct{})

	for _, col := range columns {
		colSet[col] = struct{}{}
	}
	for i, col := range t.cols {
		colMap[col] = i
	}

	colIdxs := make([]int, 0)
	for _, col := range columns {
		idx := colMap[col]
		colIdxs = append(colIdxs, idx)
	}

	result = append(result, columns)

	rowIt := t.Pager.RowIterator()

	for {

		Srow, err := rowIt.Next()

		if Srow == nil || err != nil {
			break
		}

		row := DeserializeRow(Srow)

		match := true
		for _, v := range colEquals {
			if row[v.ColIdx] != v.Value {

				match = false

			}
		}

		if match {

			res := make([]string, 0)

			for idx := range colIdxs {

				res = append(res, row[idx])

			}

			result = append(result, res)
		}

	}

	return result, nil
}

func (t *Table) Delete(filters []ColEq) error {

	rowIt := t.Pager.RowIterator()

	for {

		rowText, err := rowIt.Next()
		if rowText == nil || err != nil {
			break
		}
		row := DeserializeRow(rowText)

		match := true
		for _, v := range filters {
			if row[v.ColIdx] != v.Value {
				match = false
			}
		}

		if match {
			matchInfo := rowIt.GetCurrentInfo()
			page, err := t.Pager.GetPage(matchInfo.PageID)

			if err != nil {
				return pager.ErrorReadingPage
			}

			page.DeleteRow(matchInfo.SlotID)

			err = t.Pager.Flush(page)

		}
	}
	return nil
}

func (t *Table) Update(filters []ColEq, toUpdate []UpdateValue) error {

	rowIt := t.Pager.RowIterator()

	for {

		data, _ := rowIt.Next()

		if data == nil {
			break
		}

		row := DeserializeRow(data)

		match := true

		for _, v := range filters {

			if row[v.ColIdx] != v.Value {
				match = false
				break
			}

		}

		if match {

			for _, v := range toUpdate {
				row[v.ColIdx] = v.Value
			}

			info := rowIt.GetCurrentInfo()

			newData := SerializeRow(row)

			page, err := t.Pager.GetPage(info.PageID)

			if err != nil {
				return err
			}

			if len(newData) == len(data) {
				page.Overwrite(info.SlotID, newData)
				t.Pager.Flush(page)
			} else {
				page.DeleteRow(info.SlotID)
				t.Pager.Flush(page)
				err := t.Insert(row)
				if err != nil {
					return err
				}
			}

		}

	}

	return nil
}
