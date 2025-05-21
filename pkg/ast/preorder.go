package ast

import (
	"fmt"
	"slices"
)

// PreorderVisitor is a [Visitor] that calls visit methods on
// a delegate visitor by performing a preorder traversal of nodes.
type PreorderVisitor struct {
	Delegate Visitor
}

// VisitAccess visits the [Access] node then all children nodes and returns an
// error if any call returns an error.
func (v PreorderVisitor) VisitAccess(a *Access) error {
	if err := v.Delegate.VisitAccess(a); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, a.Trivia); err != nil {
		return err
	}
	if err := VisitExpression(v, a.Value); err != nil {
		return fmt.Errorf("value: %w", err)
	}
	if err := v.VisitToken(a.Operator); err != nil {
		return fmt.Errorf("operator: %w", err)
	}
	if err := v.VisitIdentifier(a.Name); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if err := visitTrailingComments(v, a.Trivia); err != nil {
		return err
	}
	return nil
}

// PreorderVisitor visits the [Argument] node then all children nodes and
// returns an error if any call returns an error.
func (v PreorderVisitor) VisitArgument(a *Argument) error {
	if err := v.Delegate.VisitArgument(a); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, a.Trivia); err != nil {
		return err
	}
	if err := v.VisitIdentifier(a.Name); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if err := v.VisitToken(a.Operator); err != nil {
		return fmt.Errorf("operator: %w", err)
	}
	if err := VisitExpression(v, a.Value); err != nil {
		return fmt.Errorf("value: %w", err)
	}
	if err := visitTrailingComments(v, a.Trivia); err != nil {
		return err
	}
	return nil
}

// VisitArrayCreation visits the [ArrayCreation] node then all children nodes
// and returns an error if any call returns an error.
func (v PreorderVisitor) VisitArrayCreation(a *ArrayCreation) error {
	if err := v.Delegate.VisitArrayCreation(a); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, a.Trivia); err != nil {
		return err
	}
	if err := v.VisitToken(a.NewOperator); err != nil {
		return fmt.Errorf("new operator: %w", err)
	}
	if err := v.VisitTypeLiteral(a.Type); err != nil {
		return fmt.Errorf("type: %w", err)
	}
	if err := v.VisitToken(a.Open); err != nil {
		return fmt.Errorf("open: %w", err)
	}
	if err := v.VisitIntLiteral(a.Size); err != nil {
		return fmt.Errorf("size: %w", err)
	}
	if err := v.VisitToken(a.Close); err != nil {
		return fmt.Errorf("close: %w", err)
	}
	if err := visitTrailingComments(v, a.Trivia); err != nil {
		return err
	}
	return nil
}

// VisitAssignment visits the [Assignment] node then all children nodes and
// returns an error if any call returns an error.
func (v PreorderVisitor) VisitAssignment(a *Assignment) error {
	if err := v.Delegate.VisitAssignment(a); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, a.Trivia); err != nil {
		return err
	}
	if err := VisitExpression(v, a.Assignee); err != nil {
		return fmt.Errorf("assignee: %w", err)
	}
	if err := v.VisitToken(a.Operator); err != nil {
		return fmt.Errorf("operator: %w", err)
	}
	if err := VisitExpression(v, a.Value); err != nil {
		return fmt.Errorf("value: %w", err)
	}
	if err := visitTrailingComments(v, a.Trivia); err != nil {
		return err
	}
	return nil
}

// VisitBinary visits the [Binary] node then all children nodes and
// returns an error if any call returns an error.
func (v PreorderVisitor) VisitBinary(b *Binary) error {
	if err := v.Delegate.VisitBinary(b); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, b.Trivia); err != nil {
		return err
	}
	if err := VisitExpression(v, b.LeftOperand); err != nil {
		return fmt.Errorf("left operand: %w", err)
	}
	if err := v.VisitToken(b.Operator); err != nil {
		return fmt.Errorf("operator: %w", err)
	}
	if err := VisitExpression(v, b.RightOperand); err != nil {
		return fmt.Errorf("right operand: %w", err)
	}
	if err := visitTrailingComments(v, b.Trivia); err != nil {
		return err
	}
	return nil
}

