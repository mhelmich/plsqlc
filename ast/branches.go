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

	"github.com/llir/llvm/ir/value"
)

func NewConditionalBranch(condition *BinOp, trueTarget *Block, falseTarget *Block) *ConditionalBranch {
	return &ConditionalBranch{
		Condition:   condition,
		TrueTarget:  trueTarget,
		FalseTarget: falseTarget,
	}
}

type ConditionalBranch struct {
	Condition   *BinOp
	TrueTarget  *Block
	FalseTarget *Block
}

func (b *ConditionalBranch) GenIR(cc *CompilerContext) value.Value {
	cond := b.Condition.GenIR(cc)
	ttBlk := cc.functionBlocks[b.TrueTarget]
	ftBlk := cc.functionBlocks[b.FalseTarget]
	cc.currentLlvmBlock.NewCondBr(cond, ttBlk, ftBlk)
	return nil
}

func (b *ConditionalBranch) String() string {
	var falseTargetName string
	if b.FalseTarget == nil {
		falseTargetName = "<empty>"
	} else {
		falseTargetName = b.FalseTarget.Name
	}

	return fmt.Sprintf("%s\n%s\n%s", b.Condition, b.TrueTarget.Name, falseTargetName)
}

func NewBranch(blk *Block) *Branch {
	return &Branch{
		Blk: blk,
	}
}

type Branch struct {
	Blk *Block
}

func (b *Branch) GenIR(cc *CompilerContext) value.Value {
	llvmBlk := cc.functionBlocks[b.Blk]
	cc.currentLlvmBlock.NewBr(llvmBlk)
	return nil
}

func (b *Branch) String() string {
	return b.Blk.Name
}
