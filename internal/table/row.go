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

// func (t *Table) Delete(filters []ColEq) error {
// 	t.file.Seek(0, 0)
// 	filescanner := bufio.NewScanner(t.file)
//
// 	if filescanner.Err() != nil {
// 		return filescanner.Err()
// 	}
//
// 	var rows [][]string
//
// 	filescanner.Scan() // Skip the first line (column names)
// 	for filescanner.Scan() {
// 		rowText := filescanner.Text()
// 		row := strings.Split(rowText, ",")
//
// 		match := true
// 		for _, v := range filters {
// 			if row[v.ColIdx] != v.Value {
//
// 				match = false
// 			}
// 		}
//
// 		if !match {
// 			rows = append(rows, row)
// 		}
// 	}
// 	err := t.Truncate()
//
// 	if err != nil {
// 		return err
// 	}
// 	for _, row := range rows {
// 		t.file.WriteString(strings.Join(row, ",") + "\n")
// 	}
// 	return nil
// }
//
// func (t *Table) Truncate() error {
//
// 	err := t.file.Truncate(0)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = t.file.Seek(0, 0)
// 	if err != nil {
// 		return err
// 	}
// 	_, err = t.file.WriteString(string(strings.Join(t.GetColumns(), ",") + "\n"))
//
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
//
// func (t *Table) Update(filters []ColEq, toUpdate []UpdateValue) error {
//
// 	var rows [][]string
//
// 	t.file.Seek(0, 0)
//
// 	filescanner := bufio.NewScanner(t.file)
//
// 	if filescanner.Err() != nil {
// 		return filescanner.Err()
// 	}
//
// 	filescanner.Scan()
//
// 	for filescanner.Scan() {
//
// 		text := filescanner.Text()
//
// 		row := strings.Split(text, ",")
//
// 		match := true
//
// 		for _, v := range filters {
//
// 			if row[v.ColIdx] != v.Value {
// 				match = false
// 				break
// 			}
//
// 		}
//
// 		if match {
//
// 			for _, v := range toUpdate {
// 				row[v.ColIdx] = v.Value
// 			}
//
// 		}
//
// 		rows = append(rows, row)
// 	}
//
// 	err := t.Truncate()
//
// 	if err != nil {
// 		return err
// 	}
// 	for _, row := range rows {
// 		t.file.WriteString(strings.Join(row, ",") + "\n")
// 	}
// 	return nil
// }
