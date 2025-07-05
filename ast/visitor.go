package ast

import "fmt"

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
	// VisitDocumentation visits a [Documentation] node.
	VisitDocumentation(*Documentation) error
	// VisitBlockComment visits a [BlockComment] node.
	VisitBlockComment(*BlockComment) error
	// VisitCommentStatement visits a [CommentStatement] node.
	VisitCommentStatement(*CommentStatement) error
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
	// VisitExpressionStatement visits an [ExpressionStatement] node.
	VisitExpressionStatement(*ExpressionStatement) error
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
	// VisitErrorStatement visits an [ErrorStatement] node.
	VisitErrorStatement(*ErrorStatement) error
}

func visitComments(v Visitor, node interface{ Comments() *Comments }) error {
	if err := visitPrefixComments(v, node); err != nil {
		return err
	}
	return visitSuffixComments(v, node)
}

func visitPrefixComments(v Visitor, node interface{ Comments() *Comments }) error {
	comments := node.Comments()
	if comments == nil {
		return nil
	}
	for _, c := range comments.Prefix() {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	return nil
}

func visitSuffixComments(v Visitor, node interface{ Comments() *Comments }) error {
	comments := node.Comments()
	if comments == nil {
		return nil
	}
	for _, c := range comments.Suffix() {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	return nil
}
