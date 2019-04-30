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

package lexer

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	basicExample = `
  CREATE OR REPLACE PACKAGE BODY MAIN AS

    PROCEDURE MAIN IS
    BEGIN
      DBMS.PRINT('Hello World!');
      DBMS.PRINT(99);
      N := 100;
      IF N > 100 THEN
        DBMS.PRINT(100);
      ELSE
        DBMS.PRINT(N);
      END IF;
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

func TestIdentifier(t *testing.T) {
	_, items := NewLexer("test-identifier", "input")
	i := <-items
	assert.Equal(t, IdentifierType, i.Typ)
	assert.Equal(t, "INPUT", i.Value)
	i = <-items
	assert.Equal(t, EofType, i.Typ)
	assert.Equal(t, "", i.Value)
}

func TestString(t *testing.T) {
	_, items := NewLexer("test-string", "'input'")
	i := <-items
	assert.Equal(t, StringType, i.Typ)
	assert.Equal(t, "'input'", i.Value)
	i = <-items
	assert.Equal(t, EofType, i.Typ)
	assert.Equal(t, "", i.Value)
}

func TestNumeric(t *testing.T) {
	_, items := NewLexer("test-numberic", "9876")
	i := <-items
	assert.Equal(t, NumericType, i.Typ)
	assert.Equal(t, "9876", i.Value)
}

func TestCommentAndThenString(t *testing.T) {
	_, items := NewLexer("", "-- narf narf narf \n' 9876'")
	i := <-items
	assert.Equal(t, StringType, i.Typ)
	assert.Equal(t, "' 9876'", i.Value)
	i = <-items
	assert.Equal(t, EofType, i.Typ)
	assert.Equal(t, "", i.Value)
}

func TestKeywordVsIdentifier(t *testing.T) {
	_, items := NewLexer("", "BODY MAIN AS ")
	i := <-items
	assert.Equal(t, KeywordType, i.Typ, i.String())
	i = <-items
	assert.Equal(t, IdentifierType, i.Typ, i.String())
	i = <-items
	assert.Equal(t, KeywordType, i.Typ, i.String())
	i = <-items
	assert.Equal(t, EofType, i.Typ, i.String())
}

func TestStringEscapedQuote(t *testing.T) {
	_, items := NewLexer("", "'narf narf\\'narf'")
	i := <-items
	assert.Equal(t, StringType, i.Typ, i.String())
	i = <-items
	assert.Equal(t, EofType, i.Typ, i.String())
}

func TestNarf(t *testing.T) {
	assert.True(t, strings.IndexRune(alphaChars, 'N') >= 0)
}

func TestAssignment(t *testing.T) {
	_, items := NewLexer("", "N := 123;")
	i := <-items
	assert.Equal(t, IdentifierType, i.Typ, i.String())
	i = <-items
	assert.Equal(t, OperatorType, i.Typ, i.String())
	assert.Equal(t, ":=", i.Value, i.String())
	i = <-items
	assert.Equal(t, NumericType, i.Typ, i.String())
	i = <-items
	assert.Equal(t, SeparatorType, i.Typ, i.String())
	i = <-items
	assert.Equal(t, EofType, i.Typ, i.String())
}

func TestBasicExample(t *testing.T) {
	_, items := NewLexer("", basicExample)
	var expectedItemTypes = []ItemType{
		// first line 'CREATE OR REPLACE PACKAGE BODY MAIN AS'
		KeywordType,
		KeywordType,
		KeywordType,
		KeywordType,
		KeywordType,
		IdentifierType,
		KeywordType,
		// second line 'PROCEDURE MAIN IS'
		KeywordType,
		IdentifierType,
		KeywordType,
		// third line 'BEGIN'
		KeywordType,
		// fourth line 'DBMS.PRINT('Hello World!');'
		IdentifierType,
		OperatorType,
		IdentifierType,
		SeparatorType,
		StringType,
		SeparatorType,
		SeparatorType,
		// fifth line 'DBMS.PRINT(99);'
		IdentifierType,
		OperatorType,
		IdentifierType,
		SeparatorType,
		NumericType,
		SeparatorType,
		SeparatorType,
		// sixth line 'N := 100;'
		IdentifierType,
		OperatorType,
		NumericType,
		SeparatorType,
		// seventh line 'IF N > 100 THEN'
		KeywordType,
		IdentifierType,
		OperatorType,
		NumericType,
		KeywordType,
		// eigth line 'DBMS.PRINT(100);'
		IdentifierType,
		OperatorType,
		IdentifierType,
		SeparatorType,
		NumericType,
		SeparatorType,
		SeparatorType,
		// ninth line 'ELSE'
		KeywordType,
		// tenth line 'DBMS.PRINT(N);'
		IdentifierType,
		OperatorType,
		IdentifierType,
		SeparatorType,
		IdentifierType,
		SeparatorType,
		SeparatorType,
		// eleventh line 'END IF;'
		KeywordType,
		KeywordType,
		SeparatorType,
		// twelveth line 'END;'
		KeywordType,
		SeparatorType,
		// thirteenth line 'END MAIN;'
		KeywordType,
		IdentifierType,
		SeparatorType,
		// fourteenth line '/'
		SeparatorType,
		EofType,
	}

	runTest(expectedItemTypes, items, t)
}

func TestControlFlowExample(t *testing.T) {
	_, items := NewLexer("", controlFlowExample)
	var expectedItemTypes = []ItemType{
		// first line 'CREATE OR REPLACE PACKAGE BODY MAIN AS'
		KeywordType,
		KeywordType,
		KeywordType,
		KeywordType,
		KeywordType,
		IdentifierType,
		KeywordType,
		// second line 'PROCEDURE MAIN IS'
		KeywordType,
		IdentifierType,
		KeywordType,
		// third line 'li INT := 99;'
		IdentifierType,
		IdentifierType,
		OperatorType,
		NumericType,
		SeparatorType,
		// fourth line 'BEGIN'
		KeywordType,
		// fifth line 'dbms.print(li);'
		IdentifierType,
		OperatorType,
		IdentifierType,
		SeparatorType,
		IdentifierType,
		SeparatorType,
		SeparatorType,
		// sixth line 'END;'
		KeywordType,
		SeparatorType,
		// seventh line 'END MAIN;'
		KeywordType,
		IdentifierType,
		SeparatorType,
		// eighth line '/'
		SeparatorType,
		EofType,
	}

	runTest(expectedItemTypes, items, t)
}

func runTest(itemTypes []ItemType, items <-chan *Item, t *testing.T) {
	var i *Item
	for idx := range itemTypes {
		i = <-items
		assert.Equal(t, itemTypes[idx], i.Typ, fmt.Sprintf("idx: %d item: %s", idx, i.String()))
	}
}
