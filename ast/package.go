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

func NewPackage(name string) *Package {
	return &Package{
		Name: name,
	}
}

type Package struct {
	Name      string
	functions []*Function
}

func (p *Package) GenIR(cc *CompilerContext) error {
	cc.currentPackageName = p.Name
	// TDOD: first declare all types
	// secondly declare all functions
	for idx := range p.functions {
		p.functions[idx].GenIRForProtos(cc)
	}
	// thirdly compile all the code
	for idx := range p.functions {
		p.functions[idx].GenIR(cc)
	}
	cc.currentPackageName = ""
	return nil
}

func (p *Package) AddFunction(f *Function) {
	p.functions = append(p.functions, f)
}

func (p *Package) HasMainFunction() bool {
	for idx := range p.functions {
		if p.functions[idx].Proto.Name == "MAIN" {
			return true
		}
	}
	return false
}

func (p *Package) String() string {
	var sb strings.Builder
	sb.WriteString(p.Name)
	sb.WriteString("\n")
	for idx := range p.functions {
		sb.WriteString("  ")
		sb.WriteString(p.Name)
		sb.WriteString(".")
		sb.WriteString(p.functions[idx].String())
		sb.WriteString("\n")
	}
	return sb.String()
}
