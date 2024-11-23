package main
import (
	"errors"
	"fmt"
	"os"
	"bytes"
)


func parseStringLit(source string, start_pos int) (literal string, end_pos int, err error) {
	if source[0] != '"' {
		return "", 0, errors.New("Error: expected string to start with double quotes")
	}
	var literal_buf bytes.Buffer
	i := 1
	for i < len(source) - 1 {
		if source[i] == '"' {
			break
		}
		if source[i] < ' ' || source[i] == '\\' || source[i] > '~' {
			return "", 0, fmt.Errorf("Error: unexpected character in string: '%c'", source[i])
		}
		literal_buf.WriteByte(source[i])
		i++
	}
	if i == len(source) {
		return "", 0, errors.New("Error: expected string to end with double quotes")
	}
	return literal_buf.String(), i + 1, nil
}

type Program struct {
	message string
}

func parseProgram(source string) (Program, error) {
	message, end_pos, err := parseStringLit(source, 0)
	if err != nil {
		return Program {message: ""}, err
	}
	if end_pos != len(source) {
		return Program {message: ""}, fmt.Errorf("Error: Unexpected character after end of program")
	}
	return Program {message: message}, nil
}

func main() {
	source_bytes, err := os.ReadFile("main.vy")
	if err != nil {
		panic(err)
	}
	source := string(source_bytes)

	program, err := parseProgram(source)
	if err != nil {
		panic(err)
	}
	message := program.message

	llvmir := fmt.Sprintf(
`@.str = private unnamed_addr constant [%[2]d x i8] c"%[1]s\00", align 1

define dso_local i32 @main() #0 {
	%%1 = call i64 @write(i32 1, i8* getelementptr inbounds ([%[2]d x i8], [%[2]d x i8]* @.str, i64 0, i64 0), i64 %[3]d)
	ret i32 0
}

declare dso_local i64 @write(i32, i8*, i64) #1`, message, len(message) + 1, len(message))


	llvmir_bytes := []byte(llvmir)
	err = os.WriteFile("bin/mir.ll", llvmir_bytes, 0600)
	if err != nil {
		panic(err)
	}

	fmt.Println("Vy file compiled!")
}