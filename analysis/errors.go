package analysis

import (
	"github.com/TLBuf/papyrus/issue"
)

var (
	intenalInvalidState = issue.NewInternal(
		"CHKR0001",
		`Checker is in an invalid state.`,
	)
	errorScriptUnknownParent = issue.NewError(
		"CHKR1001",
		`Script extends an unknown script; was the script loaded?`,
	)
	errorScriptCycle = issue.NewError(
		"CHKR1002",
		`Script is a descendant of itself, scripts form an inheritance cycle.`,
	)
	errorScriptNameCollision = issue.NewError(
		"CHKR1003",
		`Script has a name that matches another script.`,
	)
)
