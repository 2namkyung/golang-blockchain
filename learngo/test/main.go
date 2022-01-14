package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("list.txt")
	if err != nil {
		panic(err)
	}

	output, err := os.Create("output.txt")
	if err != nil {
		panic(err)
	}

	defer file.Close()
	defer output.Close()

	reader := bufio.NewReader(file)
	writer := bufio.NewWriter(output)

	for {
		line, isPrefix, err := reader.ReadLine()
		if isPrefix || err != nil {
			break
		}
		words := strings.Fields(string(line))

		if strings.Contains(words[1], "&") {
			words[1] = strings.ReplaceAll(words[1], "&", "\\&")
		}

		stmt := "insert into companyList(ticker, compName) values('" + words[0] + "','" + words[1] + "');\n"
		t, err := writer.WriteString(stmt)
		writer.Flush()
		fmt.Println(t)
	}
}
