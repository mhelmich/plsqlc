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
	"fmt"
	"os"
	"os/exec"
	"testing"

	"io/ioutil"

	"github.com/llir/llvm/ir"
	"github.com/stretchr/testify/assert"
)

var basicOutput = "5432\n10\n\nHello World!\n\n0\n0\n1\n"

func TestBasic(t *testing.T) {
	mod := ir.NewModule()
	GenerateInModule(mod)
	generateTestMain(mod)
	ir := mod.String()
	fmt.Printf("%s", ir)
	err := ioutil.WriteFile("./runtime.ll", []byte(ir), 0644)
	defer os.Remove("runtime.ll")
	defer os.Remove("runtime")
	assert.Nil(t, err)

	clangArgs := []string{
		"runtime.ll",
		"-Wno-override-module", // Disable override target triple warnings
		"-o", "runtime",        // Output path
		"-O3",
	}
	cmd := exec.Command("clang", clangArgs...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println(string(output))
		assert.Nil(t, err)
	}

	if len(output) > 0 {
		fmt.Println(string(output))
	}

	cmd = exec.Command("./runtime")
	output, _ = cmd.CombinedOutput()
	soutput := string(output)
	fmt.Println(soutput)
	assert.Equal(t, basicOutput, soutput)
}
