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

func NewFunction(fn string, isProcedure bool) *Function {
	return &Function{
		Proto:       NewFunctionProto(fn),
		Locals:      make([]*FunctionLocal, 0),
		Blocks:      make([]*Block, 0),
		isProcedure: isProcedure,
	}
}

type Function struct {
	Proto       *FunctionProto
	Locals      []*FunctionLocal
	Blocks      []*Block
	isProcedure bool
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
	return f.Proto.GenIR(cc)
}

func (f *Function) GenIR(cc *CompilerContext) value.Value {
	llvmFunc := cc.getFuncByName(cc.currentPackageName + "." + f.Proto.Name)
	cc.currentLlvmFunc = llvmFunc
	if len(f.Locals) > 0 {
		localsBlock := cc.currentLlvmFunc.NewBlock("locals")
		cc.currentLlvmBlock = localsBlock
		// locals have their own block
		for idx := range f.Locals {
			f.Locals[idx].GenIR(cc)
		}

		// create all llvm blocks ahead of time
		cc.functionBlocks = make(map[*Block]*ir.Block)
		for idx := range f.Blocks {
			cc.functionBlocks[f.Blocks[idx]] = cc.currentLlvmFunc.NewBlock(f.Blocks[idx].Name)
		}

		// generate llvm ir for all blocks
		for idx := range f.Blocks {
			f.Blocks[idx].GenIR(cc)
		}

		// link locals block to method entry block
		// that is the first block in the list
		localsBlock.NewBr(cc.functionBlocks[f.Blocks[0]])
	} else {
		cc.functionBlocks = make(map[*Block]*ir.Block)
		for idx := range f.Blocks {
			cc.functionBlocks[f.Blocks[idx]] = cc.currentLlvmFunc.NewBlock(f.Blocks[idx].Name)
		}

		for idx := range f.Blocks {
			f.Blocks[idx].GenIR(cc)
		}
	}

	cc.currentLlvmBlock = nil
	cc.currentLlvmFunc = nil
	cc.functionBlocks = nil
	return llvmFunc
}

func (f *Function) String() string {
	var sb strings.Builder
	sb.WriteString("<func definition> ")
	sb.WriteString(f.Proto.String())
	sb.WriteString("\n")
	sb.WriteString("Locals:\n")

	for idx := range f.Locals {
		sb.WriteString(f.Locals[idx].String())
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString("Blocks:\n")

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
