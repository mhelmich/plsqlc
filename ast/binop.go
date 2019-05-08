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
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mhelmich/plsqlc/runtime"
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

func (bo *BinOp) expressionType() expressionType {
	return binOpExpression
}

func (bo *BinOp) GenIR(cc *CompilerContext) value.Value {
	l := bo.Left.GenIR(cc)
	r := bo.Right.GenIR(cc)

	switch bo.Op {
	case ">":
		if bo.Left.expressionType() == numberExpression || bo.Right.expressionType() == numberExpression {
			return cc.currentLlvmBlock.NewICmp(enum.IPredSGT, l, r)
		} else {
			log.Panicf("Not implemented yet - '>' '%v' '%v'", bo.Left.expressionType(), bo.Right.expressionType())
		}
	case "-":
		if bo.Left.expressionType() == numberExpression || bo.Right.expressionType() == numberExpression {
			return cc.currentLlvmBlock.NewSub(l, r)
		} else {
			log.Panicf("Not implemented yet - '-' '%v' '%v'", bo.Left.expressionType(), bo.Right.expressionType())
		}
	case "=":
		if bo.Left.expressionType() == numberExpression || bo.Right.expressionType() == numberExpression {
			return cc.currentLlvmBlock.NewICmp(enum.IPredEQ, l, r)
		} else if bo.Left.expressionType() == stringExpression || bo.Right.expressionType() == stringExpression {
			if types.IsPointer(l.Type()) {
				l = cc.currentLlvmBlock.NewLoad(l)
			}

			if types.IsPointer(r.Type()) {
				r = cc.currentLlvmBlock.NewLoad(r)
			}

			return cc.currentLlvmBlock.NewCall(runtime.EqualStringFunc, l, r)
		} else {
			log.Panicf("Not implemented yet - '-' '%v' '%v'", bo.Left.expressionType(), bo.Right.expressionType())
		}
	default:
		log.Panicf("Operation '%s' hasn't been implemented yet", bo.Op)
	}
	return nil
}

func (bo *BinOp) String() string {
	return fmt.Sprintf("<left> %s <op> '%s' <right> %s", bo.Left.String(), bo.Op, bo.Right.String())
}
