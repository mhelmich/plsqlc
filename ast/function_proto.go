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
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/types"
)

func NewFunctionProto(name string) *FunctionProto {
	return &FunctionProto{
		Name:   name,
		Params: make([]*FunctionParam, 0),
	}
}

type FunctionProto struct {
	Name   string
	Params []*FunctionParam
}

func (fp *FunctionProto) AddParam(name string, ownership string, t string) {
	fp.Params = append(fp.Params, &FunctionParam{
		Name:      name,
		Ownership: ownership,
		Type:      t,
	})
}

func (fp *FunctionProto) GenIR(cc *CompilerContext) *ir.Func {
	params := make([]*ir.Param, 0)
	for idx := range fp.Params {
		params = append(params, fp.Params[idx].GenIR(cc))
	}

	qualifiedFuncName := cc.currentPackageName + "." + fp.Name
	llvmFunc := cc.llvmModule.NewFunc(qualifiedFuncName, types.Void, params...)
	cc.scopes.addMember(qualifiedFuncName, llvmFunc)
	return llvmFunc
}

func (fp *FunctionProto) String() string {
	var sb strings.Builder
	sb.WriteString(fp.Name)
	sb.WriteString("(")
	for idx := range fp.Params {
		sb.WriteString(fp.Params[idx].String())
		sb.WriteString(" ")
	}
	sb.WriteString(")")
	return sb.String()
}

type FunctionParam struct {
	Name      string
	Ownership string // IN, OUT, INOUT
	Type      string
}

func (fp *FunctionParam) GenIR(cc *CompilerContext) *ir.Param {
	return nil
}

func (fp *FunctionParam) String() string {
	return fmt.Sprintf("%s %s %s", fp.Name, fp.Ownership, fp.Type)
}
