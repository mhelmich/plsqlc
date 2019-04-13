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
	"log"
	"strings"
	"unicode/utf8"
)

type Lexer struct {
	name  string
	input string
	start int
	pos   int
	width int
	items chan *Item
}

func NewLexer(name string, input string) (*Lexer, <-chan *Item) {
	log.Printf("%s\n", input)
	l := &Lexer{
		name:  name,
		input: input,
		items: make(chan *Item),
	}

	go l.run()
	return l, l.items
}

func (l *Lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}

	close(l.items)
}

func (l *Lexer) currentLexItem() string {
	return l.input[l.start:l.pos]
}

func (l *Lexer) lastLexed() rune {
	r, _ := utf8.DecodeRuneInString(l.input[l.pos-l.width-1 : l.pos-1])
	return r
}

func (l *Lexer) emit(t ItemType) {
	var txt string
	if t == KeywordType || t == IdentifierType {
		txt = strings.ToUpper(l.input[l.start:l.pos])
	} else {
		txt = l.input[l.start:l.pos]
	}
	l.items <- &Item{t, txt}
	l.start = l.pos
}

func (l *Lexer) next() (rune rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return int32(EofType)
	}

	rune, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return rune
}

func (l *Lexer) ignore() {
	l.start = l.pos
}

func (l *Lexer) backup() {
	l.pos -= l.width
}

func (l *Lexer) peek() rune {
	rune := l.next()
	l.backup()
	return rune
}

func (l *Lexer) acceptUntilOneOf(until string) {
	r := l.next()
	for strings.IndexRune(until, r) < 0 {
		if r == int32(EofType) {
			l.errorf("Can't find any of '%s' in remaining text", until)
			return
		}
		r = l.next()
	}
}

func (l *Lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *Lexer) acceptMany(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *Lexer) errorf(format string, args ...interface{}) stateFunc {
	l.items <- &Item{
		Typ:   ErrorType,
		Value: fmt.Sprintf(format, args...),
	}
	return nil
}