// VisitCall visits the [Call] node then all children nodes and returns an
// error if any call returns an error.
func (v PreorderVisitor) VisitCall(c *Call) error {
	if err := v.Delegate.VisitCall(c); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, c.Trivia); err != nil {
		return err
	}
	if err := VisitExpression(v, c.Reciever); err != nil {
		return fmt.Errorf("reciever: %w", err)
	}
	if err := v.VisitToken(c.Open); err != nil {
		return fmt.Errorf("open: %w", err)
	}
	for _, a := range c.Arguments {
		if err := v.VisitArgument(a); err != nil {
			return fmt.Errorf("argument: %w", err)
		}
	}
	if err := v.VisitToken(c.Close); err != nil {
		return fmt.Errorf("close: %w", err)
	}
	if err := visitTrailingComments(v, c.Trivia); err != nil {
		return err
	}
	return nil
}

// VisitCast visits the [Cast] node then all children nodes and returns an error
// if any call returns an error.
func (v PreorderVisitor) VisitCast(c *Cast) error {
	if err := v.Delegate.VisitCast(c); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, c.Trivia); err != nil {
		return err
	}
	if err := VisitExpression(v, c.Value); err != nil {
		return fmt.Errorf("value: %w", err)
	}
	if err := v.VisitToken(c.Operator); err != nil {
		return fmt.Errorf("operator: %w", err)
	}
	if err := v.VisitTypeLiteral(c.Type); err != nil {
		return fmt.Errorf("type: %w", err)
	}
	if err := visitTrailingComments(v, c.Trivia); err != nil {
		return err
	}
	return nil
}

// VisitDocComment visits the [DocComment] node then all children nodes and
// returns an error if any call returns an error.
func (v PreorderVisitor) VisitDocComment(c *DocComment) error {
	if err := v.Delegate.VisitDocComment(c); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := v.VisitToken(c.Open); err != nil {
		return fmt.Errorf("open: %w", err)
	}
	if err := v.VisitToken(c.Text); err != nil {
		return fmt.Errorf("text: %w", err)
	}
	if err := v.VisitToken(c.Close); err != nil {
		return fmt.Errorf("type: %w", err)
	}
	return nil
}

// VisitBlockComment visits the [BlockComment] node then all children nodes and
// returns an error if any call returns an error.
func (v PreorderVisitor) VisitBlockComment(c *BlockComment) error {
	if err := v.Delegate.VisitBlockComment(c); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := v.VisitToken(c.Open); err != nil {
		return fmt.Errorf("open: %w", err)
	}
	if err := v.VisitToken(c.Text); err != nil {
		return fmt.Errorf("text: %w", err)
	}
	if err := v.VisitToken(c.Close); err != nil {
		return fmt.Errorf("type: %w", err)
	}
	return nil
}

// VisitLineComment visits the [LineComment] node then all children nodes and
// returns an error if any call returns an error.
func (v PreorderVisitor) VisitLineComment(c *LineComment) error {
	if err := v.Delegate.VisitLineComment(c); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := v.VisitToken(c.Open); err != nil {
		return fmt.Errorf("open: %w", err)
	}
	if err := v.VisitToken(c.Text); err != nil {
		return fmt.Errorf("text: %w", err)
	}
	return nil
}

// VisitEvent visits the [Event] node then all children nodes and returns an
// error if any call returns an error.
func (v PreorderVisitor) VisitEvent(e *Event) error {
	if err := v.Delegate.VisitEvent(e); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, e.Trivia); err != nil {
		return err
	}
	if err := v.VisitToken(e.Keyword); err != nil {
		return fmt.Errorf("keyword: %w", err)
	}
	if err := v.VisitIdentifier(e.Name); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if err := v.VisitToken(e.Open); err != nil {
		return fmt.Errorf("open: %w", err)
	}
	for _, p := range e.Parameters {
		if err := v.VisitParameter(p); err != nil {
			return fmt.Errorf("parameter: %w", err)
		}
	}
	if err := v.VisitToken(e.Close); err != nil {
		return fmt.Errorf("close: %w", err)
	}
	for _, t := range e.Native {
		if err := v.VisitToken(t); err != nil {
			return fmt.Errorf("native: %w", err)
		}
	}
	if err := v.VisitDocComment(e.Comment); err != nil {
		return fmt.Errorf("doc comment: %w", err)
	}
	for _, s := range e.Statements {
		if err := VisitFunctionStatement(v, s); err != nil {
			return fmt.Errorf("statement: %w", err)
		}
	}
	if err := v.VisitToken(e.EndKeyword); err != nil {
		return fmt.Errorf("end keyword: %w", err)
	}
	if err := visitTrailingComments(v, e.Trivia); err != nil {
		return err
	}
	return nil
}

