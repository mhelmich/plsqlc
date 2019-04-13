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

import "fmt"

type ConstantType int

const (
	StringType ConstantType = iota
	NumericType
)

func NewConstant(v string, t ConstantType) *Constant {
	return &Constant{
		Value: v,
		T:     t,
	}
}

type Constant struct {
	Value string
	T     ConstantType
}

func (c *Constant) String() string {
	return fmt.Sprintf("constant %s %d", c.Value, c.T)
}
