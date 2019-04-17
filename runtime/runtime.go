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

package runtime

import (
	"log"

	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
	"github.com/llir/llvm/ir/value"
)

const (
	digits = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var (
	llvmZeroI32 = constant.NewInt(types.I32, 0)
	llvmOneI32  = constant.NewInt(types.I32, 1)

	llvmZeroI64 = constant.NewInt(types.I64, 0)
	llvmOneI64  = constant.NewInt(types.I64, 1)

	StringType        types.Type
	StringPointerType types.Type
)

func GenerateInModule(mod *ir.Module) {
	mod.NewFunc("putchar", types.I32, ir.NewParam("c", types.I8))
	mod.NewGlobalDef("_runtime.digits", constant.NewCharArrayFromString(digits))

	stringStruct := types.NewStruct(types.NewPointer(types.I8), types.I64)
	stringStruct.SetName("_runtime._string")
	StringType = mod.NewTypeDef("_runtime._string", stringStruct)
	StringPointerType = types.NewPointer(StringType)

	generate_printInt(mod)
	generateprintInt(mod)
	generateprintStr(mod)
}

func GenerateMain(mod *ir.Module) {
	userMain := getFuncByName("MAIN.MAIN", mod)
	main := mod.NewFunc("main", types.I32)
	b := main.NewBlock("plsql-main")
	b.NewCall(userMain)
	b.NewRet(constant.NewInt(types.I32, 0))
}

func generateTestMain(mod *ir.Module) {
	printInt := getFuncByName("_runtime.printInt", mod)
	printStr := getFuncByName("_runtime.printStr", mod)

	main := mod.NewFunc("main", types.I32)
	bmain := main.NewBlock("main-main")
	bmain.NewCall(printInt, constant.NewInt(types.I32, 5432))

	strStruct := makeStringWithAlloca("\nHello World!\n", bmain)
	bmain.NewCall(printStr, bmain.NewLoad(strStruct))

	bmain.NewRet(constant.NewInt(types.I32, 0))
}

func makeStringWithAlloca(s string, b *ir.Block) value.Named {
	strStruct := b.NewAlloca(StringType)
	return MakeString(s, b, strStruct)
}

func MakeString(s string, b *ir.Block, strStruct *ir.InstAlloca) value.Named {
	dataPtr := b.NewGetElementPtr(strStruct, llvmZeroI32, llvmZeroI32)
	lenPtr := b.NewGetElementPtr(strStruct, llvmZeroI32, llvmOneI32)
	b.NewStore(constant.NewInt(types.I64, int64(len(s))), lenPtr)

	a := b.NewAlloca(types.NewArray(uint64(len(s)), types.I8))
	b.NewStore(constant.NewCharArrayFromString(s), a)
	strPtr := b.NewBitCast(a, types.NewPointer(types.I8))
	b.NewStore(strPtr, dataPtr)

	return strStruct
}

func generate_printInt(mod *ir.Module) {
	gDigits := getGlobalByName("_runtime.digits", mod)
	putchar := getFuncByName("putchar", mod)

	input := ir.NewParam("input", types.I64)
	base := ir.NewParam("base", types.I64)
	internalPrintInt := mod.NewFunc("_runtime._printInt", types.Void, input, base)
	entry := internalPrintInt.NewBlock("entry")

	cmp := entry.NewICmp(enum.IPredUGT, input, base)
	thenBlock := internalPrintInt.NewBlock("then")
	elseBlock := internalPrintInt.NewBlock("merge")
	entry.NewCondBr(cmp, thenBlock, elseBlock)

	tmpDiv := thenBlock.NewUDiv(input, base)
	thenBlock.NewCall(internalPrintInt, tmpDiv, base)
	thenBlock.NewBr(elseBlock)

	rem := elseBlock.NewURem(input, base)
	charPtr := elseBlock.NewGetElementPtr(gDigits, llvmZeroI32, rem)
	elseBlock.NewCall(putchar, elseBlock.NewLoad(charPtr))
	elseBlock.NewRet(nil)
}

func generateprintInt(mod *ir.Module) {
	internalPrintInt := getFuncByName("_runtime._printInt", mod)
	putchar := getFuncByName("putchar", mod)

	input := ir.NewParam("input", types.I64)
	printInt := mod.NewFunc("_runtime.printInt", types.Void, input)
	entry := printInt.NewBlock("entry")

	alloca := entry.NewAlloca(types.I64)
	cmp := entry.NewICmp(enum.IPredSGT, constant.NewInt(types.I64, 0), input)
	thenBlock := printInt.NewBlock("then")
	elseBlock := printInt.NewBlock("else")
	mergeBlock := printInt.NewBlock("merge")
	entry.NewCondBr(cmp, thenBlock, elseBlock)

	thenBlock.NewCall(putchar, constant.NewInt(types.I8, '-'))
	mul := thenBlock.NewMul(input, constant.NewInt(types.I64, -1))
	thenBlock.NewStore(mul, alloca)
	thenBlock.NewBr(mergeBlock)

	elseBlock.NewStore(input, alloca)
	elseBlock.NewBr(mergeBlock)

	i := mergeBlock.NewLoad(alloca)
	mergeBlock.NewCall(internalPrintInt, i, constant.NewInt(types.I64, 10))
	mergeBlock.NewCall(putchar, constant.NewInt(types.I8, '\n'))
	mergeBlock.NewRet(nil)
}

func generateprintStr(mod *ir.Module) {
	putchar := getFuncByName("putchar", mod)
	stringType := getTypeByName("_runtime._string", mod)

	strInput := ir.NewParam("s", stringType)
	printStr := mod.NewFunc("_runtime.printStr", types.Void, strInput)
	entryBB := printStr.NewBlock("entry")
	whileBB := printStr.NewBlock("loop-body")
	mergeBB := printStr.NewBlock("loop-merge")

	i := entryBB.NewAlloca(types.I64)
	entryBB.NewStore(llvmZeroI64, i)
	str := entryBB.NewExtractValue(strInput, 0)
	len := entryBB.NewExtractValue(strInput, 1)
	cmp := entryBB.NewICmp(enum.IPredSLT, entryBB.NewLoad(i), len)
	entryBB.NewCondBr(cmp, whileBB, mergeBB)

	iLoaded := whileBB.NewLoad(i)
	charPtr := whileBB.NewGetElementPtr(str, iLoaded)
	whileBB.NewCall(putchar, whileBB.NewLoad(charPtr))
	// i++
	whileBB.NewStore(whileBB.NewAdd(llvmOneI64, iLoaded), i)
	cmp = whileBB.NewICmp(enum.IPredSLT, whileBB.NewLoad(i), len)
	whileBB.NewCondBr(cmp, whileBB, mergeBB)

	mergeBB.NewCall(putchar, constant.NewInt(types.I8, '\n'))
	mergeBB.NewRet(nil)
}

func getGlobalByName(n string, mod *ir.Module) *ir.Global {
	for idx := range mod.Globals {
		if mod.Globals[idx].Name() == n {
			return mod.Globals[idx]
		}
	}

	log.Panicf("Can't find global %s", n)
	return nil
}

func getFuncByName(n string, mod *ir.Module) *ir.Func {
	for idx := range mod.Funcs {
		if mod.Funcs[idx].Name() == n {
			return mod.Funcs[idx]
		}
	}

	log.Panicf("Can't find function %s", n)
	return nil
}

func getTypeByName(n string, mod *ir.Module) types.Type {
	for idx := range mod.TypeDefs {
		if mod.TypeDefs[idx].Name() == n {
			return mod.TypeDefs[idx]
		}
	}

	log.Panicf("Can't find type %s", n)
	return nil
}
