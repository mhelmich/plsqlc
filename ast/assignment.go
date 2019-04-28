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

	"github.com/llir/llvm/ir/value"
)

func NewAssignment(varName string, expr Expression) *Assignment {
	return &Assignment{
		VarName: varName,
		Expr:    expr,
	}
}

type Assignment struct {
	VarName string
	Expr    Expression
}

func (a *Assignment) GenIR(cc *CompilerContext) value.Value {
	varValue, ok := cc.scopes.findMember(a.VarName)
	if !ok {
		return nil
	}

	exprValue := a.Expr.GenIR(cc)
	cc.currentLlvmBlock.NewStore(exprValue, varValue)
	return nil
}

func (a *Assignment) String() string {
	return fmt.Sprintf("%s := %s", a.VarName, a.Expr.String())
}
