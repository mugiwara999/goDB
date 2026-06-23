package table

import (
	"fmt"

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
		buf = append(buf, 0)
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

func (t *Table) Insert(values []string) error {
	if len(values) != len(t.cols) {
		return fmt.Errorf("insert into table %q: expected %d values, got %d", t.Name, len(t.cols), len(values))
	}

	var lastPage *pager.Page
	var err error

	if t.Pager.GetNumPages() <= 1 {
		lastPage, err = t.Pager.NewPage()
		if err != nil {
			return fmt.Errorf("insert into table %q: %w", t.Name, err)
		}
	} else {
		lastPage, err = t.Pager.GetPage(t.Pager.GetNumPages() - 1)
		if err != nil {
			return fmt.Errorf("insert into table %q: %w", t.Name, err)
		}

		if !lastPage.CanFit(len(SerializeRow(values))) {
			lastPage, err = t.Pager.NewPage()
			if err != nil {
				return fmt.Errorf("insert into table %q: %w", t.Name, err)
			}
		}
	}

	rowData := SerializeRow(values)
	if err := lastPage.AddRow(rowData); err != nil {
		return fmt.Errorf("insert into table %q: %w", t.Name, err)
	}

	if err := t.Pager.Flush(lastPage); err != nil {
		return fmt.Errorf("insert into table %q: %w", t.Name, err)
	}

	return nil
}

func (t *Table) Select(columns []string, colEquals []ColEq) ([][]string, error) {
	result := make([][]string, 0)

	if len(columns) > 0 && columns[0] == "*" {
		columns = t.GetColumns()
	}

	if len(columns) == 0 {
		return nil, fmt.Errorf("select from table %q: no columns requested", t.Name)
	}

	colMap := make(map[string]int)
	for i, col := range t.cols {
		colMap[col] = i
	}

	colIdxs := make([]int, 0, len(columns))
	for _, col := range columns {
		idx, ok := colMap[col]
		if !ok {
			return nil, fmt.Errorf("select from table %q: column %q does not exist", t.Name, col)
		}
		colIdxs = append(colIdxs, idx)
	}

	result = append(result, columns)

	rowIt := t.Pager.RowIterator()
	for {
		rowData, err := rowIt.Next()
		if err != nil {
			return nil, fmt.Errorf("select from table %q: %w", t.Name, err)
		}
		if rowData == nil {
			break
		}

		row := DeserializeRow(rowData)
		if len(row) < len(t.cols) {
			return nil, fmt.Errorf("select from table %q: corrupt row has %d values, expected %d", t.Name, len(row), len(t.cols))
		}

		match := true
		for _, v := range colEquals {
			if v.ColIdx < 0 || v.ColIdx >= len(row) {
				return nil, fmt.Errorf("select from table %q: filter column index %d is out of range for row with %d values", t.Name, v.ColIdx, len(row))
			}
			if row[v.ColIdx] != v.Value {
				match = false
				break
			}
		}

		if match {
			res := make([]string, 0, len(colIdxs))
			for _, idx := range colIdxs {
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
		rowData, err := rowIt.Next()
		if err != nil {
			return fmt.Errorf("delete from table %q: %w", t.Name, err)
		}
		if rowData == nil {
			break
		}

		row := DeserializeRow(rowData)
		if len(row) < len(t.cols) {
			return fmt.Errorf("delete from table %q: corrupt row has %d values, expected %d", t.Name, len(row), len(t.cols))
		}

		match := true
		for _, v := range filters {
			if v.ColIdx < 0 || v.ColIdx >= len(row) {
				return fmt.Errorf("delete from table %q: filter column index %d is out of range for row with %d values", t.Name, v.ColIdx, len(row))
			}
			if row[v.ColIdx] != v.Value {
				match = false
				break
			}
		}

		if match {
			matchInfo := rowIt.GetCurrentInfo()
			page, err := t.Pager.GetPage(matchInfo.PageID)
			if err != nil {
				return fmt.Errorf("delete from table %q page %d slot %d: %w", t.Name, matchInfo.PageID, matchInfo.SlotID, err)
			}

			if err := page.DeleteRow(matchInfo.SlotID); err != nil {
				return fmt.Errorf("delete from table %q page %d slot %d: %w", t.Name, matchInfo.PageID, matchInfo.SlotID, err)
			}

			if err := t.Pager.Flush(page); err != nil {
				return fmt.Errorf("delete from table %q page %d slot %d: %w", t.Name, matchInfo.PageID, matchInfo.SlotID, err)
			}
		}
	}

	return nil
}

func (t *Table) Update(filters []ColEq, toUpdate []UpdateValue) error {
	rowIt := t.Pager.RowIterator()

	for {
		rowData, err := rowIt.Next()
		if err != nil {
			return fmt.Errorf("update table %q: %w", t.Name, err)
		}
		if rowData == nil {
			break
		}

		row := DeserializeRow(rowData)
		if len(row) < len(t.cols) {
			return fmt.Errorf("update table %q: corrupt row has %d values, expected %d", t.Name, len(row), len(t.cols))
		}

		match := true
		for _, v := range filters {
			if v.ColIdx < 0 || v.ColIdx >= len(row) {
				return fmt.Errorf("update table %q: filter column index %d is out of range for row with %d values", t.Name, v.ColIdx, len(row))
			}
			if row[v.ColIdx] != v.Value {
				match = false
				break
			}
		}

		if !match {
			continue
		}

		for _, v := range toUpdate {
			if v.ColIdx < 0 || v.ColIdx >= len(row) {
				return fmt.Errorf("update table %q: update column index %d is out of range for row with %d values", t.Name, v.ColIdx, len(row))
			}
			row[v.ColIdx] = v.Value
		}

		info := rowIt.GetCurrentInfo()
		newData := SerializeRow(row)

		page, err := t.Pager.GetPage(info.PageID)
		if err != nil {
			return fmt.Errorf("update table %q page %d slot %d: %w", t.Name, info.PageID, info.SlotID, err)
		}

		if len(newData) == len(rowData) {
			if err := page.Overwrite(info.SlotID, newData); err != nil {
				return fmt.Errorf("update table %q page %d slot %d: %w", t.Name, info.PageID, info.SlotID, err)
			}
			if err := t.Pager.Flush(page); err != nil {
				return fmt.Errorf("update table %q page %d slot %d: %w", t.Name, info.PageID, info.SlotID, err)
			}
			continue
		}

		if err := page.DeleteRow(info.SlotID); err != nil {
			return fmt.Errorf("update table %q page %d slot %d: %w", t.Name, info.PageID, info.SlotID, err)
		}
		if err := t.Pager.Flush(page); err != nil {
			return fmt.Errorf("update table %q page %d slot %d: %w", t.Name, info.PageID, info.SlotID, err)
		}

		if err := t.Insert(row); err != nil {
			return fmt.Errorf("update table %q page %d slot %d: %w", t.Name, info.PageID, info.SlotID, err)
		}
	}

	return nil
}
