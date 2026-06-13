package main

import (
	"bufio"
	"fmt"
	// "log"
	"bytes"
	"os"
	"strings"
)

var (
	ErrorCreatingFile  = fmt.Errorf("Error creating file")
	ErrorWritingToFile = fmt.Errorf("Error writing to file")
	ErrorReadingFile   = fmt.Errorf("Error reading file")
	ErrorTakingInput   = fmt.Errorf("Error taking input")
)
var file *os.File = nil
var cols [][]byte = [][]byte{}
var scanner *bufio.Scanner

func createTable(scanner *bufio.Scanner) *os.File {

	fmt.Println("Enter table name")

	var tableName string

	fmt.Scanln(&tableName)

	file, err := os.Create("../data/" + tableName + ".txt")

	if err != nil {
		fmt.Println(ErrorCreatingFile, err)
	}

	fmt.Println("Table created successfully")

	fmt.Println("Enter column names separated by comma")

	var columnNames string

	// scanner := bufio.NewScanner(os.Stdin)

	if scanner.Err() != nil {
		fmt.Println("Error reading column names", scanner.Err())
	}

	if scanner.Scan() {
		columnNames = scanner.Text()
		columnNames = strings.ReplaceAll(columnNames, " ", "")
	}

	if err != nil {
		fmt.Println("Error reading column names", err)
	}

	fmt.Println(columnNames)

	_, err = file.WriteString(columnNames + "\n")

	if err != nil {
		fmt.Println(ErrorWritingToFile, err)
	}

	fmt.Println("Column names added successfully")
	return file

}

func useTable(name string) *os.File {

	file, err := os.OpenFile("../data/"+name+".txt", os.O_APPEND|os.O_RDWR, 0644)

	if err != nil {
		fmt.Println(ErrorReadingFile, err)
		return nil
	}

	return file

}

func insertRows(scanner *bufio.Scanner) {

	if file == nil {
		fmt.Println(ErrorReadingFile)
		return
	}

	if scanner.Err() != nil {
		fmt.Println("Error reading data rows", scanner.Err())
	}

	fmt.Println("Enter data rows seperated by comma  (type 'exit' to finish)")
	for scanner.Scan() {

		text := scanner.Text()

		if text == "exit" {
			return
		}

		_, err := file.WriteString(text + "\n")

		if err != nil {
			fmt.Println(ErrorWritingToFile, err)
			return
		}

		fmt.Println("Data row added successfully")

	}
}

func selectRows() {
	file.Seek(0, 0)

	scanner := bufio.NewScanner(file)

	if scanner.Err() != nil {
		fmt.Println(ErrorReadingFile, scanner.Err())
		return
	}

	for scanner.Scan() {
		if scanner.Err() != nil {
			fmt.Println(ErrorReadingFile, scanner.Err())
			return
		}
		fmt.Println(scanner.Text())
	}

}

func deleteRows() {

	file.Truncate(0)
	file.Seek(0, 0)
	file.WriteString(string(bytes.Join(cols, []byte(","))) + "\n")
}

func main() {

	fmt.Println("DB is starting")

	fmt.Println("Enter 'create' to create a new table or 'use' to use an existing table")

	var command string

	scanner = bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		command = scanner.Text()
		if scanner.Err() != nil {
			fmt.Println(ErrorTakingInput, scanner.Err())
		}

		if command != "create" && command != "use" {
			fmt.Println("Invalid command")
			return
		}
	}

	if command == "create" {
		file = createTable(scanner)
	} else {
		fmt.Println("Enter table name to use")
		scanner.Scan()
		if scanner.Err() != nil {
			fmt.Println(ErrorTakingInput, scanner.Err())
		}
		tableName := scanner.Text()
		file = useTable(tableName)
	}

	if file == nil {
		fmt.Println("Exiting due to error")
		return
	}

	file.Seek(0, 0)

	fileScanner := bufio.NewScanner(file)
	if fileScanner.Scan() {
		colLine := fileScanner.Bytes()
		cols = bytes.Split(colLine, []byte(","))
	}
	if fileScanner.Err() != nil {
		fmt.Println(ErrorReadingFile, fileScanner.Err())
		return
	}

	defer file.Close()

	if scanner.Err() != nil {
		fmt.Println(ErrorTakingInput, scanner.Err())
	}

	fmt.Println("Enter 'insert' or 'delete' or 'select' or 'exit' to perform operations on the table")
	for scanner.Scan() {

		if scanner.Err() != nil {
			fmt.Println(ErrorTakingInput, scanner.Err())
		}
		command = scanner.Text()

		if command == "exit" {
			fmt.Println("Exiting...")
			return
		}

		if command != "insert" && command != "delete" && command != "select" {
			fmt.Println("Invalid command")
			return
		}

		switch command {
		case "insert":
			insertRows(scanner)
		case "select":
			selectRows()
		case "delete":
			deleteRows()
		}
		fmt.Println("Enter 'insert' or 'delete' or 'select' or 'exit' to perform operations on the table")
	}

}
