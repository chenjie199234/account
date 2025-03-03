package model

import (
	"os"

	"github.com/chenjie199234/Corelib/util/name"
)

// Warning!!!!!!!!!!!
// This file is readonly!
// Don't modify this file!

const pkg = "github.com/chenjie199234/account"
const Name = "account"

var Group = os.Getenv("GROUP")
var Project = os.Getenv("PROJECT")

func init() {
	if Group == "" || Group == "<GROUP>" {
		panic("missing env:GROUP")
	}
	if Project == "" || Project == "<PROJECT>" {
		panic("missing env:PROJECT")
	}
	if e := name.SetSelfFullName(Project, Group, Name); e != nil {
		panic(e)
	}
}

