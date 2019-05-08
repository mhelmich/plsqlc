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
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

var printIR = false
var deleteTmpFile = true

func TestBasic(t *testing.T) {
	Compile("../examples/test.sql", "./test", false, false)
	defer os.Remove("test")
}

var fixture1Output = "Hello World!\n99\n"

func TestFixture1(t *testing.T) {
	Compile("./test01.sql", "./test", printIR, deleteTmpFile)
	output, err := executeBinary("./test")
	assert.Equal(t, fixture1Output, output)
	assert.Nil(t, err)
	// err = os.Remove("./_temp_llvm_.ll")
	// assert.Nil(t, err)
	err = os.Remove("./test")
	assert.Nil(t, err)
}

var fixture2Output = "is_narf\n"

func TestFixture2(t *testing.T) {
	Compile("./test02.sql", "./test", printIR, deleteTmpFile)
	output, err := executeBinary("./test")
	assert.Equal(t, fixture2Output, output)
	assert.Nil(t, err)
	// err = os.Remove("./_temp_llvm_.ll")
	// assert.Nil(t, err)
	err = os.Remove("./test")
	assert.Nil(t, err)
}

var fixture3Output = "15\n14\n13\n12\n11\n"

func TestFixture3(t *testing.T) {
	Compile("./test03.sql", "./test", printIR, deleteTmpFile)
	output, err := executeBinary("./test")
	assert.Equal(t, fixture3Output, output)
	assert.Nil(t, err)
	// err = os.Remove("./_temp_llvm_.ll")
	// assert.Nil(t, err)
	err = os.Remove("./test")
	assert.Nil(t, err)
}

var fixture4Output = "10\n"

func TestFixture4(t *testing.T) {
	Compile("./test04.sql", "./test", printIR, deleteTmpFile)
	output, err := executeBinary("./test")
	assert.Equal(t, fixture4Output, output)
	assert.Nil(t, err)
	// err = os.Remove("./_temp_llvm_.ll")
	// assert.Nil(t, err)
	err = os.Remove("./test")
	assert.Nil(t, err)
}

var fixture5Output = "is_15\nend\n"

func TestFixture5(t *testing.T) {
	Compile("./test05.sql", "./test", printIR, deleteTmpFile)
	output, err := executeBinary("./test")
	assert.Equal(t, fixture5Output, output)
	assert.Nil(t, err)
	// err = os.Remove("./_temp_llvm_.ll")
	// assert.Nil(t, err)
	err = os.Remove("./test")
	assert.Nil(t, err)
}

var fixture6Output = "Hello_from_P1!\n"

func TestFixture6(t *testing.T) {
	Compile("./test06.sql", "./test", printIR, deleteTmpFile)
	output, err := executeBinary("./test")
	assert.Equal(t, fixture6Output, output)
	assert.Nil(t, err)
	// err = os.Remove("./_temp_llvm_.ll")
	// assert.Nil(t, err)
	err = os.Remove("./test")
	assert.Nil(t, err)
}

func executeBinary(file string) (string, error) {
	cmd := exec.Command(file)
	output, err := cmd.CombinedOutput()
	return string(output), err
}
