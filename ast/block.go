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
)

func NewBlock(name string) *Block {
	return &Block{
		Name: name,
	}
}

type Block struct {
	Name         string
	Instructions []Instruction
}

func (b *Block) AddInstruction(i Instruction) {
	b.Instructions = append(b.Instructions, i)
}

func (b *Block) GenIR(cc *CompilerContext) error {
	cc.currentLlvmBlock.SetName(b.Name)
	for idx := range b.Instructions {
		b.Instructions[idx].GenIR(cc)
	}

	if cc.currentLlvmBlock.Term == nil {
		log.Printf("Didn't find a terminator! Filled in empty return.\n")
		cc.currentLlvmBlock.NewRet(nil)
	}

	return nil
}

func (b *Block) String() string {
	var sb strings.Builder
	for idx := range b.Instructions {
		sb.WriteString(b.Instructions[idx].String())
		sb.WriteString("\n")
	}
	return sb.String()
}
