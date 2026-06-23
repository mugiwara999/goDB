package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/mugiwara999/goDB/internal/parser"
	"github.com/mugiwara999/goDB/internal/table"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("DB is starting")
	fmt.Print(">")

	for scanner.Scan() {
		command := strings.TrimSpace(scanner.Text())
		if command == "" {
			fmt.Print(">")
			continue
		}

		query, err := parser.Parse(command)
		if err != nil {
			fmt.Println(err)
			fmt.Print(">")
			continue
		}

		if query == nil {
			fmt.Print(">")
			continue
		}

		switch query.Type {
		case "create":
			tbl, err := table.Create(query.Table, query.Columns)
			if err != nil {
				fmt.Println(err)
				break
			}
			if err := tbl.Close(); err != nil {
				fmt.Println(err)
			}

		case "insert", "select", "delete", "update":
			tbl, err := table.Open(query.Table)
			if err != nil {
				fmt.Println(err)
				break
			}

			switch query.Type {
			case "insert":
				if err := tbl.Insert(query.Values); err != nil {
					fmt.Println(err)
				}
			case "select":
				filters, err := resolveFilters(tbl, query.Filters)
				if err != nil {
					fmt.Println(err)
					break
				}
				res, err := tbl.Select(query.Columns, filters)
				if err != nil {
					fmt.Println(err)
					break
				}
				for _, row := range res {
					fmt.Println(strings.Join(row, ","))
				}
			case "delete":
				filters, err := resolveFilters(tbl, query.Filters)
				if err != nil {
					fmt.Println(err)
					break
				}
				if err := tbl.Delete(filters); err != nil {
					fmt.Println(err)
				}
			case "update":
				filters, err := resolveFilters(tbl, query.Filters)
				if err != nil {
					fmt.Println(err)
					break
				}
				updates, err := resolveUpdates(tbl, query.Updates)
				if err != nil {
					fmt.Println(err)
					break
				}
				if err := tbl.Update(filters, updates); err != nil {
					fmt.Println(err)
				}
			}

			if err := tbl.Close(); err != nil {
				fmt.Println(err)
			}
		default:
			fmt.Printf("unsupported query type %q\n", query.Type)
		}

		fmt.Print(">")
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(table.ErrorTakingInput, err)
	}
}

func resolveFilters(tbl *table.Table, filters []parser.Pair) ([]table.ColEq, error) {
	resolved := make([]table.ColEq, 0, len(filters))
	columns := tbl.GetColumns()

	for _, f := range filters {
		colIdx := slices.Index(columns, f.ColName)
		if colIdx == -1 {
			return nil, fmt.Errorf("table %q: column %q does not exist", tbl.Name, f.ColName)
		}
		resolved = append(resolved, table.ColEq{
			ColIdx: colIdx,
			Value:  f.ColValue,
		})
	}

	return resolved, nil
}

func resolveUpdates(tbl *table.Table, updates []parser.Pair) ([]table.UpdateValue, error) {
	resolved := make([]table.UpdateValue, 0, len(updates))
	columns := tbl.GetColumns()

	for _, u := range updates {
		colIdx := slices.Index(columns, u.ColName)
		if colIdx == -1 {
			return nil, fmt.Errorf("table %q: column %q does not exist", tbl.Name, u.ColName)
		}
		resolved = append(resolved, table.UpdateValue{
			ColIdx: colIdx,
			Value:  u.ColValue,
		})
	}

	return resolved, nil
}
