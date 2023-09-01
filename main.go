package main

import (
	"fmt"
	"strconv"
)

func main() {
	str := "(-1+(4+5+2)-3)+(6+8)"
	result := calculate(str)

	fmt.Println(result)
}

func calculate(s string) int {
	tokens := getTokens(s)
	list := TokenList{tokens, 0}
	result, ok := parseExpression(&list)
	if !ok {
		panic("AAAA")
	}

	return result
}

func parseExpression(tokens *TokenList) (int, bool) {
	return parseAddSub(tokens)
}

func parseAddSub(tokens *TokenList) (int, bool) {
	left, ok := parseMulDiv(tokens)
	if !ok {
		return 0, false
	}

	for {
		current, ok := tokens.current()
		if !ok {
			return left, true
		}

		if current.kind != Add && current.kind != Sub {
			return left, true
		}

		tokens.advance()

		right, ok := parseMulDiv(tokens)
		if !ok {
			panic("add/sub invalid")
		}

		if current.kind == Add {
			left += right
		} else if current.kind == Sub {
			left -= right
		} else {
			panic("AAA")
		}
	}
}

func parseMulDiv(tokens *TokenList) (int, bool) {
	left, ok := parseUnaryMinus(tokens)
	if !ok {
		return 0, false
	}

	for {
		current, ok := tokens.current()
		if !ok {
			return left, true
		}

		if current.kind != Mul && current.kind != Div {
			return left, true
		}

		tokens.advance()

		right, ok := parseUnaryMinus(tokens)
		if !ok {
			panic("mul/div invalid")
		}

		if current.kind == Mul {
			left *= right
		} else if current.kind == Div {
			left /= right
		} else {
			panic("AAA")
		}
	}
}

func parseUnaryMinus(tokens *TokenList) (int, bool) {
	current, ok := tokens.current()
	if !ok {
		return 0, false
	}

	isUnary := current.kind == UnaryMinus
	if isUnary {
		tokens.advance()
	}

	operand, ok := parseOperand(tokens)
	if !ok {
		panic("AAA")
	}

	if isUnary {
		return -operand, true
	}

	return operand, true
}

func parseOperand(tokens *TokenList) (int, bool) {
	operand, ok := parseNumber(tokens)
	if !ok {
		operand, ok = parseBrackets(tokens)
	}

	if !ok {
		panic("Can't parse operand")
	}

	return operand, true
}

func parseBrackets(tokens *TokenList) (int, bool) {
	current, ok := tokens.current()
	if !ok || current.kind != OpenBracket {
		return 0, false
	}

	tokens.advance()

	expression, ok := parseExpression(tokens)
	if !ok {
		return 0, false
	}

	current, ok = tokens.current()
	if !ok || current.kind != CloseBracket {
		panic("CloseBracket is missing")
	}

	tokens.advance()

	return expression, true
}

func parseNumber(tokens *TokenList) (int, bool) {
	current, ok := tokens.current()
	if !ok || current.kind != Number {
		return 0, false
	}

	tokens.advance()

	return current.number, true
}

func getTokens(s string) []Token {
	result := make([]Token, len(s))
	index := 0

	for i := 0; i < len(s); i++ {
		el := s[i]

		if el == '+' {
			result[index] = Token{kind: Add}
		} else if el == '-' {
			if index == 0 || result[index-1].kind == OpenBracket || result[index-1].kind == Add || result[index-1].kind == Sub || result[index-1].kind == Mul || result[index-1].kind == Div {
				result[index] = Token{kind: UnaryMinus}
			} else {
				result[index] = Token{kind: Sub}
			}
		} else if el == '*' {
			result[index] = Token{kind: Mul}
		} else if el == '/' {
			result[index] = Token{kind: Div}
		} else if el == '(' {
			result[index] = Token{kind: OpenBracket}
		} else if el == ')' {
			result[index] = Token{kind: CloseBracket}
		} else if isDigit(el) {
			startIndex := i
			endIndex := i + 1

			for ; endIndex < len(s) && isDigit(s[endIndex]); endIndex++ {
				i++
			}

			number, _ := strconv.Atoi(s[startIndex:endIndex])
			result[index] = Token{kind: Number, number: number}
		} else {
			continue
		}

		index++
	}

	return result[:index]
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func getPriority(token Token) int {
	if token.kind == Add || token.kind == Sub {
		return 1
	}

	if token.kind == Mul || token.kind == Div {
		return 2
	}

	if token.kind == UnaryMinus {
		return 3
	}

	panic("Invalid token")
}

type Token struct {
	kind   TokenKind
	number int
}

type TokenKind int

const (
	Number TokenKind = iota
	Add
	Sub
	Mul
	Div
	UnaryMinus
	OpenBracket
	CloseBracket
)

type TokenList struct {
	tokens []Token
	index  int
}

func (list TokenList) current() (*Token, bool) {
	if list.index < 0 || list.index >= len(list.tokens) {
		return nil, false
	}

	return &list.tokens[list.index], true
}

func (list *TokenList) advance() {
	list.index++
}
