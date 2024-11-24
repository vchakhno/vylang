package main

import (
	"fmt"
	"strings"
	"errors"
)

type Symbole interface {
	// ~Atom | ~SExpr
}

type Atom = string

type SExpr struct {
	children []Symbole
}

type Parser struct {
	str []byte
	index int
}

func isSpace(c byte) bool {
	return c == ' ' || ('\t' <= c && c <= '\r' )
}

func (p *Parser) skipSpace() {
	for p.index < len(p.str) && isSpace(p.str[p.index]) {
		p.index++
	}
}

func (p *Parser) atom() (Atom) {
	var storage strings.Builder;

	for p.index < len(p.str) && !isSpace(p.str[p.index]) && p.str[p.index] != '(' &&  p.str[p.index] != ')' {
		char := p.str[p.index]
		storage.WriteByte(char)
		p.index++
	}
	return storage.String()
}

func (p *Parser) sexpr() (SExpr, error) {
	var output SExpr;

	if (p.str[p.index] != '(') {
		fmt.Printf("%s %d [%c]\n", string(p.str), p.index, p.str[p.index - 1])
		return SExpr{}, errors.New("Missing open parenthesis")
	}
	p.index++
	p.skipSpace()
	for p.index < len(p.str) {
		switch p.str[p.index]{
		case '(':
			sexpr, err := p.sexpr()
			if err != nil {
				return SExpr{}, err;
			}
			output.children = append(output.children, sexpr)
			break;
		case ')':
			p.index++
			return output, nil
		default:
			atom := p.atom();
			output.children = append(output.children, atom)
			break;
		}
		p.skipSpace()
	}
	return SExpr{}, errors.New("Missing close parenthesis")
}

func (p *Parser) Parse() ([]SExpr, error) {
	var output []SExpr;

	p.skipSpace()
	for p.index < len(p.str) {
		sexpr, err := p.sexpr()
		if err != nil {
			return []SExpr{}, err;
		}
		output = append(output, sexpr)
		p.skipSpace()
	}
	return output, nil
}




func main() {
	input := `(abc (def) ghi)
		()
		(aaaa)
		()
		(banane)
	`

	parser := Parser {[]byte(input), 0};
	ast, err := parser.Parse()
	if err != nil {
		panic(err)
	}
	// print ast

	fmt.Println("Resulting list:", len(ast))
}
