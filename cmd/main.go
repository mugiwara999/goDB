package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"

	"github.com/mugiwara999/goDB/internal/parser"
	"github.com/mugiwara999/goDB/internal/table"

	// "slices"
	"strings"
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
				continue
			}
			for _, row := range res {
				fmt.Println(strings.Join(row, ","))
			}
			defer tableInstance.Close()
		}
		fmt.Print(">")
	}

	// fmt.Println("Enter 'insert' or 'delete' or 'select' or 'update' or 'exit' to perform operations on the table")
	// for scanner.Scan() {
	//
	// 	if scanner.Err() != nil {
	// 		fmt.Println(table.ErrorTakingInput, scanner.Err())
	// 	}
	// 	command = scanner.Text()
	//
	// 	if command == "exit" {
	// 		fmt.Println("Exiting...")
	// 		return
	// 	}
	//
	// 	if command != "insert" && command != "delete" && command != "select" && command != "update" {
	// 		fmt.Println("Invalid command")
	// 		return
	// 	}
	//
	// 	switch command {
	// 	case "insert":
	// 		fmt.Println("Enter data rows seperated by comma  (type 'exit' to finish)")
	// 		for scanner.Scan() {
	//
	// 			text := scanner.Text()
	//
	// 			if text == "exit" {
	// 				break
	// 			}
	//
	// 			vals := strings.Split(text, ",")
	//
	// 			if len(vals) != len(tableInstance.GetColumns()) {
	// 				fmt.Println(table.ErrorInvalidInput, "Expected", len(tableInstance.GetColumns()), "values but got", len(vals))
	// 				continue
	// 			}
	//
	// 			err := tableInstance.Insert(vals)
	//
	// 			if err != nil {
	// 				fmt.Println(table.ErrorWritingToFile, err)
	// 			}
	// 			fmt.Println("Data row added successfully")
	// 		}
	// 	case "select":
	// 		fmt.Println("Enter column name and value to filter by in the format 'column=value (type 'exit' to finish)")
	// 		var colEquals []table.ColEq
	// 		for scanner.Scan() {
	// 			text := scanner.Text()
	//
	// 			if text == "exit" {
	// 				slices.SortFunc(colEquals, func(a, b table.ColEq) int {
	//
	// 					return a.ColIdx - b.ColIdx
	//
	// 				})
	// 				res, err := tableInstance.Select(tableInstance.GetColumns(), colEquals)
	//
	// 				if err != nil {
	// 					fmt.Println(err)
	// 					break
	// 				}
	//
	// 				for _, row := range res {
	// 					fmt.Println(strings.Join(row, ","))
	// 				}
	// 				break
	// 			}
	//
	// 			parts := strings.SplitN(text, "=", 2)
	// 			if len(parts) != 2 {
	// 				fmt.Println(table.ErrorInvalidInput, "Expected format 'column=value'")
	// 				continue
	// 			}
	//
	// 			colName := parts[0]
	// 			colIdx := slices.Index(tableInstance.GetColumns(), colName)
	// 			if colIdx == -1 {
	// 				fmt.Println(table.ErrorInvalidInput, "Column name not found")
	// 				continue
	// 			}
	// 			x := table.ColEq{
	// 				ColIdx: colIdx,
	// 				Value:  parts[1],
	// 			}
	//
	// 			colEquals = append(colEquals, x)
	// 		}
	//
	// 	case "delete":
	//
	// 		fmt.Println("Enter column name and value to filter by in the format 'column=value (type 'exit' to finish)")
	// 		var colEquals []table.ColEq
	// 		for scanner.Scan() {
	// 			text := scanner.Text()
	//
	// 			if text == "exit" {
	// 				slices.SortFunc(colEquals, func(a, b table.ColEq) int {
	// 					return a.ColIdx - b.ColIdx
	//
	// 				})
	// 				err := tableInstance.Delete(colEquals)
	//
	// 				if err != nil {
	// 					fmt.Println(err)
	// 					break
	// 				}
	// 				fmt.Println("Deletion successful")
	// 				break
	// 			}
	//
	// 			parts := strings.SplitN(text, "=", 2)
	// 			if len(parts) != 2 {
	// 				fmt.Println(table.ErrorInvalidInput, "Expected format 'column=value'")
	// 				continue
	// 			}
	//
	// 			colName := parts[0]
	// 			colIdx := slices.Index(tableInstance.GetColumns(), colName)
	// 			if colIdx == -1 {
	// 				fmt.Println(table.ErrorInvalidInput, "Column name not found")
	// 				continue
	// 			}
	// 			x := table.ColEq{
	// 				ColIdx: colIdx,
	// 				Value:  parts[1],
	// 			}
	//
	// 			colEquals = append(colEquals, x)
	// 		}
	// 	case "update":
	//
	// 		fmt.Println("Enter column name and value to filter by in the format 'column=value (type 'exit' to finish)")
	// 		var colEquals []table.ColEq
	// 		var toUpdate []table.UpdateValue
	// 		for scanner.Scan() {
	// 			text := scanner.Text()
	//
	// 			if text == "exit" {
	// 				slices.SortFunc(colEquals, func(a, b table.ColEq) int {
	// 					return a.ColIdx - b.ColIdx
	//
	// 				})
	// 				fmt.Println("Enter column name and value to update in the format 'column=value'")
	// 				for scanner.Scan() {
	// 					text := scanner.Text()
	// 					if text == "exit" {
	// 						err := tableInstance.Update(colEquals, toUpdate)
	// 						if err != nil {
	// 							fmt.Println(err)
	// 						}
	// 						fmt.Println("Update successful")
	// 						break
	// 					}
	// 					parts := strings.SplitN(text, "=", 2)
	// 					if len(parts) != 2 {
	// 						fmt.Println(table.ErrorInvalidInput, "Expected format heeee 'column=value'")
	// 						continue
	// 					}
	//
	// 					colName := parts[0]
	// 					colIdx := slices.Index(tableInstance.GetColumns(), colName)
	// 					if colIdx == -1 {
	// 						fmt.Println(table.ErrorInvalidInput, "Column name not found")
	// 						continue
	// 					}
	// 					newVal := parts[1]
	//
	// 					toUpdate = append(toUpdate, table.UpdateValue{
	// 						ColIdx: colIdx,
	// 						Value:  newVal,
	// 					})
	// 				}
	// 				break
	// 			}
	//
	// 			parts := strings.SplitN(text, "=", 2)
	// 			if len(parts) != 2 {
	// 				fmt.Println(table.ErrorInvalidInput, "Expected format 'column=value'")
	// 				continue
	// 			}
	// 			colName := parts[0]
	// 			colIdx := slices.Index(tableInstance.GetColumns(), colName)
	// 			if colIdx == -1 {
	// 				fmt.Println(table.ErrorInvalidInput, "Column name not found")
	// 				continue
	// 			}
	// 			x := table.ColEq{
	// 				ColIdx: colIdx,
	// 				Value:  parts[1],
	// 			}
	//
	// 			colEquals = append(colEquals, x)
	// 		}
	//
	// 	}
	// 	fmt.Println("Enter 'insert' or 'delete' or 'select' or 'update' or 'exit' to perform operations on the table")
	// }

}
