// (c) Copyright 2016 Hewlett Packard Enterprise Development LP
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package core

import (
	"encoding/json"
	"fmt"
	"go/ast"
	"os"
)

// Score type used by severity and confidence values
type Score int

const (
	Low    Score = iota // Low value
	Medium              // Medium value
	High                // High value
)

// An Issue is returnd by a GAS rule if it discovers an issue with the scanned code.
type Issue struct {
	Severity   Score  `json:"severity"`   // issue severity (how problematic it is)
	Confidence Score  `json:"confidence"` // issue confidence (how sure we are we found it)
	What       string `json:"details"`    // Human readable explanation
	File       string `json:"file"`       // File name we found it in
	Code       string `json:"code"`       // Impacted code line
	Line       int    `json:"line"`       // Line number in file
}

// MetaData is embedded in all GAS rules. The Severity, Confidence and What message
// will be passed tbhrough to reported issues.
type MetaData struct {
	Severity   Score
	Confidence Score
	What       string
}

// MarshalJSON is used convert a Score object into a JSON representation
func (c Score) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}


// String converts a Score into a string
func (c Score) String() string {
	switch c {
	case High:
		return "xxx"
	case Medium:
		return "xxx"
	case Low:
		return "xxx"
	}
	return "UNDEFINED"
}

func codeSnippet(file *os.File, start int64, end int64, n ast.Node) (string, error) {
	if n == nil {
		return "", fmt.Errorf("Invalid AST node provided")
	}

	size := (int)(end - start) // Go bug, os.File.Read should return int64 ...
	file.Seek(start, 0)

	buf := make([]byte, size)
	if nread, err := file.Read(buf); err != nil || nread != size {
		return "", fmt.Errorf("Unable to read code")
	}
	return string(buf), nil
}

// NewIssue creates a new Issue
func NewIssue(ctx *Context, node ast.Node, desc string, severity Score, confidence Score) *Issue {
	var code string
	fobj := ctx.FileSet.File(node.Pos())
	name := fobj.Name()
	line := fobj.Line(node.Pos())

	if file, err := os.Open(fobj.Name()); err == nil {
		defer file.Close()
		s := (int64)(fobj.Position(node.Pos()).Offset) // Go bug, should be int64
		e := (int64)(fobj.Position(node.End()).Offset) // Go bug, should be int64
		code, err = codeSnippet(file, s, e, node)
		if err != nil {
			code = err.Error()
		}
	}

	return &Issue{
		File:       name,
		Line:       line,
		What:       desc,
		Confidence: confidence,
		Severity:   severity,
		Code:       code,
	}
}