// VisitFunction visits the [Function] node then all children nodes and returns
// an error if any call returns an error.
func (v PreorderVisitor) VisitFunction(f *Function) error {
	if err := v.Delegate.VisitFunction(f); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, f.Trivia); err != nil {
		return err
	}
	if f.ReturnType != nil {
		if err := v.VisitTypeLiteral(f.ReturnType); err != nil {
			return fmt.Errorf("return type: %w", err)
		}
	}
	if err := v.VisitToken(f.Keyword); err != nil {
		return fmt.Errorf("keyword: %w", err)
	}
	if err := v.VisitIdentifier(f.Name); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if err := v.VisitToken(f.Open); err != nil {
		return fmt.Errorf("open: %w", err)
	}
	for _, p := range f.Parameters {
		if err := v.VisitParameter(p); err != nil {
			return fmt.Errorf("parameter: %w", err)
		}
	}
	if err := v.VisitToken(f.Close); err != nil {
		return fmt.Errorf("close: %w", err)
	}
	flags := append(f.Global, f.Native...)
	slices.SortFunc(flags, func(a, b *Token) int { return a.Location.ByteOffset - b.Location.ByteOffset })
	for _, t := range flags {
		if err := v.VisitToken(t); err != nil {
			return fmt.Errorf("flag: %w", err)
		}
	}
	if err := v.VisitDocComment(f.Comment); err != nil {
		return fmt.Errorf("doc comment: %w", err)
	}
	for _, s := range f.Statements {
		if err := VisitFunctionStatement(v, s); err != nil {
			return fmt.Errorf("statement: %w", err)
		}
	}
	if err := v.VisitToken(f.EndKeyword); err != nil {
		return fmt.Errorf("end keyword: %w", err)
	}
	if err := visitTrailingComments(v, f.Trivia); err != nil {
		return err
	}
	return nil
}

// VisitIdentifier visits the [Identifier] node then all children nodes and
// returns an error if any call returns an error.
func (v PreorderVisitor) VisitIdentifier(i *Identifier) error {
	if err := v.Delegate.VisitIdentifier(i); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, i.Trivia); err != nil {
		return err
	}
	if err := v.VisitToken(i.Text); err != nil {
		return fmt.Errorf("text: %w", err)
	}
	if err := visitTrailingComments(v, i.Trivia); err != nil {
		return err
	}
	return nil
}

// VisitIf visits the [If] node then all children nodes and returns an error if
// any call returns an error.
func (v PreorderVisitor) VisitIf(i *If) error {
	if err := v.Delegate.VisitIf(i); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, i.Trivia); err != nil {
		return err
	}
	if err := v.VisitToken(i.Keyword); err != nil {
		return fmt.Errorf("keyword: %w", err)
	}
	if err := VisitExpression(v, i.Condition); err != nil {
		return fmt.Errorf("condition: %w", err)
	}
	for _, s := range i.Statements {
		if err := VisitFunctionStatement(v, s); err != nil {
			return fmt.Errorf("statement: %w", err)
		}
	}
	for _, e := range i.ElseIfs {
		if err := v.VisitElseIf(e); err != nil {
			return fmt.Errorf("else if: %w", err)
		}
	}
	if i.Else != nil {
		if err := v.VisitElse(i.Else); err != nil {
			return fmt.Errorf("else: %w", err)
		}
	}
	if err := v.VisitToken(i.EndKeyword); err != nil {
		return fmt.Errorf("end keyword: %w", err)
	}
	if err := visitTrailingComments(v, i.Trivia); err != nil {
		return err
	}
	return nil
}

