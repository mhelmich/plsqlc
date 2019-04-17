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
	"strings"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
	"github.com/mhelmich/plsqlc/runtime"
)

func NewFunction(fn string) *Function {
	return &Function{
		Proto:  NewFunctionProto(fn),
		Locals: make([]*FunctionLocal, 0),
		Blocks: make([]*Block, 0),
	}
}

type Function struct {
	Proto  *FunctionProto
	Locals []*FunctionLocal
	Blocks []*Block
}

func (f *Function) AddParam(name string, ownership string, t string) {
	f.Proto.AddParam(name, ownership, t)
}

func (f *Function) AddLocal(name string, typ string, value string) {
	fl := &FunctionLocal{
		Name:  name,
		Typ:   typ,
		Value: value,
	}
	f.Locals = append(f.Locals, fl)
}

func (f *Function) AddBlock(b *Block) {
	f.Blocks = append(f.Blocks, b)
}

func (f *Function) GenIRForProtos(cc *CompilerContext) value.Value {
	f.Proto.GenIR(cc)
	return nil
}

func (f *Function) GenIR(cc *CompilerContext) value.Value {
	llvmFunc := cc.getFuncByName(cc.currentPackageName + "." + f.Proto.Name)
	cc.currentLlvmFunc = llvmFunc
	hasLocals := len(f.Locals) > 0
	var blocks []*ir.Block
	var additionalBlocks int
	if hasLocals {
		// if we have locals in this function, we add an extra block just for
		// the declaration and definition of all locals
		blocks = make([]*ir.Block, len(f.Blocks)+1)
		additionalBlocks = 1

		blocks[0] = cc.currentLlvmFunc.NewBlock("locals")
		cc.currentLlvmBlock = blocks[0]

		// locals have their own block
		for idx := range f.Locals {
			f.Locals[idx].GenIR(cc)
		}
	} else {
		blocks = make([]*ir.Block, len(f.Blocks))
		additionalBlocks = 0
	}

	// create code blocks
	for idx := range f.Blocks {
		blocks[idx+additionalBlocks] = cc.currentLlvmFunc.NewBlock("")
		cc.currentLlvmBlock = blocks[idx+additionalBlocks]
		f.Blocks[idx].GenIR(cc)
	}

	// wire up the first two blocks
	if hasLocals {
		blocks[0].NewBr(blocks[1])
	}

	cc.currentLlvmBlock = nil
	cc.currentLlvmFunc = nil
	return llvmFunc
}

func (f *Function) String() string {
	var sb strings.Builder
	sb.WriteString("<func definition> ")
	sb.WriteString(f.Proto.String())
	sb.WriteString("\n")

	for idx := range f.Locals {
		sb.WriteString(f.Locals[idx].String())
		sb.WriteString("\n")
	}

	sb.WriteString("\n")

	for idx := range f.Blocks {
		sb.WriteString(f.Blocks[idx].String())
		sb.WriteString("\n")
	}

	return sb.String()
}

type FunctionLocal struct {
	Name  string
	Typ   string
	Value string
}

func (fl *FunctionLocal) GenIR(cc *CompilerContext) value.Value {
	alloca := cc.currentLlvmBlock.NewAlloca(plsqlTypeToLLVMType(fl.Typ))
	cc.scopes.addMember(fl.Name, alloca)
	if fl.Value != "" {
		switch fl.Typ {

		case "INT":
			i, err := strconv.ParseInt(fl.Value, 10, 64)
			if err != nil {
				log.Panicf("%s", err.Error())
			}
			cc.currentLlvmBlock.NewStore(constant.NewInt(types.I64, i), alloca)

		case "VARCHAR":
			// chop off the 's on both ends
			s := fl.Value[1 : len(fl.Value)-1]
			runtime.MakeString(s, cc.currentLlvmBlock, alloca)

		default:
			log.Panicf("Local for type '%s' not implemented", fl.Typ)
		}

	}
	return alloca
}

func (fl *FunctionLocal) String() string {
	return fmt.Sprintf("%s %s %s", fl.Name, fl.Typ, fl.Value)
}
