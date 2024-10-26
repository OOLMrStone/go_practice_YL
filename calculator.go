package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Calc(expression string) (float64, error) {
	tokens := tokenize(expression)
	postfix, err := toPostfix(tokens)
	if err != nil {
		return 0, err
	}
	return evalPostfix(postfix)
}

func tokenize(expression string) []string {
	var tokens []string
	var currentNum strings.Builder

	for _, ch := range expression {
		if isSpace(ch) {
			continue
		}
		if isDigit(ch) || ch == '.' {
			currentNum.WriteRune(ch)
		} else {
			if currentNum.Len() > 0 {
				tokens = append(tokens, currentNum.String())
				currentNum.Reset()
			}
			tokens = append(tokens, string(ch))
		}
	}
	if currentNum.Len() > 0 {
		tokens = append(tokens, currentNum.String())
	}

	return tokens
}

func toPostfix(tokens []string) ([]string, error) {
	var result []string
	var operators []string

	prec := map[string]int{
		"+": 1,
		"-": 1,
		"*": 2,
		"/": 2,
		"(": 0,
	}

	for _, token := range tokens {
		if isDigit(rune(token[0])) {
			result = append(result, token)
		} else if token == "(" {
			operators = append(operators, token)
		} else if token == ")" {
			for len(operators) > 0 && operators[len(operators)-1] != "(" {
				result = append(result, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			if len(operators) == 0 {
				return nil, errors.New("Неверная запись выражения")
			}
			operators = operators[:len(operators)-1]
		} else {
			for len(operators) > 0 && prec[operators[len(operators)-1]] >= prec[token] {
				result = append(result, operators[len(operators)-1])
				operators = operators[:len(operators)-1]
			}
			operators = append(operators, token)
		}
	}

	for len(operators) > 0 {
		if operators[len(operators)-1] == "(" {
			return nil, errors.New("В выражении есть незакрытые скобки")
		}
		result = append(result, operators[len(operators)-1])
		operators = operators[:len(operators)-1]
	}

	return result, nil
}

func evalPostfix(tokens []string) (float64, error) {
	stack := []float64{}

	for _, token := range tokens {
		if isDigit(rune(token[0])) {
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, err
			}
			stack = append(stack, num)
		} else {
			if len(stack) < 2 {
				return 0, errors.New("Неверная запись выражения")
			}
			b := stack[len(stack)-1]
			a := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			switch token {
			case "+":
				stack = append(stack, a+b)
			case "-":
				stack = append(stack, a-b)
			case "*":
				stack = append(stack, a*b)
			case "/":
				if b == 0 {
					return 0, errors.New("Деление на ноль")
				}
				stack = append(stack, a/b)
			default:
				return 0, errors.New("Неизвестный знак: " + token)
			}
		}
	}

	if len(stack) != 1 {
		return 0, errors.New("Неверная запись выражения")
	}

	return stack[0], nil
}

func isSpace(ch rune) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

func isDigit(ch rune) bool {
	return (ch >= '0' && ch <= '9') || ch == '.'
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Введите математическое выражение: ")
	expression, _ := reader.ReadString('\n')
	expression = strings.TrimSpace(expression)

	result, err := Calc(expression)
	if err != nil {
		fmt.Println("Ошибка:", err)
	} else {
		fmt.Printf("Результат: %v\n", result)
	}
}
