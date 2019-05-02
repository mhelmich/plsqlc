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

func parseInsideBlock(p *parser, pc *parserContext) {
	pkg := pc.pkg
	f := pc.function
	blk := pc.block
	for {
		switch i := p.next(); i.Typ {
		case lexer.IdentifierType:
			// could be a qualified function call ('package.func()'), a local function call ('func()') or an assignment ('a:=12')
			if p.acceptValue(".") {
				fc := parseQualifiedFunctionCall(p, i.Value)
				blk.AddInstruction(fc)
				continue

			} else if p.acceptValue("(") {
				fc := parseLocalFunctionCall(p, pkg.Name, i.Value)
				blk.AddInstruction(fc)
				continue

			} else if p.acceptValue(":=") {
				a := parseAssignment(p, i.Value)
				blk.AddInstruction(a)
				continue
			}

			log.Panicf("Shouldn't reach here!")

		case lexer.KeywordType:
			switch i.Value {
			case "END":
				// eat a potential 'IF'
				p.acceptValue("IF")
				// eat a potential 'LOOP'
				p.acceptValue("LOOP")

				if ok := p.acceptValue(";"); !ok {
					log.Panicf("Can't find ';' lex item")
				}
				pc.block = nil
				return

			case "IF":
				cond := parseBinOp(p)
				if ok := p.acceptValue("THEN"); !ok {
					log.Panicf("Can't find 'THEN' lex item")
				}

				// put if block into context (the old one is still saved in 'blk')
				ifBlk := ast.NewBlock("if-block")
				pc.block = ifBlk
				f.AddBlock(ifBlk)

				parseInsideBlock(p, pc)
				// create new merge block for after if branches
				// and prime pointers
				mergeBlk := ast.NewBlock("merge-block")
				pc.block = mergeBlk
				f.AddBlock(mergeBlk)

				ifBlk.Terminator = ast.NewBranch(mergeBlk)
				blk.Terminator = ast.NewConditionalBranch(cond, ifBlk, mergeBlk)

				blk = mergeBlk
				continue

			case "WHILE":
				cond := parseBinOp(p)
				if ok := p.acceptValue("LOOP"); !ok {
					log.Panicf("Can't find 'LOOP' lex item")
				}

				loopBlk := ast.NewBlock("loop-block")
				pc.block = loopBlk
				f.AddBlock(loopBlk)

				parseInsideBlock(p, pc)

				mergeBlk := ast.NewBlock("merge-block")
				f.AddBlock(mergeBlk)
				blk.Terminator = ast.NewConditionalBranch(cond, loopBlk, mergeBlk)
				loopBlk.Terminator = ast.NewConditionalBranch(cond, loopBlk, mergeBlk)

				blk = mergeBlk
				continue

			default:
				log.Panicf("Can't match lex item '%s'", i.Value)
			}
		default:
			log.Panicf("Can't match lex item '%s'", i.Value)
		}
	}
}

func parseBinOp(p *parser) *ast.BinOp {
	return parseBinOpExpression(p, p.next())
}

// an expression can be:
// - a function call
// - a variable
// - string
// - a number
// - a binop (todo)
func parseExpression(p *parser) ast.Expression {
	return parseExpressionFromLexItem(p, p.next())
}

func parseExpressionFromLexItem(p *parser, i *lexer.Item) ast.Expression {
	switch i.Typ {
	case lexer.StringType:
		return ast.NewStringLiteral(i.Value)

	case lexer.NumericType:
		return ast.NewNumericLiteral(i.Value)

	case lexer.IdentifierType:
		// this could be a function call or a variable
		if p.peek().Value == "(" {
			// function call
			log.Panicf("not implemented yet")
		} else if p.peek().Typ == lexer.OperatorType {
			// bin op
			return parseBinOpExpression(p, i)
		} else {
			// variable
			return ast.NewVariable(i.Value)
		}
		log.Panicf("not implemented yet")

	default:
		log.Panicf("Can't match lex item '%s'", i.Value)
	}

	return nil
}

func parseBinOpExpression(p *parser, leftItem *lexer.Item) *ast.BinOp {
	// eat all the items from the parser before parsing them
	opItem := p.next()
	rightItem := p.next()

	if opItem.Typ != lexer.OperatorType {
		log.Panicf("Lex item '%s' is not an operator\n", opItem.Value)
	}

	left := parseExpressionFromLexItem(p, leftItem)
	right := parseExpressionFromLexItem(p, rightItem)
	return ast.NewBinOp(left, opItem.Value, right)
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

func parseLocalFunctionCall(p *parser, moduleName string, funcName string) ast.Expression {
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

func parseAssignment(p *parser, identifier string) *ast.Assignment {
	expr := parseExpression(p)
	a := ast.NewAssignment(identifier, expr)

	if ok := p.acceptValue(";"); !ok {
		log.Panicf("Can't find ';' lex item")
	}
	return a
}
