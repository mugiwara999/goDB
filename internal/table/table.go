package table

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Table struct {
	file *os.File
	cols []string
}

var (
	ErrorCreatingFile  = fmt.Errorf("Error creating file")
	ErrorWritingToFile = fmt.Errorf("Error writing to file")
	ErrorReadingFile   = fmt.Errorf("Error reading file")
	ErrorTakingInput   = fmt.Errorf("Error taking input")
	ErrorInvalidInput  = fmt.Errorf("Invalid input")
)

func Open(name string) (*Table, error) {

	err := godotenv.Load()
	var DataDir string

	if err != nil {
		DataDir = "../data"
	} else {
		DataDir = os.Getenv("DATA_DIR")
	}
	file, err := os.OpenFile(DataDir+"/"+name+".txt", os.O_APPEND|os.O_RDWR, 0644)

	if err != nil {
		return nil, err
	}

	fileScanner := bufio.NewScanner(file)

	fileScanner.Scan()

	cols := strings.Split(fileScanner.Text(), ",")

	table := &Table{
		file: file,
		cols: cols,
	}

	return table, nil

}
func Create(name string, cols []string) (*Table, error) {

	err := godotenv.Load()
	var DataDir string

	if err != nil {
		DataDir = "../data"
	} else {
		DataDir = os.Getenv("DATA_DIR")
	}
	file, err := os.Create(DataDir + "/" + name + ".txt")

	if err != nil {
		return nil, fmt.Errorf(ErrorCreatingFile.Error(), err)
	}

	_, err = file.WriteString(strings.Join(cols, ",") + "\n")

	if err != nil {
		fmt.Println(ErrorWritingToFile, err)
	}

	table := &Table{
		file: file,
		cols: cols,
	}
	return table, nil
}

func (t *Table) Close() error {

	return t.file.Close()

}
