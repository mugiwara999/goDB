package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/mugiwara999/goDB/internal/table"
)

var scanner *bufio.Scanner

var tableInstance *table.Table

func main() {

	fmt.Println("DB is starting")

	fmt.Println("Enter 'create' to create a new table or 'use' to use an existing table")

	var command string
	var err error

	scanner = bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		command = scanner.Text()
		if scanner.Err() != nil {
			fmt.Println(table.ErrorTakingInput, scanner.Err())
		}

		if command != "create" && command != "use" {
			fmt.Println("Invalid command")
			return
		}
	}

	if command == "create" {

		fmt.Println("Enter table name")

		scanner.Scan()
		var tableName string = scanner.Text()

		fmt.Println("Enter column names separated by comma")

		var columnNames string

		if scanner.Err() != nil {
			fmt.Println("Error reading column names", scanner.Err())
		}

		var cols []string
		if scanner.Scan() {
			columnNames = scanner.Text()
			cols = strings.Split(columnNames, ",")
		}

		tableInstance, err = table.Create(tableName, cols)
		if err != nil {
			fmt.Println(table.ErrorReadingFile, err)
			return
		}
	} else {
		fmt.Println("Enter table name to use")
		scanner.Scan()
		if scanner.Err() != nil {
			fmt.Println(table.ErrorTakingInput, scanner.Err())
		}
		tableName := scanner.Text()
		tableInstance, err = table.Open(tableName)

		if err != nil {
			fmt.Println(table.ErrorReadingFile, err)
			return
		}
	}

	if tableInstance == nil {
		fmt.Println("Exiting due to error")
		return
	}

	defer tableInstance.Close()

	fmt.Println("Enter 'insert' or 'delete' or 'select' or 'update' or 'exit' to perform operations on the table")
	for scanner.Scan() {

		if scanner.Err() != nil {
			fmt.Println(table.ErrorTakingInput, scanner.Err())
		}
		command = scanner.Text()

		if command == "exit" {
			fmt.Println("Exiting...")
			return
		}

		if command != "insert" && command != "delete" && command != "select" && command != "update" {
			fmt.Println("Invalid command")
			return
		}

		switch command {
		case "insert":
			fmt.Println("Enter data rows seperated by comma  (type 'exit' to finish)")
			for scanner.Scan() {

				text := scanner.Text()

				if text == "exit" {
					break
				}

				vals := strings.Split(text, ",")

				if len(vals) != len(tableInstance.GetColumns()) {
					fmt.Println(table.ErrorInvalidInput, "Expected", len(tableInstance.GetColumns()), "values but got", len(vals))
					continue
				}

				err := tableInstance.Insert(vals)

				if err != nil {
					fmt.Println(table.ErrorWritingToFile, err)
				}
				fmt.Println("Data row added successfully")
			}
		case "select":
			fmt.Println("Enter column name and value to filter by in the format 'column=value (type 'exit' to finish)")
			var colEquals []table.ColEq
			for scanner.Scan() {
				text := scanner.Text()

				if text == "exit" {
					slices.SortFunc(colEquals, func(a, b table.ColEq) int {

						return a.ColIdx - b.ColIdx

					})
					res, err := tableInstance.Select(tableInstance.GetColumns(), colEquals)

					if err != nil {
						fmt.Println(err)
						break
					}

					for _, row := range res {
						fmt.Println(strings.Join(row, ","))
					}
					break
				}

				parts := strings.SplitN(text, "=", 2)
				if len(parts) != 2 {
					fmt.Println(table.ErrorInvalidInput, "Expected format 'column=value'")
					continue
				}

				colName := parts[0]
				colIdx := slices.Index(tableInstance.GetColumns(), colName)
				if colIdx == -1 {
					fmt.Println(table.ErrorInvalidInput, "Column name not found")
					continue
				}
				x := table.ColEq{
					ColIdx: colIdx,
					Value:  parts[1],
				}

				colEquals = append(colEquals, x)
			}

		case "delete":

			fmt.Println("Enter column name and value to filter by in the format 'column=value (type 'exit' to finish)")
			var colEquals []table.ColEq
			for scanner.Scan() {
				text := scanner.Text()

				if text == "exit" {
					slices.SortFunc(colEquals, func(a, b table.ColEq) int {
						return a.ColIdx - b.ColIdx

					})
					err := tableInstance.Delete(colEquals)

					if err != nil {
						fmt.Println(err)
						break
					}
					fmt.Println("Deletion successful")
					break
				}

				parts := strings.SplitN(text, "=", 2)
				if len(parts) != 2 {
					fmt.Println(table.ErrorInvalidInput, "Expected format 'column=value'")
					continue
				}

				colName := parts[0]
				colIdx := slices.Index(tableInstance.GetColumns(), colName)
				if colIdx == -1 {
					fmt.Println(table.ErrorInvalidInput, "Column name not found")
					continue
				}
				x := table.ColEq{
					ColIdx: colIdx,
					Value:  parts[1],
				}

				colEquals = append(colEquals, x)
			}
		case "update":

			fmt.Println("Enter column name and value to filter by in the format 'column=value (type 'exit' to finish)")
			var colEquals []table.ColEq
			var toUpdate []table.UpdateValue
			for scanner.Scan() {
				text := scanner.Text()

				if text == "exit" {
					slices.SortFunc(colEquals, func(a, b table.ColEq) int {
						return a.ColIdx - b.ColIdx

					})
					fmt.Println("Enter column name and value to update in the format 'column=value'")
					for scanner.Scan() {
						text := scanner.Text()
						if text == "exit" {
							err := tableInstance.Update(colEquals, toUpdate)
							if err != nil {
								fmt.Println(err)
							}
							fmt.Println("Update successful")
							break
						}
						parts := strings.SplitN(text, "=", 2)
						if len(parts) != 2 {
							fmt.Println(table.ErrorInvalidInput, "Expected format heeee 'column=value'")
							continue
						}

						colName := parts[0]
						colIdx := slices.Index(tableInstance.GetColumns(), colName)
						if colIdx == -1 {
							fmt.Println(table.ErrorInvalidInput, "Column name not found")
							continue
						}
						newVal := parts[1]

						toUpdate = append(toUpdate, table.UpdateValue{
							ColIdx: colIdx,
							Value:  newVal,
						})
					}
					break
				}

				parts := strings.SplitN(text, "=", 2)
				if len(parts) != 2 {
					fmt.Println(table.ErrorInvalidInput, "Expected format 'column=value'")
					continue
				}
				colName := parts[0]
				colIdx := slices.Index(tableInstance.GetColumns(), colName)
				if colIdx == -1 {
					fmt.Println(table.ErrorInvalidInput, "Column name not found")
					continue
				}
				x := table.ColEq{
					ColIdx: colIdx,
					Value:  parts[1],
				}

				colEquals = append(colEquals, x)
			}

		}
		fmt.Println("Enter 'insert' or 'delete' or 'select' or 'update' or 'exit' to perform operations on the table")
	}

}
