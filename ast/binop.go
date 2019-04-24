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

	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/value"
)

func NewBinOp(left Expression, op string, right Expression) *BinOp {
	return &BinOp{
		Left:  left,
		Op:    op,
		Right: right,
	}
}

type BinOp struct {
	Left  Expression
	Op    string
	Right Expression
}

func (bo *BinOp) typ() expressionType {
	return binOpExpression
}

func (bo *BinOp) GenIR(cc *CompilerContext) value.Value {
	l := bo.Left.GenIR(cc)
	r := bo.Right.GenIR(cc)

	switch bo.Op {
	case ">":
		return cc.currentLlvmBlock.NewICmp(enum.IPredSGT, l, r)
	default:
		log.Panicf("Operation '%s' hasn't been implemented yet", bo.Op)
	}
	return nil
}

func (bo *BinOp) String() string {
	return fmt.Sprintf("<left> %s <op> '%s' <right> %s", bo.Left.String(), bo.Op, bo.Right.String())
}
