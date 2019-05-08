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

package ast

import (
	"fmt"
	"log"
	"strconv"

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

func NewNumericLiteral(value string) *NumericLiteral {
	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Panicf("Can't convert '%s' into number!", value)
	}
	return &NumericLiteral{
		Value: i,
	}
}

type NumericLiteral struct {
	Value int64
}

func (nl *NumericLiteral) expressionType() expressionType {
	return numberExpression
}

func (nl *NumericLiteral) GenIR(cc *CompilerContext) value.Value {
	return constant.NewInt(types.I64, nl.Value)
}

func (nl *NumericLiteral) String() string {
	return fmt.Sprintf("<numeric literal> %d", nl.Value)
}
