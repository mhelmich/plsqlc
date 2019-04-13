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
	"strings"
)

func NewFunction(fn string) *Function {
	return &Function{
		Proto:  NewFunctionProto(fn),
		Blocks: make([]*Block, 0),
	}
}

type Function struct {
	Proto  *FunctionProto
	Blocks []*Block
}

func (f *Function) AddParam(name string, ownership string, t string) {
	f.Proto.AddParam(name, ownership, t)
}

func (f *Function) AddBlock(b *Block) {
	f.Blocks = append(f.Blocks, b)
}

func (f *Function) GenIRForProtos(cc *CompilerContext) error {
	f.Proto.GenIR(cc)
	return nil
}

func (f *Function) GenIR(cc *CompilerContext) error {
	cc.currentLlvmFunc = cc.getFuncByName(cc.currentPackageName + "." + f.Proto.Name)
	for idx := range f.Blocks {
		f.Blocks[idx].GenIR(cc)
	}
	cc.currentLlvmFunc = nil
	return nil
}

func (f *Function) String() string {
	var sb strings.Builder
	sb.WriteString("<func definition> ")
	sb.WriteString(f.Proto.String())
	sb.WriteString("\n")
	for idx := range f.Blocks {
		sb.WriteString(f.Blocks[idx].String())
		sb.WriteString("\n")
	}
	return sb.String()
}
