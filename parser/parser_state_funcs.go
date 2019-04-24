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

type stateFunc func(*parser, *parserContext) (stateFunc, *parserContext)

func parseText(p *parser, pc *parserContext) (stateFunc, *parserContext) {
	switch i := p.next(); i.Value {
	case "CREATE":
		if ok := p.acceptValue("OR"); !ok {
			log.Panicf("Can't find 'or' lex item instead got %s", p.peek().Value)
		}
		if ok := p.acceptValue("REPLACE"); !ok {
			log.Panicf("Can't find 'replace' lex item instead got %s", p.peek().Value)
		}
		switch i2 := p.next(); i2.Value {
		case "PACKAGE":
			return parseCreatePackage, pc
		default:
			log.Panicf("Can't match item %s", i.String())
		}

	default:
		log.Panicf("Can't match item %s", i.String())
	}
	return nil, nil
}

func parseCreatePackage(p *parser, pc *parserContext) (stateFunc, *parserContext) {
	switch i := p.next(); i.Value {
	case "BODY":
		packageNameItem := p.next()
		if ok := p.acceptValue("AS"); !ok {
			log.Panicf("Can't find 'as' lex item")
		}
		pkg := ast.NewPackage(packageNameItem.Value)
		p.packages[packageNameItem.Value] = pkg
		log.Printf("Found package: %s\n", pkg.Name)
		pc.pkg = pkg
		return parseInsidePackage, pc

	default:
		log.Panicf("Can't match lex item '%s'", i.String())
	}

	return nil, nil
}

func parseInsidePackage(p *parser, pc *parserContext) (stateFunc, *parserContext) {
	pkg := pc.pkg
	switch i := p.next(); i.Value {
	case "PROCEDURE":
		fName := p.next().Value
		f := ast.NewFunction(fName, true)
		pkg.AddFunction(f)
		pc.function = f
		return parseFunction, pc

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

func parseFunction(p *parser, pc *parserContext) (stateFunc, *parserContext) {
	f := pc.function
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
		return parseFunctionBody, pc

	case lexer.KeywordType:
		if i.Value != "IS" {
			log.Panicf("Can't find 'is' lex item but is %s", p.next().Value)
		}

		// parse function locals
		for p.peek().Typ == lexer.IdentifierType {
			localName := p.next().Value
			localType := p.next().Value
			p.next()
			localValue := p.next().Value
			p.next()
			f.AddLocal(localName, localType, localValue)
		}

		return parseFunctionBody, pc

	default:
		log.Panicf("Can't match token %s", i.Value)
	}

	pc.function = nil
	return parseInsidePackage, pc
}

func parseFunctionBody(p *parser, pc *parserContext) (stateFunc, *parserContext) {
	f := pc.function
	switch i := p.next(); i.Value {
	case "BEGIN":
		blk := ast.NewBlock(f.Proto.Name + "-entry")
		f.AddBlock(blk)
		pc.block = blk
		parseInsideBlock(p, pc)

	default:
		log.Panicf("Can't match lex item '%s'", i.Value)
	}
	log.Printf("parsed function body for %s\n", f.Proto.Name)
	pc.function = nil
	return parseInsidePackage, pc
}
