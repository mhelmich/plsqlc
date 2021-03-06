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
	"log"
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mhelmich/plsqlc/runtime"
)

func NewFunctionCall(moduleName string, functionName string) *FunctionCall {
	return &FunctionCall{
		ModuleName:   moduleName,
		FunctionName: functionName,
		Args:         make([]Expression, 0),
	}
}

type FunctionCall struct {
	ModuleName   string
	FunctionName string
	Args         []Expression
}

func (fc *FunctionCall) AddArg(expr Expression) {
	fc.Args = append(fc.Args, expr)
}

func (fc *FunctionCall) expressionType() expressionType {
	return functionCallExpression
}

func (fc *FunctionCall) GenIR(cc *CompilerContext) value.Value {
	var fn *ir.Func
	if fc.ModuleName == "DBMS" {
		// aha!
		switch fc.FunctionName {
		case "PRINT":
			switch fc.Args[0].expressionType() {
			case stringExpression:
				fn = cc.getFuncByName(runtime.PrintStringFuncName)
			case numberExpression:
				fn = cc.getFuncByName(runtime.PrintIntFuncName)
			case variableExpression:
				variable := fc.Args[0].(*Variable)
				v, ok := cc.scopes.findMember(variable.Name)
				if !ok {
					log.Panicf("Can't find variable '%s' in scope", variable.Name)
				}

				if v.Type().Equal(types.I64Ptr) {
					fn = cc.getFuncByName(runtime.PrintIntFuncName)
				} else if v.Type().Equal(runtime.StringPointerType) {
					fn = cc.getFuncByName(runtime.PrintStringFuncName)
				} else {
					log.Panicf("Can't find type '%s'", v.Type().String())
				}

			default:
				log.Panicf("Can't find type '%d'", fc.Args[0].expressionType())
			}

		default:
			log.Panicf("Don't recognize runtime function '%s'", fc.FunctionName)
		}
	} else {
		fn = cc.getFuncByName(fc.ModuleName + "." + fc.FunctionName)
	}

	args := make([]value.Value, 0)
	for idx := range fc.Args {
		v := fc.Args[idx].GenIR(cc)

		if fc.Args[idx].expressionType() == stringExpression {
			v = cc.currentLlvmBlock.NewLoad(v)
		}

		args = append(args, v)
	}

	funcCall := cc.currentLlvmBlock.NewCall(fn, args...)
	return funcCall
}

func (fc *FunctionCall) String() string {
	var sb strings.Builder
	sb.WriteString("<func call> ")
	sb.WriteString(fc.ModuleName)
	sb.WriteString(".")
	sb.WriteString(fc.FunctionName)
	sb.WriteString("(")
	for idx := range fc.Args {
		sb.WriteString(fc.Args[idx].String())
		sb.WriteString(",")
	}
	sb.WriteString(")")
	return sb.String()
}
