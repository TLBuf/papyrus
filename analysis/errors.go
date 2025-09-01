package analysis

import (
	"github.com/TLBuf/papyrus/issue"
)

var (
	internalInvalidState = issue.NewInternal(
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
	errorStateNameCollision = issue.NewError(
		"CHKR1004",
		`Another state with the same name already exists in the script.`,
	)
	errorValueNameCollision = issue.NewError(
		"CHKR1005",
		`Another property or variable with the same name already exists in this scope.`,
	)
	errorFunctionNameCollision = issue.NewError(
		"CHKR1006",
		`Another function or event with the same name already exists in the script.`,
	)
	errorParameterNameCollision = issue.NewError(
		"CHKR1007",
		`Another parameter with the same name is already defined for this function or event.`,
	)
	errorBoolParseLiteral = issue.NewError(
		"CHKR1008",
		`A Bool literal could not be parsed as a boolean value.`,
	)
	errorIntParseLiteral = issue.NewError(
		"CHKR1009",
		`An Int literal could not be parsed as a 32-bit signed integer value.`,
	)
	errorFloatParseLiteral = issue.NewError(
		"CHKR1010",
		`A Float literal could not be parsed as a 32-bit floating-point value.`,
	)
	errorStringParseLiteral = issue.NewError(
		"CHKR1011",
		`A String literal could not be parsed as a UTF-8 string value.`,
	)
	errorInvalidArrayAccess = issue.NewError(
		"CHKR1012",
		`An array value has an invalid access; arrays can only have 'Length' accessed via the dot operator.`,
	)
	errorUnknownFunction = issue.NewError(
		"CHKR1013",
		`Access operator references an unknown function.`,
	)
	errorCannotCallEvent = issue.NewError(
		"CHKR1014",
		`Access operator references an event for a call, but events cannot be called like functions.`,
	)
	errorUnknownProperty = issue.NewError(
		"CHKR1015",
		`Access operator references an unknown property.`,
	)
	errorCannotAccessVariable = issue.NewError(
		"CHKR1016",
		`Access operator references a variable, but variable cannot be referenced outside the script that defines them.`,
	)
	errorCannotAccessBool = issue.NewError(
		"CHKR1017",
		`Bool values cannot be the target of an access operation.`,
	)
	errorCannotAccessInt = issue.NewError(
		"CHKR1018",
		`Int values cannot be the target of an access operation.`,
	)
	errorCannotAccessFloat = issue.NewError(
		"CHKR1019",
		`Float values cannot be the target of an access operation.`,
	)
	errorCannotAccessString = issue.NewError(
		"CHKR1020",
		`String values cannot be the target of an access operation.`,
	)
	errorCannotAccessFunction = issue.NewError(
		"CHKR1021",
		`Functions cannot be the target of an access operation.`,
	)
	errorCannotAccessEvent = issue.NewError(
		"CHKR1022",
		`Events cannot be the target of an access operation.`,
	)
	errorCannotAccessNone = issue.NewError(
		"CHKR1023",
		`'None' cannot be the target of an access operation.`,
	)
	errorTypeReferencesUnknownScript = issue.NewError(
		"CHKR1024",
		`Type references an unknown script; was the script loaded?`,
	)
	errorIndexTargetNotArray = issue.NewError(
		"CHKR1025",
		`Non-array values cannot be indexed.`,
	)
	errorIndexNotInt = issue.NewError(
		"CHKR1026",
		`The expression used to index an array must evaluate to an integer.`,
	)
	errorNegationNotNumeric = issue.NewError(
		"CHKR1027",
		`Only numeric expressions can be negated with the minus operator, '-'.`,
	)
	errorCastNotConvertible = issue.NewError(
		"CHKR1028",
		`Expression is of a typed that is incompatible with the desired cast type.`,
	)
)
