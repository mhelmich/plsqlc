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

	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mhelmich/plsqlc/runtime"
)

func NewStringLiteral(value string) *StringLiteral {
	value = value[1 : len(value)-1]
	return &StringLiteral{
		Value: value,
	}
}

type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) expressionType() expressionType {
	return stringExpression
}

func (sl *StringLiteral) GenIR(cc *CompilerContext) value.Value {
	stringType := cc.getTypeByName(runtime.StringTypeName)

	strStruct := cc.currentLlvmBlock.NewAlloca(stringType)
	dataPtr := cc.currentLlvmBlock.NewGetElementPtr(strStruct, llvmZero, llvmZero)
	lenPtr := cc.currentLlvmBlock.NewGetElementPtr(strStruct, llvmZero, llvmOne)
	cc.currentLlvmBlock.NewStore(constant.NewInt(types.I64, int64(len(sl.Value))), lenPtr)

	a := cc.currentLlvmBlock.NewAlloca(types.NewArray(uint64(len(sl.Value)), types.I8))
	cc.currentLlvmBlock.NewStore(constant.NewCharArrayFromString(sl.Value), a)
	strPtr := cc.currentLlvmBlock.NewBitCast(a, types.NewPointer(types.I8))
	cc.currentLlvmBlock.NewStore(strPtr, dataPtr)

	return strStruct
}

func (sl *StringLiteral) String() string {
	return fmt.Sprintf("<string literal> %s", sl.Value)
}