// VisitElseIf visits the [ElseIf] node then all children nodes and returns an
// error if any call returns an error.
func (v PreorderVisitor) VisitElseIf(e *ElseIf) error {
	if err := v.Delegate.VisitElseIf(e); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, e.Trivia); err != nil {
		return err
	}
	if err := v.VisitToken(e.Keyword); err != nil {
		return fmt.Errorf("keyword: %w", err)
	}
	if err := VisitExpression(v, e.Condition); err != nil {
		return fmt.Errorf("condition: %w", err)
	}
	for _, s := range e.Statements {
		if err := VisitFunctionStatement(v, s); err != nil {
			return fmt.Errorf("statement: %w", err)
		}
	}
	if err := visitTrailingComments(v, e.Trivia); err != nil {
		return err
	}
	return nil
}

// VisitElse visits the [Else] node then all children nodes and returns an error
// if any call returns an error.
func (v PreorderVisitor) VisitElse(e *Else) error {
	if err := v.Delegate.VisitElse(e); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, e.Trivia); err != nil {
		return err
	}
	if err := v.VisitToken(e.Keyword); err != nil {
		return fmt.Errorf("keyword: %w", err)
	}
	for _, s := range e.Statements {
		if err := VisitFunctionStatement(v, s); err != nil {
			return fmt.Errorf("statement: %w", err)
		}
	}
	if err := visitTrailingComments(v, e.Trivia); err != nil {
		return err
	}
	return nil
}

// VisitImport visits the [Import] node then all children nodes and returns an
// error if any call returns an error.
func (v PreorderVisitor) VisitImport(i *Import) error {
	if err := v.Delegate.VisitImport(i); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, i.Trivia); err != nil {
		return err
	}
	if err := v.VisitToken(i.Keyword); err != nil {
		return fmt.Errorf("keyword: %w", err)
	}
	if err := v.VisitIdentifier(i.Name); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if err := visitTrailingComments(v, i.Trivia); err != nil {
		return err
	}
	return nil
}

// VisitIndex visits the [Index] node then all children nodes and returns an
// error if any call returns an error.
func (v PreorderVisitor) VisitIndex(i *Index) error {
	if err := v.Delegate.VisitIndex(i); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, i.Trivia); err != nil {
		return err
	}
	if err := VisitExpression(v, i.Value); err != nil {
		return fmt.Errorf("value: %w", err)
	}
	if err := v.VisitToken(i.Open); err != nil {
		return fmt.Errorf("open: %w", err)
	}
	if err := VisitExpression(v, i.Index); err != nil {
		return fmt.Errorf("index: %w", err)
	}
	if err := v.VisitToken(i.Close); err != nil {
		return fmt.Errorf("close: %w", err)
	}
	if err := visitTrailingComments(v, i.Trivia); err != nil {
		return err
	}
	return nil
}

// VisitBoolLiteral visits the [BoolLiteral] node then all children nodes and
// returns an error if any call returns an error.
func (v PreorderVisitor) VisitBoolLiteral(l *BoolLiteral) error {
	if err := v.Delegate.VisitBoolLiteral(l); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, l.Trivia); err != nil {
		return err
	}
	if err := v.VisitToken(l.Text); err != nil {
		return fmt.Errorf("text: %w", err)
	}
	if err := visitTrailingComments(v, l.Trivia); err != nil {
		return err
	}
	return nil
}

// VisitIntLiteral visits the [IntLiteral] node then all children nodes and
// returns an error if any call returns an error.
func (v PreorderVisitor) VisitIntLiteral(l *IntLiteral) error {
	if err := v.Delegate.VisitIntLiteral(l); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, l.Trivia); err != nil {
		return err
	}
	if err := v.VisitToken(l.Text); err != nil {
		return fmt.Errorf("text: %w", err)
	}
	if err := visitTrailingComments(v, l.Trivia); err != nil {
		return err
	}
	return nil
}

// VisitFloatLiteral visits the [FloatLiteral] node then all children nodes and
// returns an error if any call returns an error.
func (v PreorderVisitor) VisitFloatLiteral(l *FloatLiteral) error {
	if err := v.Delegate.VisitFloatLiteral(l); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, l.Trivia); err != nil {
		return err
	}
	if err := v.VisitToken(l.Text); err != nil {
		return fmt.Errorf("text: %w", err)
	}
	if err := visitTrailingComments(v, l.Trivia); err != nil {
		return err
	}
	return nil
}

