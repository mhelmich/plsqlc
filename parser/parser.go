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
	"fmt"
	"log"

	"github.com/mhelmich/plsqlc/ast"
	"github.com/mhelmich/plsqlc/lexer"
)

type parser struct {
	input        <-chan *lexer.Item
	packages     map[string]*ast.Package
	peekableItem *lexer.Item
	sem          chan interface{}
}

func NewParser(input <-chan *lexer.Item) *parser {
	p := newParser(input)
	go p.run()
	return p
}

func newParser(input <-chan *lexer.Item) *parser {
	p := &parser{
		input:    input,
		packages: make(map[string]*ast.Package),
		sem:      make(chan interface{}),
	}

	return p
}

func (p *parser) GetPackageAsts() map[string]*ast.Package {
	if p.sem != nil {
		<-p.sem
		p.sem = nil
	}

	return p.packages
}

func (p *parser) run() {
	state := parseText
	pc := &parserContext{}
	for state != nil {
		state, pc = state(p, pc)
	}

	close(p.sem)
}

func (p *parser) next() *lexer.Item {
	if p.peekableItem != nil {
		tmp := p.peekableItem
		p.peekableItem = nil
		return tmp
	}
	return <-p.input
}

func (p *parser) peek() *lexer.Item {
	if p.peekableItem == nil {
		p.peekableItem = <-p.input
	}
	return p.peekableItem
}

func (p *parser) addPackage(pkg *ast.Package) {
	p.packages[pkg.Name] = pkg
}

func (p *parser) acceptValue(valid string) bool {
	if p.peek().Value == valid {
		p.next()
		return true
	}
	return false
}

func (p *parser) acceptType(valid lexer.ItemType) (bool, string) {
	if p.peek().Typ == valid {
		i := p.next()
		return true, i.Value
	}
	return false, ""
}

func (p *parser) errorf(format string, args ...interface{}) stateFunc {
	log.Printf(format+"\n", args...)
	return nil
}

type parserContext struct {
	pkg      *ast.Package
	function *ast.Function
	block    *ast.Block
}

func (pc *parserContext) String() string {
	return fmt.Sprintf("pkg: %s, func: %s, blk: %s", pc.pkg.Name, pc.function.Proto.Name, pc.block.Name)
}
