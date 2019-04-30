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

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

type expressionType int

const (
	stringExpression expressionType = iota
	numberExpression
	functionCallExpression
	variableExpression
	binOpExpression
)

var (
	llvmZero = constant.NewInt(types.I32, 0)
	llvmOne  = constant.NewInt(types.I32, 1)
)

type Node interface {
	String() string
}

type CompilerContext struct {
	llvmModule         *ir.Module
	currentPackageName string
	currentLlvmFunc    *ir.Func
	functionBlocks     map[*Block]*ir.Block
	currentLlvmBlock   *ir.Block
	scopes             *scope
}

func NewCompilerContext(mod *ir.Module) *CompilerContext {
	return &CompilerContext{
		llvmModule: mod,
		scopes:     newScope(),
	}
}

func (cc *CompilerContext) getFuncByName(n string) *ir.Func {
	for idx := range cc.llvmModule.Funcs {
		if cc.llvmModule.Funcs[idx].Name() == n {
			return cc.llvmModule.Funcs[idx]
		}
	}

	log.Panicf("Can't find function '%s'", n)
	return nil
}

func (cc *CompilerContext) getGlobalByName(n string) *ir.Global {
	for idx := range cc.llvmModule.Globals {
		if cc.llvmModule.Globals[idx].Name() == n {
			return cc.llvmModule.Globals[idx]
		}
	}

	log.Panicf("Can't find global %s", n)
	return nil
}

func (cc *CompilerContext) getTypeByName(n string) types.Type {
	for idx := range cc.llvmModule.TypeDefs {
		if cc.llvmModule.TypeDefs[idx].Name() == n {
			return cc.llvmModule.TypeDefs[idx]
		}
	}

	log.Panicf("Can't find type %s", n)
	return nil
}

func (cc *CompilerContext) GetIRModule() *ir.Module {
	return cc.llvmModule
}

func (cc *CompilerContext) pushScope() *scope {
	ns := newScope()
	ns.Parent = cc.scopes
	cc.scopes = ns
	return ns
}

func (cc *CompilerContext) popScope() *scope {
	s := cc.scopes
	if s.Parent == nil {
		log.Panicf("Can't pop root scope!")
	}

	x := cc.scopes.Parent
	s.valid = false
	s.Parent = nil
	cc.scopes = x
	return x
}

func newScope() *scope {
	return &scope{
		Members: make(map[string]value.Value),
		valid:   true,
	}
}

type scope struct {
	Parent  *scope
	Members map[string]value.Value
	valid   bool
}

func (s *scope) addMember(name string, val value.Value) {
	s.Members[name] = val
}

func (s *scope) findMember(name string) (value.Value, bool) {
	if !s.valid {
		log.Panicf("Scope is not valid!")
	}

	hasParent := true
	x := s
	for hasParent {
		v, ok := x.Members[name]
		if ok {
			return v, true
		}
		hasParent = x.Parent != nil
		x = x.Parent
	}

	return nil, false
}

type Instruction interface {
	GenIR(cc *CompilerContext) value.Value
	String() string
}

type Expression interface {
	GenIR(cc *CompilerContext) value.Value
	typ() expressionType
	String() string
}