// VisitStringLiteral visits the [StringLiteral] node then all children nodes
// and returns an error if any call returns an error.
func (v PreorderVisitor) VisitStringLiteral(l *StringLiteral) error {
	if err := v.Delegate.VisitStringLiteral(l); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, l.Trivia); err != nil {
		return err
	}
	if err := v.VisitToken(l.Text); err != nil {
		return fmt.Errorf("text: %w", err)
	}
	if err := visitTrailingComments(v, l.Trivia); err != nil {
		return err
	}
	return nil
}

// VisitNoneLiteral visits the [NoneLiteral] node then all children nodes and
// returns an error if any call returns an error.
func (v PreorderVisitor) VisitNoneLiteral(l *NoneLiteral) error {
	if err := v.Delegate.VisitNoneLiteral(l); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, l.Trivia); err != nil {
		return err
	}
	if err := v.VisitToken(l.Text); err != nil {
		return fmt.Errorf("text: %w", err)
	}
	if err := visitTrailingComments(v, l.Trivia); err != nil {
		return err
	}
	return nil
}

// VisitParameter visits the [Parameter] node then all children nodes and
// returns an error if any call returns an error.
func (v PreorderVisitor) VisitParameter(p *Parameter) error {
	if err := v.Delegate.VisitParameter(p); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, p.Trivia); err != nil {
		return err
	}
	if err := v.VisitTypeLiteral(p.Type); err != nil {
		return fmt.Errorf("type: %w", err)
	}
	if err := v.VisitIdentifier(p.Name); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if p.Operator != nil {
		if err := v.VisitToken(p.Operator); err != nil {
			return fmt.Errorf("operator: %w", err)
		}
	}
	if p.Value != nil {
		if err := VisitLiteral(v, p.Value); err != nil {
			return fmt.Errorf("value: %w", err)
		}
	}
	if err := visitTrailingComments(v, p.Trivia); err != nil {
		return err
	}
	return nil
}

// VisitParenthetical visits the [Parenthetical] node then all children nodes
// and returns an error if any call returns an error.
func (v PreorderVisitor) VisitParenthetical(p *Parenthetical) error {
	if err := v.Delegate.VisitParenthetical(p); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, p.Trivia); err != nil {
		return err
	}
	if err := v.VisitToken(p.Open); err != nil {
		return fmt.Errorf("open: %w", err)
	}
	if err := VisitExpression(v, p.Value); err != nil {
		return fmt.Errorf("value: %w", err)
	}
	if err := v.VisitToken(p.Close); err != nil {
		return fmt.Errorf("close: %w", err)
	}
	if err := visitTrailingComments(v, p.Trivia); err != nil {
		return err
	}
	return nil
}

// VisitProperty visits the [Property] node then all children nodes and returns
// an error if any call returns an error.
func (v PreorderVisitor) VisitProperty(p *Property) error {
	if err := v.Delegate.VisitProperty(p); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, p.Trivia); err != nil {
		return err
	}
	if err := v.VisitTypeLiteral(p.Type); err != nil {
		return fmt.Errorf("type: %w", err)
	}
	if err := v.VisitToken(p.Keyword); err != nil {
		return fmt.Errorf("keyword: %w", err)
	}
	if err := v.VisitIdentifier(p.Name); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if p.Operator != nil {
		if err := v.VisitToken(p.Operator); err != nil {
			return fmt.Errorf("operator: %w", err)
		}
	}
	if p.Value != nil {
		if err := VisitExpression(v, p.Value); err != nil {
			return fmt.Errorf("value: %w", err)
		}
	}
	if p.Auto != nil {
		if err := v.VisitToken(p.Auto); err != nil {
			return fmt.Errorf("auto: %w", err)
		}
	}
	if p.AutoReadOnly != nil {
		if err := v.VisitToken(p.AutoReadOnly); err != nil {
			return fmt.Errorf("autoreadonly: %w", err)
		}
	}
	flags := append(p.Hidden, p.Conditional...)
	slices.SortFunc(flags, func(a, b *Token) int { return a.Location.ByteOffset - b.Location.ByteOffset })
	for _, t := range flags {
		if err := v.VisitToken(t); err != nil {
			return fmt.Errorf("flag: %w", err)
		}
	}
	if err := v.VisitDocComment(p.Comment); err != nil {
		return fmt.Errorf("doc comment: %w", err)
	}
	first := p.Get
	second := p.Set
	if p.Get != nil && p.Set != nil && second.Location.ByteOffset < first.Location.ByteOffset {
		first = p.Set
		second = p.Get
	}
	if first != nil {
		if err := v.VisitFunction(first); err != nil {
			return fmt.Errorf("function: %w", err)
		}
	}
	if second != nil {
		if err := v.VisitFunction(second); err != nil {
			return fmt.Errorf("function: %w", err)
		}
	}
	if err := v.VisitToken(p.EndKeyword); err != nil {
		return fmt.Errorf("end keyword: %w", err)
	}
	if err := visitTrailingComments(v, p.Trivia); err != nil {
		return err
	}
	return nil
}

