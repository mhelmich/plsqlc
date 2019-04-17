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

	"github.com/llir/llvm/ir/value"
)

func NewVariable(name string) *Variable {
	return &Variable{
		Name: name,
	}
}

type Variable struct {
	Name string
}

func (v *Variable) typ() expressionType {
	return variableExpression
}

func (v *Variable) GenIR(cc *CompilerContext) value.Value {
	mem, ok := cc.scopes.findMember(v.Name)
	if !ok {
		log.Panicf("Can't find '%s' in scope", v.Name)
	}
	return cc.currentLlvmBlock.NewLoad(mem)
}

func (v *Variable) String() string {
	return fmt.Sprintf("<variable> %s", v.Name)
}
