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
	"reflect"
	"runtime"
	"strings"
	"testing"

	"github.com/mhelmich/plsqlc/ast"
	"github.com/mhelmich/plsqlc/lexer"
	"github.com/stretchr/testify/assert"
)

const (
	createPackage = `
  CREATE OR REPLACE PACKAGE BODY MAIN AS
  END MAIN;
  /
	`

	parseFunctionTest1 = `
  (
    IENTITYIDS IN INT,
    ODETAIL_CV OUT VARCHAR,
    ISHARINGACCESS IN INT
  ) IS
	`

	parseFunctionTest2 = `
  IS
	`

	parseFunctionTest3 = `
		PROCEDURE main IS
			li INT := 99;
			lstr varchar := 'this is a string';
		BEGIN
			dbms.print(li);
			dbms.print(lstr);
			IF li > 50 THEN
				dbms.print(50);
			END IF;
			dbms.print(47);
		END;
	`

	parseIfBranches = `
	IF li > 50 THEN
		dbms.print(50);
	END IF;
	END; -- to make the parser end gracefully
	`

	parseIfBranches2 = `
	IF li > 50 THEN
		dbms.print(50);
	END IF;
	dbms.print(99);
	END; -- to make the parser end gracefully
	`
)

func TestParseFunction1(t *testing.T) {
	_, items := lexer.NewLexer("", parseFunctionTest1)
	p := newParser(items)

	pkg := ast.NewPackage("pkg_name")
	f := ast.NewFunction("f_name", true)
	pc := &parserContext{
		pkg:      pkg,
		function: f,
	}
	var pf stateFunc

	pf, pc = parseFunction(p, pc)
	assert.Equal(t, 3, len(f.Proto.Params))
	assert.Equal(t, "parseFunctionBody", getFunctionNameTest(pf))
	assert.NotNil(t, pc.pkg)
	assert.NotNil(t, pc.function)
}

func TestParseFunction2(t *testing.T) {
	_, items := lexer.NewLexer("", parseFunctionTest2)
	p := newParser(items)

	pkg := ast.NewPackage("pkg_name")
	f := ast.NewFunction("f_name", true)
	pc := &parserContext{
		pkg:      pkg,
		function: f,
	}
	var pf stateFunc

	pf, pc = parseFunction(p, pc)
	assert.Equal(t, 0, len(f.Proto.Params))
	assert.Equal(t, "parseFunctionBody", getFunctionNameTest(pf))
	assert.NotNil(t, pc.pkg)
	assert.NotNil(t, pc.function)
}

func TestParseFunction3(t *testing.T) {
	_, items := lexer.NewLexer("", parseFunctionTest3)
	p := newParser(items)

	pkg := ast.NewPackage("pkg_name")
	pc := &parserContext{
		pkg: pkg,
	}
	var pf stateFunc

	pf, pc = parseInsidePackage(p, pc)
	assert.Equal(t, "parseFunction", getFunctionNameTest(pf))
	assert.NotNil(t, pc.pkg)
	assert.NotNil(t, pc.function)
	parsedFunc := pc.function

	pf, pc = parseFunction(p, pc)
	assert.Equal(t, "parseFunctionBody", getFunctionNameTest(pf))
	assert.NotNil(t, pc.pkg)
	assert.NotNil(t, pc.function)

	pf, pc = parseFunctionBody(p, pc)
	assert.Equal(t, "parseInsidePackage", getFunctionNameTest(pf))
	assert.NotNil(t, pc.pkg)
	assert.Nil(t, pc.function)

	log.Printf("%s\n", parsedFunc.String())
}

func TestParseIfBranches(t *testing.T) {
	_, items := lexer.NewLexer("", parseIfBranches)
	p := newParser(items)

	pkg := ast.NewPackage("pkg_name")
	f := ast.NewFunction("f_name", true)
	blk := ast.NewBlock("entry-block")
	f.AddBlock(blk)
	pc := &parserContext{
		pkg:      pkg,
		function: f,
		block:    blk,
	}

	parseInsideBlock(p, pc)
	assert.Equal(t, 3, len(f.Blocks))
	assert.NotNil(t, f.Blocks[0].Terminator)
}

func TestParseIfBranches2(t *testing.T) {
	_, items := lexer.NewLexer("", parseIfBranches2)
	p := newParser(items)

	pkg := ast.NewPackage("pkg_name")
	f := ast.NewFunction("f_name", true)
	blk := ast.NewBlock("entry-block")
	f.AddBlock(blk)
	pc := &parserContext{
		pkg:      pkg,
		function: f,
		block:    blk,
	}

	parseInsideBlock(p, pc)
	assert.Equal(t, 3, len(f.Blocks))
	assert.NotNil(t, f.Blocks[0].Terminator)
	_, ok := f.Blocks[0].Terminator.(*ast.ConditionalBranch)
	assert.True(t, ok)
}

func getFunctionNameTest(i interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	slices := strings.Split(name, ".")
	return slices[len(slices)-1]
}
