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

package compiler

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/llir/llvm/ir"
	"github.com/mhelmich/plsqlc/ast"
	"github.com/mhelmich/plsqlc/lexer"
	"github.com/mhelmich/plsqlc/parser"
	"github.com/mhelmich/plsqlc/runtime"
)

func Compile(inputPath string, outputPath string, printIR bool, deleteLlvmIR bool) {
	mod := ir.NewModule()
	runtime.GenerateInModule(mod)
	mod = compileCode(inputPath, mod)
	runtime.GenerateMain(mod)

	// tmpFile, err := ioutil.TempFile("", "_plsqlc")
	tmpFile, err := os.Create("_temp_llvm_.ll")
	if err != nil {
		log.Panicf("%s", err.Error())
	}

	fileName := tmpFile.Name()
	tmpFile.Close()
	if deleteLlvmIR {
		defer os.Remove(fileName)
	}

	ir := mod.String()

	if printIR {
		log.Printf("%s", ir)
	}

	err = ioutil.WriteFile(fileName, []byte(ir), 0644)
	if err != nil {
		log.Panicf("%s", err.Error())
	}

	clangArgs := []string{
		fileName,               // input path of temp file
		"-Wno-override-module", // disable override target triple warnings
		"-o", outputPath,       // output path
		"-O3",
	}

	if printIR {
		log.Printf("clang %v\n", clangArgs)
	}

	cmd := exec.Command("clang", clangArgs...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		log.Panicf("%s", string(output))
	}

	if len(output) > 0 {
		log.Panicf("%s", string(output))
	}
}

func compileCode(in string, mod *ir.Module) *ir.Module {
	file, err := os.Open(in)
	if err != nil {
		log.Panicf("failed reading file: %s", err)
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Panicf("%s", err.Error())
	}

	_, items := lexer.NewLexer("", string(data))
	p := parser.NewParser(items)
	// map[string]*ast.Package
	namesToPackages := p.GetPackageAsts()
	for name, pkg := range namesToPackages {
		if name == "MAIN" {
			if !pkg.HasMainFunction() {
				log.Panicf("Can't find 'main' function")
			}
			cc := ast.NewCompilerContext(mod)
			pkg.GenIR(cc)
			return cc.GetIRModule()
		}
	}

	log.Panicf("Can't find 'main' package")
	return nil
}
