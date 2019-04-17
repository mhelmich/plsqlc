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
	"testing"

	"github.com/mhelmich/plsqlc/lexer"
	"github.com/stretchr/testify/assert"
)

const (
	basicExample = `
  CREATE OR REPLACE PACKAGE BODY MAIN AS

    PROCEDURE MAIN IS
    BEGIN
      DBMS.PRINT('Hello World!');
      DBMS.PRINT(99);
    END;

	END MAIN;
  /
`

	controlFlowExample = `
  CREATE OR REPLACE PACKAGE BODY main AS

		PROCEDURE main IS
			li INT := 99;
		BEGIN
			dbms.print(li);
		END;

  END main;
  /
`
)

func TestParserPeek(t *testing.T) {
	_, items := lexer.NewLexer("", basicExample)
	p := newParser(items)

	peekedItem := p.peek()
	consumedItem := p.next()
	assert.Equal(t, peekedItem.Value, consumedItem.Value)
	assert.Equal(t, "CREATE", consumedItem.Value)
	assert.True(t, p.acceptValue("OR"))
	assert.Equal(t, "REPLACE", p.peek().Value)
	assert.Equal(t, "REPLACE", p.peek().Value)
	assert.True(t, p.acceptValue("REPLACE"))
	b, v := p.acceptType(lexer.KeywordType)
	assert.True(t, b)
	assert.Equal(t, "PACKAGE", v)
}

func TestBasic(t *testing.T) {
	_, items := lexer.NewLexer("test-basic", basicExample)
	p := newParser(items)
	p.run()
	for k, v := range p.packages {
		log.Printf("%s %s", k, v.String())
	}
	assert.Equal(t, 1, len(p.packages))
}

func TestControlFlow(t *testing.T) {
	_, items := lexer.NewLexer("", controlFlowExample)
	p := newParser(items)
	p.run()
	for k, v := range p.packages {
		log.Printf("%s %s", k, v.String())
	}
	assert.Equal(t, 1, len(p.packages))
}
