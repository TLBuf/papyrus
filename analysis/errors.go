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
	errorUnaryNegationNotNumeric = issue.NewError(
		"CHKR1027",
		`Only numeric expressions can be negated with the minus operator, '-'.`,
	)
	errorUnaryNegationNotBool = issue.NewError(
		"CHKR1028",
		`Only expressions of a type assignable to Bool can be negated with the logical not operator, '!'.`,
	)
	errorCastNotConvertible = issue.NewError(
		"CHKR1029",
		`Expression is of a type that is incompatible with the desired cast type.`,
	)
	errorIfConditionNotBool = issue.NewError(
		"CHKR1030",
		`The If condition expression must be of a type assignable to Bool.`,
	)
	errorElseIfConditionNotBool = issue.NewError(
		"CHKR1031",
		`The ElseIf condition expression must be of a type assignable to Bool.`,
	)
	errorWhileConditionNotBool = issue.NewError(
		"CHKR1032",
		`The While condition expression must be of a type assignable to Bool.`,
	)
	errorVariableTypeMismatch = issue.NewError(
		"CHKR1033",
		`The expression assigned to a variable must be assignable to the variable's type.`,
	)
	errorFunctionReturnValueMissing = issue.NewError(
		"CHKR1034",
		`Return statements in a Function that defines a return type must specify a return value of a type assignable to the Function's return type.`,
	)
	errorFunctionReturnValueUnexpected = issue.NewError(
		"CHKR1035",
		`Return statements in a Function that does not define a return type must not specify a return value.`,
	)
	errorFunctionReturnTypeMismatch = issue.NewError(
		"CHKR1036",
		`Return statements in a Function that defines a return type must specify a return value of a type assignable to the Function's return type.`,
	)
	errorEventReturnValueUnexpected = issue.NewError(
		"CHKR1037",
		`Return statements in an Event must not specify a return value.`,
	)
	errorParameterDefaultValueTypeMismatch = issue.NewError(
		"CHKR1038",
		`Parameter default values must be of a type assignable to the parameter's type.`,
	)
	errorAssignmentTypeMismatch = issue.NewError(
		"CHKR1039",
		`Expressions assigned to variables must be of a type assignable to the variable's type.`,
	)
	errorAssignmentArithmeticAssigneeNotNumeric = issue.NewError(
		"CHKR1040",
		`Variables assigned to with an arithmetic assignment operator ('+=', '-=', '*=', or '/=') must be of a numeric type.`,
	)
	errorAssignmentArithmeticValueNotNumeric = issue.NewError(
		"CHKR1041",
		`Expressions assigned with an arithmetic assignment operator ('+=', '-=', '*=', or '/=') must be of a numeric type.`,
	)
	errorAssignmentModuloAssigneeNotInt = issue.NewError(
		"CHKR1042",
		`Variables assigned to with the modulo assignment operator ('%=') must be of type Int.`,
	)
	errorAssignmentModuloValueNotInt = issue.NewError(
		"CHKR1043",
		`Expressions assigned with the modulo assignment operator ('%=') must be of type Int.`,
	)
	errorBinaryOperandsNotEquatable = issue.NewError(
		"CHKR1044",
		`Expressions compared with equality operators ('==' or '!=') must be of equatable types.`,
	)
	errorBinaryOperandsNotComparable = issue.NewError(
		"CHKR1045",
		`Expressions compared with comparison operators ('>', '>=', '<', or '<=') must be of comparable types.`,
	)
	errorBinaryLogicalOperandNotBool = issue.NewError(
		"CHKR1046",
		`Expressions combined with logical operators ('&&' or '||') must be of a type assignable to Bool.`,
	)
	errorBinaryArithmeticOperandNotNumeric = issue.NewError(
		"CHKR1047",
		`Expressions combined with arithmetic operators ('+', '-', '*', or '/') must be of a numeric type.`,
	)
	errorBinaryModuloOperandNotInt = issue.NewError(
		"CHKR1048",
		`Expressions combined with the modulo operator ('%') must be of type Int.`,
	)
)
