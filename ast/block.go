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
	"strconv"
	"strings"

	"github.com/llir/llvm/ir/value"
)

var blockNameCounter int64

func NewBlock(name string) *Block {
	n := name + "-" + strconv.FormatInt(blockNameCounter, 10)
	blockNameCounter++
	return &Block{
		Name: n,
	}
}

type Block struct {
	Name         string
	Instructions []Instruction
	Terminator   Instruction
}

func (b *Block) AddInstruction(i Instruction) {
	b.Instructions = append(b.Instructions, i)
}

func (b *Block) GenIR(cc *CompilerContext) value.Value {
	cc.currentLlvmBlock = cc.functionBlocks[b]
	for idx := range b.Instructions {
		b.Instructions[idx].GenIR(cc)
	}

	if cc.currentLlvmBlock.Term == nil && b.Terminator == nil {
		log.Printf("Didn't find a terminator in block '%s'! Filled in empty return.\n", b.Name)
		cc.currentLlvmBlock.NewRet(nil)
	} else {
		b.Terminator.GenIR(cc)
	}

	if cc.currentLlvmBlock.Term == nil {
		log.Panicf("%s has no terminator!\n", b.Name)
	}

	return cc.currentLlvmBlock
}

func (b *Block) String() string {
	var sb strings.Builder
	sb.WriteString(b.Name)
	sb.WriteString("\n")

	if len(b.Instructions) > 0 {
		for idx := range b.Instructions {
			sb.WriteString(b.Instructions[idx].String())
			sb.WriteString("\n")
		}
	} else {
		sb.WriteString("<empty>")
	}

	sb.WriteString("Terminator:\n")
	if b.Terminator == nil {
		sb.WriteString("<empty>")
	} else {
		sb.WriteString(b.Terminator.String())
	}
	sb.WriteString("\n")

	return sb.String()
}
