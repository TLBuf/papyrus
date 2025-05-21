package ast

import (
	"fmt"
)

// Visitor is a visitor of AST nodes.
type Visitor interface {
	// VisitAccess visits an [Access] node.
	VisitAccess(*Access) error
	// VisitArgument visits an [Argument] node.
	VisitArgument(*Argument) error
	// VisitArrayCreation visits an [ArrayCreation] node.
	VisitArrayCreation(*ArrayCreation) error
	// VisitAssignment visits an [Assignment] node.
	VisitAssignment(*Assignment) error
	// VisitBinary visits a [Binary] node.
	VisitBinary(*Binary) error
	// VisitCall visits a [Call] node.
	VisitCall(*Call) error
	// VisitCast visits a [Cast] node.
	VisitCast(*Cast) error
	// VisitDocComment visits a [DocComment] node.
	VisitDocComment(*DocComment) error
	// VisitBlockComment visits a [BlockComment] node.
	VisitBlockComment(*BlockComment) error
	// VisitLineComment visits a [LineComment] node.
	VisitLineComment(*LineComment) error
	// VisitEvent visits an [Event] node.
	VisitEvent(*Event) error
	// VisitFunction visits a [Function] node.
	VisitFunction(*Function) error
	// VisitIdentifier visits an [Identifier] node.
	VisitIdentifier(*Identifier) error
	// VisitIf visits an [If] node.
	VisitIf(*If) error
	// VisitElseIf visits an [ElseIf] node.
	VisitElseIf(*ElseIf) error
	// VisitElse visits an [Else] node.
	VisitElse(*Else) error
	// VisitImport visits an [Import] node.
	VisitImport(*Import) error
	// VisitIndex visits an [Index] node.
	VisitIndex(*Index) error
	// VisitBoolLiteral visits a [BoolLiteral] node.
	VisitBoolLiteral(*BoolLiteral) error
	// VisitIntLiteral visits an [IntLiteral] node.
	VisitIntLiteral(*IntLiteral) error
	// VisitFloatLiteral visits a [FloatLiteral] node.
	VisitFloatLiteral(*FloatLiteral) error
	// VisitStringLiteral visits a [StringLiteral] node.
	VisitStringLiteral(*StringLiteral) error
	// VisitNoneLiteral visits a [NoneLiteral] node.
	VisitNoneLiteral(*NoneLiteral) error
	// VisitParameter visits a [Parameter] node.
	VisitParameter(*Parameter) error
	// VisitParenthetical visits a [Parenthetical] node.
	VisitParenthetical(*Parenthetical) error
	// VisitProperty visits a [Property] node.
	VisitProperty(*Property) error
	// VisitReturn visits a [Return] node.
	VisitReturn(*Return) error
	// VisitScript visits a [Script] node.
	VisitScript(*Script) error
	// VisitState visits a [State] node.
	VisitState(*State) error
	// VisitToken visits a [Token] node.
	VisitToken(*Token) error
	// VisitTypeLiteral visits a [TypeLiteral] node.
	VisitTypeLiteral(*TypeLiteral) error
	// VisitUnary visits a [Unary] node.
	VisitUnary(*Unary) error
	// VisitScriptVariable visits a [ScriptVariable] node.
	VisitScriptVariable(*ScriptVariable) error
	// VisitFunctionVariable visits a [FunctionVariable] node.
	VisitFunctionVariable(*FunctionVariable) error
	// VisitWhile visits a [While] node.
	VisitWhile(*While) error
	// VisitErrorScriptStatement visits an [ErrorScriptStatement] node.
	VisitErrorScriptStatement(*ErrorScriptStatement) error
	// VisitErrorFunctionStatement visits an [ErrorFunctionStatement] node.
	VisitErrorFunctionStatement(*ErrorFunctionStatement) error
}

// VisitScriptStatement calls the appropriate Visit
// method on a [Visitor] for a [ScriptStatement].
func VisitScriptStatement(v Visitor, s ScriptStatement) error {
	switch s := s.(type) {
	case *Event:
		return v.VisitEvent(s)
	case *Function:
		return v.VisitFunction(s)
	case *Import:
		return v.VisitImport(s)
	case *Property:
		return v.VisitProperty(s)
	case *ScriptVariable:
		return v.VisitScriptVariable(s)
	case *ErrorScriptStatement:
		return v.VisitErrorScriptStatement(s)
	default:
		return fmt.Errorf("unsupported ScriptStatement implementation: %T", s)
	}
}

