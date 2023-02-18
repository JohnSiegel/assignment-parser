package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type Token struct {
	Kind  string
	Value string
}

type Parser struct {
	tokens    []Token
	pos       int
	variables map[string]int
}

func main() {

	input := os.Args[1:]
	if len(input) != 1 {
		println("Usage: BinOps inputFile.txt")
		os.Exit(1)
	}
	assignments, err := readLines(input[0])
	if err != nil {
		fmt.Printf("Error with reading input file: %s", err)
	}
	CalcAllAndPrint(assignments)
}

func isDigit(c byte) bool     { return c >= '0' && c <= '9' }
func isLowerCase(c byte) bool { return c >= 'a' && c <= 'z' }
func isUpperCase(c byte) bool { return c >= 'A' && c <= 'Z' }
func isLetter(c byte) bool {
	return isLowerCase(c) || isUpperCase(c)
}

func isLetterOrDigit(c byte) bool {
	return isDigit(c) || isLetter(c)
}

func CalcAllAndPrint(exprs []string) {
	fmt.Println("The read expressions are:")
	variables := make(map[string]int)
	for i, line := range exprs {
		fmt.Println(line)
		p := NewParser(tokenize(line), variables)
		p.Parse(i == len(exprs)-1)
	}
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func NewParser(tokens []Token, variables map[string]int) *Parser {
	return &Parser{
		tokens:    tokens,
		pos:       0,
		variables: variables,
	}
}

func (p *Parser) Parse(shouldPrint bool) {
	result := p.parseAssignment()

	// Parse all the assignments
	for p.pos < len(p.tokens) {
		result = p.parseAssignment()
	}

	// Print the value of the last assignment
	if shouldPrint {
		fmt.Println(p.variables[result])
	}
}

func (p *Parser) parseAssignment() string {
	// The left-hand side of the assignment is a variable
	variable := p.expect("var").Value

	// Check that the assignment operator is present
	p.expect("=")

	// Parse the right-hand side of the assignment
	value := p.parseExpression()

	// Store the value in the variable
	p.variables[variable] = value

	return variable
}

func (p *Parser) parseExpression() int {
	// Parse the first term
	value := p.parseTerm()

	// Parse any additional terms with bitwise OR operators
	for p.accept("|") {
		value |= p.parseTerm()
	}

	return value
}

func (p *Parser) parseTerm() int {
	// Parse the first factor
	value := p.parseFactor()

	// Parse any additional factors with bitwise XOR operators
	for p.accept("^") {
		value ^= p.parseFactor()
	}

	return value
}

func (p *Parser) parseFactor() int {
	// Parse the first unary expression
	value := p.parseUnary()

	// Parse any additional unary expressions with bitwise AND operators
	for p.accept("&") {
		value &= p.parseUnary()
	}

	return value
}

func (p *Parser) parseUnary() int {
	if p.accept("~") {
		// Unary NOT operator
		return ^p.parseUnary()
	} else if p.accept("var") {
		// Variable reference
		variable := p.tokens[p.pos-1].Value
		if value, ok := p.variables[variable]; ok {
			return value
		} else {
			panic(fmt.Sprintf("Undefined variable %s", variable))
		}
	} else if p.accept("const") {
		// Constant value
		value, err := strconv.Atoi(p.tokens[p.pos-1].Value)
		if err != nil {
			panic(err)
		}
		return value
	} else if p.accept("(") {
		// Parenthesized expression
		value := p.parseExpression()
		p.expect(")")
		return value
	} else {
		// Syntax error
		panic("Syntax error")
	}
}

func (p *Parser) accept(kind string) bool {
	if p.pos < len(p.tokens) && p.tokens[p.pos].Kind == kind {
		p.pos++
		return true
	} else {
		return false
	}
}

func (p *Parser) expect(kind string) Token {
	if p.pos < len(p.tokens) && p.tokens[p.pos].Kind == kind {
		token := p.tokens[p.pos]
		p.pos++
		return token
	} else {
		panic(fmt.Sprintf("Expected %s", kind))
	}
}

func tokenize(s string) []Token {
	var tokens []Token
	var token Token
	var value string
	var i int

	for i < len(s) {
		// Skip whitespace
		for i < len(s) && (s[i] == ' ' || s[i] == '\t') {
			i++
		}

		// Check for end of input
		if i >= len(s) {
			break
		}

		// Check for a variable
		if isLetter(s[i]) {
			value = string(s[i])
			i++
			for i < len(s) && isLetterOrDigit(s[i]) {
				value += string(s[i])
				i++
			}
			token = Token{Kind: "var", Value: value}
		} else if isDigit(s[i]) {
			value = string(s[i])
			i++
			for i < len(s) && isDigit(s[i]) {
				value += string(s[i])
				i++
			}
			token = Token{Kind: "const", Value: value}
		} else {
			// Check for an operator
			switch s[i] {
			case '=':
				token = Token{Kind: "=", Value: "="}
			case '|':
				token = Token{Kind: "|", Value: "|"}
			case '^':
				token = Token{Kind: "^", Value: "^"}
			case '&':
				token = Token{Kind: "&", Value: "&"}
			case '~':
				token = Token{Kind: "~", Value: "~"}
			case '(':
				token = Token{Kind: "(", Value: "("}
			case ')':
				token = Token{Kind: ")", Value: ")"}
			default:
				panic(fmt.Sprintf("Invalid character %c", s[i]))
			}
			i++
		}

		tokens = append(tokens, token)
	}

	return tokens
}
