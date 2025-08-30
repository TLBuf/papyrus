package parser

import (
	"github.com/TLBuf/papyrus/issue"
)

var (
	intenalInvalidState = issue.NewInternal(
		"PRSR0001",
		`Parser is in an invalid state.`,
	)
	errorExpectedScriptName = issue.NewError(
		"PRSR1003",
		`Expected script to start with "ScriptName`,
	)
	errorExpectedScriptNameIdent = issue.NewError(
		"PRSR1004",
		`Expected 'ScriptName' to be followed by an identifier, the name of the script.`,
	)
	errorExpectedExtendsIdent = issue.NewError(
		"PRSR1005",
		`Expected 'Extends' to be followed by an identifier, the name of the parent script.`,
	)
	errorExpectedScriptStatementKeyword = issue.NewError(
		"PRSR1006",
		`Expected 'Property', 'Function', or an identfier for a variable to follow a type literal defined at the script level.`,
	)
	errorExpectedScriptStatement = issue.NewError(
		"PRSR1007",
		`Expected 'Import', 'Event', 'Function', 'State', 'Auto' (preceding 'State'), or a type literal to start a statement at the script level.`,
	)
	errorExpectedImportIdent = issue.NewError(
		"PRSR1008",
		`Expected 'Import' to be followed by an identifier, the name of the script being imported.`,
	)
	errorExpectedImportEnd = issue.NewError(
		"PRSR1009",
		`Expected a newline (or the end of the file) to follow the identifier after 'Import'.`,
	)
	errorExpectedAutoStateKeyword = issue.NewError(
		"PRSR1010",
		`Expected 'State' to follow 'Auto' defined at the script level.`,
	)
	errorExpectedStateIdent = issue.NewError(
		"PRSR1011",
		`Expected 'State' to be followed by an identifier, the name of the state.`,
	)
	errorExpectedStateStatementKeyword = issue.NewError(
		"PRSR1012",
		`Expected 'Function' to follow a type literal defined within a state.`,
	)
	errorExpectedStateStatement = issue.NewError(
		"PRSR1013",
		`Expected 'Event', 'Function', or a type literal to start a statement within a state.`,
	)
	errorUnclosedState = issue.NewError(
		"PRSR1014",
		`Reached the end of the file while parsing a state; did you forget a closing 'EndState'?`,
	)
	errorStateEnd = issue.NewError(
		"PRSR1015",
		`Expected a newline (or the end of the file) to follow 'EndState'.`,
	)
	errorExpectedEventIdent = issue.NewError(
		"PRSR1016",
		`Expected 'Event' to be followed by an identifier, the name of the event.`,
	)
	errorExpectedEventOpenParen = issue.NewError(
		"PRSR1017",
		`Expected '(' to follow the identifier after 'Event'`,
	)
	errorUnclosedEvent = issue.NewError(
		"PRSR1018",
		`Reached the end of the file while parsing an event; did you forget a closing 'EndEvent'?`,
	)
	errorEventEnd = issue.NewError(
		"PRSR1019",
		`Expected a newline (or the end of the file) to follow 'EndEvent'.`,
	)
	errorExpectedFunctionIdent = issue.NewError(
		"PRSR1020",
		`Expected 'Function' to be followed by an identifier, the name of the function.`,
	)
	errorExpectedFunctionOpenParen = issue.NewError(
		"PRSR1021",
		`Expected '(' to follow the identifier after 'Function'`,
	)
	errorUnclosedFunction = issue.NewError(
		"PRSR1022",
		`Reached the end of the file while parsing a function; did you forget a closing 'EndFunction'?`,
	)
	errorFunctionEnd = issue.NewError(
		"PRSR1023",
		`Expected a newline (or the end of the file) to follow 'EndFunction'.`,
	)
	errorUnclosedParamListEOF = issue.NewError(
		"PRSR1024",
		`Reached the end of the file while parsing a parameter list; did you forget a closing ')'?`,
	)
	errorUnclosedParamListNewline = issue.NewError(
		"PRSR1025",
		`Reached the end of the line while parsing a parameter list; did you forget a closing ')'?`,
	)
	errorExpectedParamTypeLiteral = issue.NewError(
		"PRSR1026",
		`Expected a type literal to start a parameter.`,
	)
	errorExpectedParamIdent = issue.NewError(
		"PRSR1027",
		`Expected the parameter's type literal to be followed by an identifier, the name of the parameter.`,
	)
	errorExpectedParamLiteral = issue.NewError(
		"PRSR1028",
		`Expected a literal expression to follow the '=' when defining a parameter's default value.`,
	)
	errorExpectedFunctionStatementExpr = issue.NewError(
		"PRSR1029",
		`Expected an expression statement.`,
	)
	errorExpectedFunctionVariableIdent = issue.NewError(
		"PRSR1030",
		`Expected type literal of a variable definition to be followed by an identifier, the name of the variable.`,
	)
	errorExpectedFunctionVariableExpr = issue.NewError(
		"PRSR1031",
		`Expected an expression to follow the '=' when defining a variable's initial value.`,
	)
	errorExpectedAssignmentAssigneeExpr = issue.NewError(
		"PRSR1032",
		`Expected an expression on the left side of the '=' in an assignment statement.`,
	)
	errorExpectedAssignmentValueExpr = issue.NewError(
		"PRSR1033",
		`Expected an expression on the right side of the '=' in an assignment statement.`,
	)
	errorExpectedReturnExpr = issue.NewError(
		"PRSR1034",
		`Expected an expression (or newline) to follow 'Return'.`,
	)
	errorExpectedIfConditionExpr = issue.NewError(
		"PRSR1035",
		`Expected an expression to follow 'If'.`,
	)
	errorUnclosedIf = issue.NewError(
		"PRSR1036",
		`Reached the end of the file while parsing an if block; did you forget a closing 'EndIf'?`,
	)
	errorExpectedElseIfConditionExpr = issue.NewError(
		"PRSR1037",
		`Expected an expression to follow 'ElseIf'; did you mean 'Else'?.`,
	)
	errorUnclosedElseIf = issue.NewError(
		"PRSR1038",
		`Reached the end of the file while parsing an else-if block; did you forget a closing 'EndIf'?`,
	)
	errorUnclosedElse = issue.NewError(
		"PRSR1039",
		`Reached the end of the file while parsing an else block; did you forget a closing 'EndIf'?`,
	)
	errorExpectedWhileExpr = issue.NewError(
		"PRSR1040",
		`Expected an expression to follow 'While'.`,
	)
	errorUnclosedWhile = issue.NewError(
		"PRSR1041",
		`Reached the end of the file while parsing a while block; did you forget a closing 'EndWhile'?`,
	)
	errorExpectedPropertyIdent = issue.NewError(
		"PRSR1042",
		`Expected 'Property' to be followed by an identifier, the name of the property.`,
	)
	errorExpectedPropertyLiteral = issue.NewError(
		"PRSR1043",
		`Expected a literal expression to follow the '=' when defining a property's default value.`,
	)
	errorExpectedPropertyReadOnlyValue = issue.NewError(
		"PRSR1044",
		`Expected property defined as 'AutoReadOnly' to define a default value.`,
	)
	errorExpectedFullPropertyStatement = issue.NewError(
		"PRSR1045",
		`Expected 'Function', or a type literal to start a function within a full property.`,
	)
	errorExpectedFullPropertyKeywordType = issue.NewError(
		"PRSR1046",
		`Expected 'Function' to follow a type literal defined within a full property.`,
	)
	errorExpectedFullPropertyGetReturnType = issue.NewError(
		"PRSR1047",
		`Expected the 'Get' function defined within a full property to have a return type.`,
	)
	errorExpectedFullPropertyGetParams = issue.NewError(
		"PRSR1048",
		`Expected the 'Get' function defined within a full property to have zero parameters.`,
	)
	errorExpectedFullPropertyGetDuplicate = issue.NewError(
		"PRSR1049",
		`Expected exactly one 'Get' function defined within a full property.`,
	)
	errorExpectedFullPropertySetReturnType = issue.NewError(
		"PRSR1050",
		`Expected the 'Set' function defined within a full property to have no return type.`,
	)
	errorExpectedFullPropertySetParams = issue.NewError(
		"PRSR1051",
		`Expected the 'Set' function defined within a full property to have one parameter.`,
	)
	errorExpectedFullPropertySetDuplicate = issue.NewError(
		"PRSR1052",
		`Expected exactly one 'Set' function defined within a full property.`,
	)
	errorExpectedFullPropertyGetOrSet = issue.NewError(
		"PRSR1053",
		`Expected the function(s) defined within a full property to be named 'Get' or 'Set'.`,
	)
	errorUnclosedFullProperty = issue.NewError(
		"PRSR1054",
		`Reached the end of the file while parsing a full property; did you forget a closing 'EndProperty'?`,
	)
	errorFullPropertyExtra = issue.NewError(
		"PRSR1055",
		`Expected 'EndProperty' to follow the function definitions of a full property.`,
	)
	errorFullPropertyEnd = issue.NewError(
		"PRSR1056",
		`Expected a newline (or the end of the file) to follow 'EndProperty'.`,
	)
	errorExpectedScriptVariableIdent = issue.NewError(
		"PRSR1057",
		`Expected type literal of a variable definition to be followed by an identifier, the name of the variable.`,
	)
	errorExpectedScriptVariableLiteral = issue.NewError(
		"PRSR1058",
		`Expected a literal expression to follow the '=' when defining a script variable's initial value.`,
	)
	errorExpectedScriptVariableEnd = issue.NewError(
		"PRSR1059",
		`Expected a newline (or the end of the file) to follow a script variable definition.`,
	)
	errorExpectedBinaryExpr = issue.NewError(
		"PRSR1060",
		`Expected an expression to follow binary operator.`,
	)
	errorExpectedUnaryExpr = issue.NewError(
		"PRSR1061",
		`Expected an expression to follow unary operator.`,
	)
	errorExpectedCastTypeLiteral = issue.NewError(
		"PRSR1062",
		`Expected the cast operator, 'As', to be followed by a type literal.`,
	)
	errorExpectedAccessIdent = issue.NewError(
		"PRSR1063",
		`Expected the access operator, '.', to be followed by an indentifer, the name of the element being accessed.`,
	)
	errorExpectedIndexExpr = issue.NewError(
		"PRSR1064",
		`Expected an expression to follow the '[' in an index expression.`,
	)
	errorExpectedIndexCloseBracket = issue.NewError(
		"PRSR1065",
		`Expected a ']' after index expression.`,
	)
	errorUnclosedArgListEOF = issue.NewError(
		"PRSR1066",
		`Reached the end of the file while parsing an argument list; did you forget a closing ')'?`,
	)
	errorUnclosedArgListNewline = issue.NewError(
		"PRSR1067",
		`Reached the end of the line while parsing an argument list; did you forget a closing ')'?`,
	)
	errorExpectedArgExpr = issue.NewError(
		"PRSR1068",
		`Expected an expression for argument value.`,
	)
	errorExpectedArrayCreationTypeLiteral = issue.NewError(
		"PRSR1069",
		`Expected the array creation operator, 'New', to be followed by a type literal.`,
	)
	errorExpectedArrayCreationOpenBracket = issue.NewError(
		"PRSR1070",
		`Expected the type literal following 'New' to be followed by a '['.`,
	)
	errorExpectedArrayCreationInt = issue.NewError(
		"PRSR1071",
		`Expected the '[' in an array creation expression to be followed by an interger literal; Papyrus does not support dynamic array sizes.`,
	)
	errorExpectedArrayCreationCloseBracket = issue.NewError(
		"PRSR1072",
		`Expected a ']' to follow the integer literal in an array creation expression.`,
	)
	errorExpectedParenExpr = issue.NewError(
		"PRSR1073",
		`Expected an expression to follow '('.`,
	)
	errorExpectedParenClose = issue.NewError(
		"PRSR1074",
		`Expected ')' after expression; did you forget to close a parenthetical?`,
	)
)
