package forth

import (
	"fmt"
	"strconv"
	"strings"
)

var wordsMap map[string]string

func Forth(input []string) ([]int, error) {
	initializeOperationsMap()

	inputString := strings.ToUpper(input[len(input)-1])

	for i := 0; i < len(input)-1; i++ {
		err := handleNewWordDefinition(strings.ToUpper(input[i]))
		if err != nil {
			return nil, err
		}
	}

	updatedInputString := replaceInputStackWithNewOperations(inputString)

	result, err := solveInputStackForResult(updatedInputString)
	if err != nil {
		return nil, err
	}

	return result, nil
}

//Initialize a map which contains all the words and their respective operations
func initializeOperationsMap() {
	wordsMap = make(map[string]string)

	wordsMap["+"] = "+"
	wordsMap["-"] = "-"
	wordsMap["*"] = "*"
	wordsMap["/"] = "/"

	wordsMap["DUP"] = "DUP"
	wordsMap["DROP"] = "DROP"
	wordsMap["SWAP"] = "SWAP"
	wordsMap["OVER"] = "OVER"
}

//Function which takes word definition as an input, define that word based on the given operations and then add them into the wordsMap
func handleNewWordDefinition(wordDefinition string) error {
	sd := strings.Split(wordDefinition, " ")

	if len(sd) < 4 || sd[0] != ":" || sd[len(sd)-1] != ";" {
		return fmt.Errorf("syntax error - word definition is not following the proper syntax")
	}

	newWord := sd[1]
	_, err := strconv.Atoi(newWord)
	if err == nil {
		return fmt.Errorf("syntax error - word to be defined can not be a number")
	}

	operation, err := solveWordOperations(sd[2 : len(sd)-1])
	if err != nil {
		return err
	}

	wordsMap[newWord] = operation

	return nil
}

//Function which takes the definition of new word and then returns the updated operations string
func solveWordOperations(operations []string) (string, error) {
	operationsString := ""
	for _, v := range operations {
		_, err := strconv.Atoi(v)
		if err == nil {
			if operationsString == "" {
				operationsString = fmt.Sprintf("%v", v)
			} else {
				operationsString = fmt.Sprintf("%v %v", operationsString, v)
			}
			continue
		}

		if word, ok := wordsMap[v]; ok {
			if operationsString == "" {
				operationsString = fmt.Sprintf("%v", word)
			} else {
				operationsString = fmt.Sprintf("%v %v", operationsString, word)
			}
		} else {
			return "", fmt.Errorf("syntax error - given word is not defined: %v", v)
		}
	}

	return operationsString, nil
}

//Function which takes solveable input string as an input and then returns the updated input string
//in which all the new user-defined words are replaced with their operations
func replaceInputStackWithNewOperations(inputString string) string {
	s := strings.Split(inputString, " ")

	updatedInputString := ""
	for _, v := range s {
		if word, ok := wordsMap[v]; ok {
			if updatedInputString == "" {
				updatedInputString = fmt.Sprintf("%v", word)
			} else {
				updatedInputString = fmt.Sprintf("%v %v", updatedInputString, word)
			}
		} else {
			if updatedInputString == "" {
				updatedInputString = fmt.Sprintf("%v", v)
			} else {
				updatedInputString = fmt.Sprintf("%v %v", updatedInputString, v)
			}
		}
	}

	return updatedInputString
}

//Function which solves the given input string with the operations and then return the result as integer slice or an error if stack not in correct format
func solveInputStackForResult(inputString string) ([]int, error) {
	inputStack := strings.Split(inputString, " ")
	result := make([]int, 0)

	for i := 0; i < len(inputStack); i++ {
		v, err := strconv.Atoi(inputStack[i])
		if err == nil {
			result = append(result, v)
			continue
		}

		result, err = doOperation(inputStack[i], result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

//Function which performs a given operation on the given result slice and then returns the updated result slice or an error
func doOperation(op string, result []int) ([]int, error) {

	if op == "+" || op == "-" || op == "*" || op == "/" || op == "SWAP" || op == "OVER" {
		if len(result) < 2 {
			return nil, fmt.Errorf("error - not enough integers found to perform the given operation")
		}

		if op == "+" {
			result = append(result[:len(result)-2], result[len(result)-2]+result[len(result)-1])
		} else if op == "-" {
			result = append(result[:len(result)-2], result[len(result)-2]-result[len(result)-1])
		} else if op == "*" {
			result = append(result[:len(result)-2], result[len(result)-2]*result[len(result)-1])
		} else if op == "/" {
			if result[len(result)-1] == 0 {
				return nil, fmt.Errorf("error - can not divide by zero")
			}
			result = append(result[:len(result)-2], result[len(result)-2]/result[len(result)-1])
		} else if op == "SWAP" {
			e1 := result[len(result)-2]
			e2 := result[len(result)-1]
			result[len(result)-2] = e2
			result[len(result)-1] = e1
		} else if op == "OVER" {
			result = append(result, result[len(result)-2])
		}

	} else if op == "DUP" || op == "DROP" {
		if len(result) < 1 {
			return nil, fmt.Errorf("error - not enough integers found to perform the given operation")
		}

		if op == "DUP" {
			result = append(result, result[len(result)-1])
		} else if op == "DROP" {
			result = result[:len(result)-1]
		}

	} else {
		return nil, fmt.Errorf("syntax error - given word is not defined as an operation: %v", op)
	}

	return result, nil
}
