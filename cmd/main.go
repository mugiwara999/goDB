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

var scanner *bufio.Scanner

var tableInstance *table.Table

func main() {
	fmt.Println("DB is starting")
	var command string
	// var err error
	fmt.Print(">")
	scanner = bufio.NewScanner(os.Stdin)
	if scanner.Err() != nil {
		fmt.Println(table.ErrorTakingInput, scanner.Err())
	}
	for scanner.Scan() {
		if scanner.Err() != nil {
			fmt.Println(table.ErrorTakingInput, scanner.Err())
		}
		command = scanner.Text()
		if command == "" {
			fmt.Print(">")
			continue
		}
		query, err := parser.Parse(strings.ToLower(command))
		if err != nil {
			fmt.Println("Error parsing query:", err)
			fmt.Print(">")
			continue
		}
		if query == nil {
			fmt.Print(">")
			continue
		}
		tableInstance, err = table.Open(query.Table)
		if err != nil {
			fmt.Println("Error opening table:", err)
			fmt.Print(">")
			continue
		}
		switch query.Type {
		case "select":
			filters := []table.ColEq{}
			for _, v := range query.Filters {
				colIdx := slices.Index(tableInstance.GetColumns(), v.ColName)
				if colIdx == -1 {
					fmt.Println(table.ErrorInvalidInput, "Invalid column name", v.ColName)
					fmt.Print(">")
					continue
				}
				x := table.ColEq{
					ColIdx: colIdx,
					Value:  v.ColValue,
				}
				filters = append(filters, x)
			}
			res, err := tableInstance.Select(query.Columns, filters)
			if err != nil {
				fmt.Println("Error selecting data:", err)
				fmt.Print(">")
				continue
			}
			for _, row := range res {
				fmt.Println(strings.Join(row, ","))
			}
			tableInstance.Close()
		}
		fmt.Print(">")
	}
}
