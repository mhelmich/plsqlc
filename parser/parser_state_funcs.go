/*
 * Copyright 2019 Marco Helmich
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package parser

import (
	"log"

	"github.com/mhelmich/plsqlc/ast"
	"github.com/mhelmich/plsqlc/lexer"
)

type stateFunc func(*parser, []interface{}) (stateFunc, []interface{})

func parseText(p *parser, args []interface{}) (stateFunc, []interface{}) {
	switch i := p.next(); i.Value {
	case "CREATE":
		if ok := p.acceptValue("OR"); !ok {
			log.Panicf("Can't find 'or' lex item")
		}
		if ok := p.acceptValue("REPLACE"); !ok {
			log.Panicf("Can't find 'replace' lex item")
		}
		switch i2 := p.next(); i2.Value {
		case "PACKAGE":
			return parseCreatePackage, nil
		default:
			log.Panicf("Can't match item %s", i.String())
		}

	default:
		log.Panicf("Can't match item %s", i.String())
	}
	return nil, nil
}

func parseCreatePackage(p *parser, args []interface{}) (stateFunc, []interface{}) {
	switch i := p.next(); i.Value {
	case "BODY":
		packageNameItem := p.next()
		if ok := p.acceptValue("AS"); !ok {
			log.Panicf("Can't find 'as' lex item")
		}
		pkg := ast.NewPackage(packageNameItem.Value)
		p.packages[packageNameItem.Value] = pkg
		log.Printf("Found package: %s\n", pkg.Name)
		return parseInsidePackage, []interface{}{pkg}

	default:
		log.Panicf("Can't match lex item '%s'", i.String())
	}

	return nil, nil
}

func parseInsidePackage(p *parser, args []interface{}) (stateFunc, []interface{}) {
	pkg := args[0].(*ast.Package)
	switch i := p.next(); i.Value {
	case "PROCEDURE":
		fName := p.next().Value
		f := ast.NewFunction(fName)
		pkg.AddFunction(f)
		return parseFunction, []interface{}{pkg, f}

	case "END":
		if ok := p.acceptValue(pkg.Name); !ok {
			log.Panicf("Can't find '%s' lex item", pkg.Name)
		}
		if ok := p.acceptValue(";"); !ok {
			log.Panicf("Can't find ';' lex item")
		}
		if ok := p.acceptValue("/"); !ok {
			log.Panicf("Can't find '/' lex item")
		}
		return nil, nil

	default:
		log.Panicf("Can't match lex item '%s'", i.Value)
	}
	return nil, nil
}

func parseFunction(p *parser, args []interface{}) (stateFunc, []interface{}) {
	pkg := args[0].(*ast.Package)
	f := args[1].(*ast.Function)
	switch i := p.next(); i.Typ {
	case lexer.SeparatorType:

		hasMore := true
		for hasMore {
			name := p.next().Value
			ownership := p.next().Value
			typ := p.next().Value
			f.AddParam(name, ownership, typ)
			sep := p.next().Value
			hasMore = sep == ","
		}

		if ok := p.acceptValue("IS"); !ok {
			log.Panicf("Can't find 'is' lex item")
		}
		return parseFunctionBody, []interface{}{pkg, f}

	case lexer.KeywordType:
		if i.Value != "IS" {
			log.Panicf("Can't find 'is' lex item but is %s", p.next().Value)
		}
		return parseFunctionBody, []interface{}{pkg, f}

	default:
		log.Panicf("Can't match token %s", i.Value)
	}
	return parseInsidePackage, []interface{}{pkg}
}

func parseFunctionBody(p *parser, args []interface{}) (stateFunc, []interface{}) {
	pkg := args[0].(*ast.Package)
	f := args[1].(*ast.Function)
	switch i := p.next(); i.Value {
	case "BEGIN":
		return parseBlock, []interface{}{pkg, f, f.Proto.Name + "-entry"}

	default:
		log.Panicf("Can't match lex item '%s'", i.Value)
	}
	log.Printf("parsed function body for %s\n", f.Proto.Name)
	return parseInsidePackage, []interface{}{pkg}
}

func parseBlock(p *parser, args []interface{}) (stateFunc, []interface{}) {
	pkg := args[0].(*ast.Package)
	f := args[1].(*ast.Function)
	blkName := args[2].(string)
	blk := ast.NewBlock(blkName)
	f.AddBlock(blk)
	return parseInsideBlock, []interface{}{pkg, blk}
}

func parseInsideBlock(p *parser, args []interface{}) (stateFunc, []interface{}) {
	pkg := args[0].(*ast.Package)
	blk := args[1].(*ast.Block)
	switch i := p.next(); i.Typ {
	case lexer.IdentifierType:
		// could be a qualified function call ('package.func()'), a local function call ('func()') or an assignment ('a:=12')
		if p.acceptValue(".") {
			fc := parseQualifiedFunctionCall(p, i.Value)
			blk.AddInstruction(fc)
			return parseInsideBlock, args
		} else if p.acceptValue("(") {
			fc := parseLocalFunctionCall(p, i.Value)
			blk.AddInstruction(fc)
			return parseInsideBlock, args
		} else if p.acceptValue(":=") {
			log.Panicf("not implemented yet")
		}

		log.Panicf("Shouldn't reach here!")

	case lexer.KeywordType:
		switch i.Value {
		case "END":
			if ok := p.acceptValue(";"); !ok {
				log.Panicf("Can't find ';' lex item")
			}
			return parseInsidePackage, []interface{}{pkg}

		case "IF":
			log.Panicf("not implemented yet")

		case "FOR":
			log.Panicf("not implemented yet")

		default:
			log.Panicf("Can't match lex item '%s'", i.Value)
		}
	default:
		log.Panicf("Can't match lex item '%s'", i.Value)
	}
	return nil, nil
}

func parseQualifiedFunctionCall(p *parser, moduleName string) ast.Expression {
	ok, funcName := p.acceptType(lexer.IdentifierType)
	if !ok {
		log.Panicf("Can't find '%s' lex item", funcName)
	}

	fc := ast.NewFunctionCall(moduleName, funcName)

	if ok := p.acceptValue("("); !ok {
		log.Panicf("Can't find '(' lex item")
	}

	for ok := p.acceptValue(")"); !ok; ok = p.acceptValue(")") {
		// there is moa parameterz
		expr := parseExpression(p)
		fc.AddArg(expr)
		p.acceptValue(",")
	}

	if ok := p.acceptValue(";"); !ok {
		log.Panicf("Can't find ';' lex item")
	}
	return fc
}

// an expression can be:
// - a function call
// - string
// - a number
func parseExpression(p *parser) ast.Expression {
	switch i := p.next(); i.Typ {
	case lexer.StringType:
		return ast.NewStringLiteral(i.Value)

	case lexer.NumericType:
		return ast.NewNumericLiteral(i.Value)

	case lexer.IdentifierType:
		log.Panicf("not implemented yet")

	default:
		log.Panicf("Can't match lex item '%s'", i.Value)
	}

	return nil
}

func parseLocalFunctionCall(p *parser, funcName string) ast.Expression {
	return nil
}

func parseAssignment(p *parser, args []interface{}) (stateFunc, []interface{}) {
	return nil, nil
}

func parseInstruction(p *parser, args []interface{}) (stateFunc, []interface{}) {
	return nil, nil
}
