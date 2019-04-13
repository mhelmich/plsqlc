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

package main

import (
	"flag"
	"log"
	"os"

	"github.com/mhelmich/plsqlc/compiler"
)

func main() {
	inFilePath, outFilePath, printIR, deleteIR := parseArgs()

	if !fileExists(inFilePath) {
		log.Panicf("File '%s' doesn't exist", inFilePath)
		return
	}

	compiler.Compile(inFilePath, outFilePath, printIR, deleteIR)
}

func parseArgs() (string, string, bool, bool) {
	inFilePathPtr := flag.String("i", "", "path to input sql file")
	outFilePathPtr := flag.String("o", "", "path to the output binary")
	printIRPtr := flag.Bool("pir", false, "whether or not to print LLVM IR onto the terminal")
	deleteIRPtr := flag.Bool("dir", true, "whether or not to delete intermediate files")
	flag.Parse()

	var inFilePath string
	if *inFilePathPtr == "" {
		inFilePath = flag.Arg(0)
	} else {
		inFilePath = *inFilePathPtr
	}

	var outFilePath string
	if *outFilePathPtr == "" {
		outFilePath = "out"
	} else {
		outFilePath = *outFilePathPtr
	}

	return inFilePath, outFilePath, *printIRPtr, *deleteIRPtr
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}
