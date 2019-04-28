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
	"log"
	"strings"
)

const (
	commentToken = "--"

	alphaChars   = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM"
	numericChars = "1234567890"
	specialChars = "_"

	separatorChars = ";(),/"
	operatorChars  = "<>:.=-" // contains ':' so that ':=' can be found
)

// if types are keywords, the parser gets more complicated
// as it might need to parse a keywords (primitive types) or
// identifiers (custom types)
var keywords = map[string]bool{
	"CREATE":    true,
	"OR":        true,
	"REPLACE":   true,
	"PACKAGE":   true,
	"BODY":      true,
	"AS":        true,
	"PROCEDURE": true,
	"IS":        true,
	"BEGIN":     true,
	"END":       true,
	"IF":        true,
	"THEN":      true,
	"ELSE":      true,
	"LOOP":      true,
	"WHILE":     true,
}

type stateFunc func(*Lexer) stateFunc

func lexText(l *Lexer) stateFunc {
	for {
		if strings.HasPrefix(l.input[l.pos:], commentToken) {
			l.pos += len(commentToken)
			l.acceptUntilOneOf("\n")
			l.ignore()
			continue
		}

		for {
			switch r := l.next(); {
			case r == int32(EofType):
				l.emit(EofType)
				return nil
			case isSpace(r):
				l.ignore()
			case contains(separatorChars, r):
				return lexSeparator
			case contains(operatorChars, r):
				return lexOperator
			case r == '\'':
				return lexString
			case contains(alphaChars, r):
				l.backup()
				return lexIdentifier
			case contains(numericChars+"+-", r):
				return lexNumeric
			default:
				log.Panicf("Found %s but can't match a rule", string(r))
			}
		}
	}
}

func lexOperator(l *Lexer) stateFunc {
	l.accept("=") // try accepting a '=' to complete a probable ':='
	l.emit(OperatorType)
	return lexText
}

func lexSeparator(l *Lexer) stateFunc {
	l.emit(SeparatorType)
	return lexText
}

func lexIdentifier(l *Lexer) stateFunc {
	l.acceptMany(alphaChars + numericChars + specialChars)
	upperTxt := strings.ToUpper(l.currentLexItem())
	if _, ok := keywords[upperTxt]; ok {
		l.emit(KeywordType)
	} else {
		l.emit(IdentifierType)
	}
	return lexText
}

func lexString(l *Lexer) stateFunc {
	for {
		l.acceptUntilOneOf("'")
		ll := l.lastLexed()
		if ll != '\\' {
			break
		}
	}

	l.emit(StringType)
	return lexText
}

func lexNumeric(l *Lexer) stateFunc {
	l.acceptMany(numericChars)
	l.emit(NumericType)
	return lexText
}

func isSpace(r rune) bool {
	if r == ' ' || r == '\t' || r == '\n' || r == '\r' {
		return true
	}
	return false
}

func contains(valid string, r rune) bool {
	return strings.IndexRune(valid, r) >= 0
}
