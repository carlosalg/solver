package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strconv"
)

func main() {
	//Check for arguments in the console
	if len(os.Args) < 2 {
		fmt.Println("Please provide a file name.")
		return
	}
	//take the name of the file from the argument in the console
	fileName := os.Args[1]
	//open the fiel and check for errors
	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("failed to open file:", err)
	}
	//ensure that the file is closed when the program finish
	defer file.Close()
	//read the file line by line calls the solver funtion an check errors
	scanner := bufio.NewScanner(file)
	var funtions []string
	var prettyFormatWritten bool

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}
		funtions = append(funtions, line)
		expAst, err := parser.ParseExpr(line)
		if err != nil {
			fmt.Println("Error parsing expresion: ", err)
			return
		}
		result, err := expSolver(expAst)
		if err != nil {
			fmt.Println("Error evaluating expresion:", err)
			return
		}
		fmt.Println("Result:", result)
		if !prettyFormatWritten {
			_, err = file.WriteString("------+Result+------\n")
			if err != nil {
				panic(err)
			}
			prettyFormatWritten = true
		}
		_, err = file.WriteString(strconv.FormatFloat(result, 'f', -1, 64) + "\n")
		if err != nil {
			panic(err)
		}
	}
	if prettyFormatWritten {
		_, err = file.WriteString("---------------------\n")
		if err != nil {
			panic(err)
		}
	}
	//fmt.Println(funtions[0])
	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

// Solve each of the expresions using AST
func expSolver(expr ast.Expr) (float64, error) {
	switch e := expr.(type) {
	case *ast.BinaryExpr:
		left, err := expSolver(e.X)
		if err != nil {
			return 0, err
		}
		right, err := expSolver(e.Y)
		if err != nil {
			return 0, err
		}
		switch e.Op {
		case token.ADD:
			return left + right, nil
		case token.SUB:
			return left - right, nil
		case token.MUL:
			return left * right, nil
		case token.QUO:
			return left / right, nil
		}
	case *ast.BasicLit:
		value, err := strconv.ParseFloat(e.Value, 64)
		if err != nil {
			return 0, err
		}
		return value, nil
	}
	return 0, fmt.Errorf("unsupported expresion: %T", expr)
}
