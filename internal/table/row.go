package table

import (
	"bufio"
	"strings"
)

type ColEq struct {
	ColIdx int
	Value  string
}

type UpdateValue struct {
	ColIdx int
	Value  string
}

func (t *Table) GetColumns() []string {
	return t.cols
}

func (t *Table) Insert(values []string) error {
	_, err := t.file.WriteString(strings.Join(values, ",") + "\n")
	return err
}

func (t *Table) Select(columns []string, colEquals []ColEq) ([][]string, error) {

	result := make([][]string, 0)

	t.file.Seek(0, 0)
	fileScanner := bufio.NewScanner(t.file)

	if fileScanner.Err() != nil {
		return nil, fileScanner.Err()
	}

	colSet := make(map[string]struct{})
	colIdxSet := make(map[int]struct{})

	for _, col := range columns {
		colSet[col] = struct{}{}
	}

	for i, v := range t.cols {
		if _, ok := colSet[v]; ok {
			colIdxSet[i] = struct{}{}
		}
	}

	fileScanner.Scan() // Skip the first line (column names)
	result = append(result, columns)

	for fileScanner.Scan() {
		rowText := fileScanner.Text()
		row := strings.Split(rowText, ",")

		match := true
		for _, v := range colEquals {
			if row[v.ColIdx] != v.Value {

				match = false

			}
		}

		if match {

			res := make([]string, 0)

			for idx := range colIdxSet {

				res = append(res, row[idx])

			}

			result = append(result, res)
		}

	}

	return result, nil
}

func (t *Table) Delete(filters []ColEq) error {
	t.file.Seek(0, 0)
	filescanner := bufio.NewScanner(t.file)

	if filescanner.Err() != nil {
		return filescanner.Err()
	}

	var rows [][]string

	filescanner.Scan() // Skip the first line (column names)
	for filescanner.Scan() {
		rowText := filescanner.Text()
		row := strings.Split(rowText, ",")

		match := true
		for _, v := range filters {
			if row[v.ColIdx] != v.Value {

				match = false
			}
		}

		if !match {
			rows = append(rows, row)
		}
	}
	err := t.Truncate()

	if err != nil {
		return err
	}
	for _, row := range rows {
		t.file.WriteString(strings.Join(row, ",") + "\n")
	}
	return nil
}

func (t *Table) Truncate() error {

	err := t.file.Truncate(0)
	if err != nil {
		return err
	}
	_, err = t.file.Seek(0, 0)
	if err != nil {
		return err
	}
	_, err = t.file.WriteString(string(strings.Join(t.GetColumns(), ",") + "\n"))

	if err != nil {
		return err
	}
	return nil
}

func (t *Table) Update(filters []ColEq, toUpdate []UpdateValue) error {

	var rows [][]string

	t.file.Seek(0, 0)

	filescanner := bufio.NewScanner(t.file)

	if filescanner.Err() != nil {
		return filescanner.Err()
	}

	filescanner.Scan()

	for filescanner.Scan() {

		text := filescanner.Text()

		row := strings.Split(text, ",")

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

		}

		rows = append(rows, row)
	}

	err := t.Truncate()

	if err != nil {
		return err
	}
	for _, row := range rows {
		t.file.WriteString(strings.Join(row, ",") + "\n")
	}
	return nil
}
