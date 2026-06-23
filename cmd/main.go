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
	// TODO : Verify database and version
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
			fmt.Println("Error parsing query:", err)
			continue
		}
		tableInstance, err = table.Open(query.Table)
		if err != nil && query.Type != "create" {
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

		case "delete":
			if err != nil {
				fmt.Println("Error opening table:", err)
				fmt.Print(">")
				continue
			}

			filters := []table.ColEq{}
			validFilter := true
			for _, v := range query.Filters {
				colIdx := slices.Index(tableInstance.GetColumns(), v.ColName)
				if colIdx == -1 {
					fmt.Println(table.ErrorInvalidInput, "Invalid column name", v.ColName)
					validFilter = false
					break
				}
				x := table.ColEq{
					ColIdx: colIdx,
					Value:  v.ColValue,
				}
				filters = append(filters, x)
			}

			if !validFilter {
				tableInstance.Close()
				fmt.Print(">")
				continue
			}

			err = tableInstance.Delete(filters)
			if err != nil {
				fmt.Println("Error deleting data:", err)
				fmt.Print(">")
				continue
			}
			tableInstance.Close()

		case "insert":
			err = tableInstance.Insert(query.Values)

			if err != nil {
				fmt.Println("Error inserting data:", err)
				fmt.Print(">")
				continue
			}
			tableInstance.Close()

		case "update":
			filters := []table.ColEq{}
			validFilter := true
			for _, v := range query.Filters {
				colIdx := slices.Index(tableInstance.GetColumns(), v.ColName)
				if colIdx == -1 {
					fmt.Println(table.ErrorInvalidInput, "Invalid column name", v.ColName)
					validFilter = false
					break
				}
				x := table.ColEq{
					ColIdx: colIdx,
					Value:  v.ColValue,
				}
				filters = append(filters, x)
			}

			if !validFilter {
				tableInstance.Close()
				fmt.Print(">")
				continue
			}

			updates := []table.UpdateValue{}
			validUpdate := true
			for _, v := range query.Updates {

				colIdx := slices.Index(tableInstance.GetColumns(), v.ColName)
				if colIdx == -1 {
					fmt.Println(table.ErrorInvalidInput, "Invalid column name", v.ColName)
					validUpdate = false
					break
				}
				x := table.UpdateValue{
					ColIdx: colIdx,
					Value:  v.ColValue,
				}
				updates = append(updates, x)
			}

			if !validUpdate {
				tableInstance.Close()
				fmt.Print(">")
				continue
			}

			err = tableInstance.Update(filters, updates)
			if err != nil {
				fmt.Println("Error updating data:", err)
				fmt.Print(">")
				continue
			}
			tableInstance.Close()

		case "create":
			_, err = table.Create(query.Table, query.Columns)
			if err != nil {
				fmt.Println("Error creating table:", err)
				fmt.Print(">")
				continue
			}

		}
		fmt.Print(">")
	}
}