// VisitReturn visits the [Return] node then all children nodes and returns an
// error if any call returns an error.
func (v PreorderVisitor) VisitReturn(r *Return) error {
	if err := v.Delegate.VisitReturn(r); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, r.Trivia); err != nil {
		return err
	}
	if err := v.VisitToken(r.Keyword); err != nil {
		return fmt.Errorf("keyword: %w", err)
	}
	if r.Value != nil {
		if err := VisitExpression(v, r.Value); err != nil {
			return fmt.Errorf("value: %w", err)
		}
	}
	return nil
}

// VisitScript visits the [Script] node then all children nodes and returns an
// error if any call returns an error.
func (v PreorderVisitor) VisitScript(s *Script) error {
	if err := v.Delegate.VisitScript(s); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, s.Trivia); err != nil {
		return err
	}
	if err := v.VisitToken(s.Keyword); err != nil {
		return fmt.Errorf("keyword: %w", err)
	}
	if err := v.VisitIdentifier(s.Name); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if s.Extends != nil {
		if err := v.VisitToken(s.Extends); err != nil {
			return fmt.Errorf("extends: %w", err)
		}
	}
	if s.Parent != nil {
		if err := v.VisitIdentifier(s.Parent); err != nil {
			return fmt.Errorf("parent: %w", err)
		}
	}
	flags := append(s.Hidden, s.Conditional...)
	slices.SortFunc(flags, func(a, b *Token) int { return a.Location.ByteOffset - b.Location.ByteOffset })
	for _, t := range flags {
		if err := v.VisitToken(t); err != nil {
			return fmt.Errorf("flag: %w", err)
		}
	}
	if err := v.VisitDocComment(s.Comment); err != nil {
		return fmt.Errorf("doc comment: %w", err)
	}
	for _, s := range s.Statements {
		if err := VisitScriptStatement(v, s); err != nil {
			return fmt.Errorf("statement: %w", err)
		}
	}
	if err := visitTrailingComments(v, s.Trivia); err != nil {
		return err
	}
	return nil
}

// VisitState visits the [State] node then all children nodes and returns an
// error if any call returns an error.
func (v PreorderVisitor) VisitState(s *State) error {
	if err := v.Delegate.VisitState(s); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, s.Trivia); err != nil {
		return err
	}
	if s.Auto != nil {
		if err := v.VisitToken(s.Auto); err != nil {
			return fmt.Errorf("auto: %w", err)
		}
	}
	if err := v.VisitToken(s.Keyword); err != nil {
		return fmt.Errorf("keyword: %w", err)
	}
	if err := v.VisitIdentifier(s.Name); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	for _, s := range s.Invokables {
		if err := VisitInvokable(v, s); err != nil {
			return fmt.Errorf("invokable: %w", err)
		}
	}
	if err := v.VisitToken(s.EndKeyword); err != nil {
		return fmt.Errorf("end keyword: %w", err)
	}
	if err := visitTrailingComments(v, s.Trivia); err != nil {
		return err
	}
	return nil
}

// VisitToken visits the [Token] node and returns an error if the delegate call
// returns an error.
func (v PreorderVisitor) VisitToken(t *Token) error {
	if err := v.Delegate.VisitToken(t); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitTypeLiteral visits the [TypeLiteral] node then all children nodes and
// returns an error if any call returns an error.
func (v PreorderVisitor) VisitTypeLiteral(t *TypeLiteral) error {
	if err := v.Delegate.VisitTypeLiteral(t); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, t.Trivia); err != nil {
		return err
	}
	if err := v.VisitToken(t.Text); err != nil {
		return fmt.Errorf("text: %w", err)
	}
	if t.Open != nil {
		if err := v.VisitToken(t.Open); err != nil {
			return fmt.Errorf("open: %w", err)
		}
	}
	if t.Close != nil {
		if err := v.VisitToken(t.Close); err != nil {
			return fmt.Errorf("close: %w", err)
		}
	}
	return nil
}

