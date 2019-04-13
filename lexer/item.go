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

import "fmt"

type ItemType uint8

const (
	EofType ItemType = iota
	ErrorType
	IdentifierType
	KeywordType
	NumericType
	StringType
	OperatorType
	SeparatorType
	TextType
	CommentType
)

type Item struct {
	Typ   ItemType
	Value string
}

func (i *Item) String() string {
	if i.Typ == EofType {
		return "EOF"
	} else if i.Typ == ErrorType {
		return i.Value
	}

	return fmt.Sprintf("%d %q", i.Typ, i.Value)
}