// VisitFunctionStatement calls the appropriate Visit
// method on a [Visitor] for a [FunctionStatement].
func VisitFunctionStatement(v Visitor, s FunctionStatement) error {
	switch s := s.(type) {
	case *Access:
		return v.VisitAccess(s)
	case *ArrayCreation:
		return v.VisitArrayCreation(s)
	case *Assignment:
		return v.VisitAssignment(s)
	case *Binary:
		return v.VisitBinary(s)
	case *Call:
		return v.VisitCall(s)
	case *Cast:
		return v.VisitCast(s)
	case *Identifier:
		return v.VisitIdentifier(s)
	case *If:
		return v.VisitIf(s)
	case *Index:
		return v.VisitIndex(s)
	case *BoolLiteral:
		return v.VisitBoolLiteral(s)
	case *IntLiteral:
		return v.VisitIntLiteral(s)
	case *FloatLiteral:
		return v.VisitFloatLiteral(s)
	case *StringLiteral:
		return v.VisitStringLiteral(s)
	case *NoneLiteral:
		return v.VisitNoneLiteral(s)
	case *Parenthetical:
		return v.VisitParenthetical(s)
	case *Return:
		return v.VisitReturn(s)
	case *FunctionVariable:
		return v.VisitFunctionVariable(s)
	case *Unary:
		return v.VisitUnary(s)
	case *While:
		return v.VisitWhile(s)
	case *ErrorFunctionStatement:
		return v.VisitErrorFunctionStatement(s)
	default:
		return fmt.Errorf("unsupported FunctionStatement implementation: %T", s)
	}
}

// VisitExpression calls the appropriate Visit
// method on a [Visitor] for an [Expression].
func VisitExpression(v Visitor, e Expression) error {
	switch e := e.(type) {
	case *Access:
		return v.VisitAccess(e)
	case *ArrayCreation:
		return v.VisitArrayCreation(e)
	case *Binary:
		return v.VisitBinary(e)
	case *Call:
		return v.VisitCall(e)
	case *Cast:
		return v.VisitCast(e)
	case *Identifier:
		return v.VisitIdentifier(e)
	case *Index:
		return v.VisitIndex(e)
	case *BoolLiteral:
		return v.VisitBoolLiteral(e)
	case *IntLiteral:
		return v.VisitIntLiteral(e)
	case *FloatLiteral:
		return v.VisitFloatLiteral(e)
	case *StringLiteral:
		return v.VisitStringLiteral(e)
	case *NoneLiteral:
		return v.VisitNoneLiteral(e)
	case *Parenthetical:
		return v.VisitParenthetical(e)
	case *Unary:
		return v.VisitUnary(e)
	default:
		return fmt.Errorf("unsupported Expression implementation: %T", e)
	}
}

// VisitLiteral calls the appropriate Visit
// method on a [Visitor] for a [Literal].
func VisitLiteral(v Visitor, l Literal) error {
	switch l := l.(type) {
	case *BoolLiteral:
		return v.VisitBoolLiteral(l)
	case *IntLiteral:
		return v.VisitIntLiteral(l)
	case *FloatLiteral:
		return v.VisitFloatLiteral(l)
	case *StringLiteral:
		return v.VisitStringLiteral(l)
	case *NoneLiteral:
		return v.VisitNoneLiteral(l)
	default:
		return fmt.Errorf("unsupported Literal implementation: %T", l)
	}
}

// VisitInvokable calls the appropriate Visit
// method on a [Visitor] for a [Invokable].
func VisitInvokable(v Visitor, i Invokable) error {
	switch i := i.(type) {
	case *Event:
		return v.VisitEvent(i)
	case *Function:
		return v.VisitFunction(i)
	default:
		return fmt.Errorf("unsupported Invokable implementation: %T", i)
	}
}

// VisitLooseComment calls the appropriate Visit
// method on a [Visitor] for a [LooseComment].
func VisitLooseComment(v Visitor, c LooseComment) error {
	switch c := c.(type) {
	case *BlockComment:
		return v.VisitBlockComment(c)
	case *LineComment:
		return v.VisitLineComment(c)
	default:
		return fmt.Errorf("unsupported LooseComment implementation: %T", c)
	}
}

// VisitError calls the appropriate Visit
// method on a [Visitor] for an [Error].
func VisitError(v Visitor, e Error) error {
	switch e := e.(type) {
	case *ErrorScriptStatement:
		return v.VisitErrorScriptStatement(e)
	case *ErrorFunctionStatement:
		return v.VisitErrorFunctionStatement(e)
	default:
		return fmt.Errorf("unsupported Error implementation: %T", e)
	}
}

func visitLeadingComments(v Visitor, t Trivia) error {
	for _, c := range t.LeadingComments {
		if err := VisitLooseComment(v, c); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range t.PrefixComments {
		if err := VisitLooseComment(v, c); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	return nil
}

func visitTrailingComments(v Visitor, t Trivia) error {
	for _, c := range t.SuffixComments {
		if err := VisitLooseComment(v, c); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range t.TrailingComments {
		if err := VisitLooseComment(v, c); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	return nil
}