// VisitUnary visits the [Unary] node then all children nodes and returns an
// error if any call returns an error.
func (v PreorderVisitor) VisitUnary(u *Unary) error {
	if err := v.Delegate.VisitUnary(u); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, u.Trivia); err != nil {
		return err
	}
	if err := v.VisitToken(u.Operator); err != nil {
		return fmt.Errorf("operator: %w", err)
	}
	if err := VisitExpression(v, u.Operand); err != nil {
		return fmt.Errorf("operand: %w", err)
	}
	return nil
}

// VisitScriptVariable visits the [ScriptVariable] node then all children nodes
// and returns an error if any call returns an error.
func (v PreorderVisitor) VisitScriptVariable(s *ScriptVariable) error {
	if err := v.Delegate.VisitScriptVariable(s); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, s.Trivia); err != nil {
		return err
	}
	if err := v.VisitTypeLiteral(s.Type); err != nil {
		return fmt.Errorf("type: %w", err)
	}
	if err := v.VisitIdentifier(s.Name); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if s.Operator != nil {
		if err := v.VisitToken(s.Operator); err != nil {
			return fmt.Errorf("operator: %w", err)
		}
	}
	if s.Value != nil {
		if err := VisitLiteral(v, s.Value); err != nil {
			return fmt.Errorf("value: %w", err)
		}
	}
	for _, c := range s.Conditional {
		if err := v.VisitToken(c); err != nil {
			return fmt.Errorf("conditional: %w", err)
		}
	}
	return nil
}

// VisitFunctionVariable visits the [FunctionVariable] node then all children
// nodes and returns an error if any call returns an error.
func (v PreorderVisitor) VisitFunctionVariable(f *FunctionVariable) error {
	if err := v.Delegate.VisitFunctionVariable(f); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, f.Trivia); err != nil {
		return err
	}
	if err := v.VisitTypeLiteral(f.Type); err != nil {
		return fmt.Errorf("type: %w", err)
	}
	if err := v.VisitIdentifier(f.Name); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if f.Operator != nil {
		if err := v.VisitToken(f.Operator); err != nil {
			return fmt.Errorf("operator: %w", err)
		}
	}
	if f.Value != nil {
		if err := VisitExpression(v, f.Value); err != nil {
			return fmt.Errorf("value: %w", err)
		}
	}
	return nil
}

// VisitWhile visits the [While] node then all children nodes and returns an
// error if any call returns an error.
func (v PreorderVisitor) VisitWhile(w *While) error {
	if err := v.Delegate.VisitWhile(w); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitLeadingComments(v, w.Trivia); err != nil {
		return err
	}
	if err := v.VisitToken(w.Keyword); err != nil {
		return fmt.Errorf("keyword: %w", err)
	}
	if err := VisitExpression(v, w.Condition); err != nil {
		return fmt.Errorf("condition: %w", err)
	}
	for _, s := range w.Statements {
		if err := VisitFunctionStatement(v, s); err != nil {
			return fmt.Errorf("statement: %w", err)
		}
	}
	if err := v.VisitToken(w.EndKeyword); err != nil {
		return fmt.Errorf("end keyword: %w", err)
	}
	if err := visitTrailingComments(v, w.Trivia); err != nil {
		return err
	}
	return nil
}

// VisitErrorScriptStatement visits the [ErrorScriptStatement] node and returns
// an error if the delegate call returns an error.
func (v PreorderVisitor) VisitErrorScriptStatement(e *ErrorScriptStatement) error {
	if err := v.Delegate.VisitErrorScriptStatement(e); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitErrorFunctionStatement visits the [ErrorFunctionStatement] node and
// returns an error if the delegate call returns an error.
func (v PreorderVisitor) VisitErrorFunctionStatement(e *ErrorFunctionStatement) error {
	if err := v.Delegate.VisitErrorFunctionStatement(e); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}
