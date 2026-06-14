package table

import (
	"bufio"
	"strings"
)

type ColEq struct {
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

func (t *Table) Delete() error {

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
