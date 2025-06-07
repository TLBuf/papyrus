package ast

import (
	"fmt"
	"slices"
)

// PostorderVisitor is a [Visitor] that calls visit methods on
// a delegate visitor by performing a postorder traversal of nodes.
type PostorderVisitor struct {
	Delegate Visitor
}

// VisitAccess visits the [Access] node then all children nodes and returns an
// error if any call returns an error.
func (v PostorderVisitor) VisitAccess(a *Access) error {
	for _, c := range a.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range a.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := a.Accept(v); err != nil {
		return fmt.Errorf("value: %w", err)
	}
	if err := a.Operator.Accept(v); err != nil {
		return fmt.Errorf("operator: %w", err)
	}
	if err := a.Name.Accept(v); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	for _, c := range a.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range a.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitAccess(a); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// PostorderVisitor visits the [Argument] node then all children nodes and
// returns an error if any call returns an error.
func (v PostorderVisitor) VisitArgument(a *Argument) error {
	for _, c := range a.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range a.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := a.Name.Accept(v); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if err := a.Operator.Accept(v); err != nil {
		return fmt.Errorf("operator: %w", err)
	}
	if err := a.Accept(v); err != nil {
		return fmt.Errorf("value: %w", err)
	}
	for _, c := range a.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range a.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitArgument(a); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitArrayCreation visits the [ArrayCreation] node then all children nodes
// and returns an error if any call returns an error.
func (v PostorderVisitor) VisitArrayCreation(a *ArrayCreation) error {
	for _, c := range a.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range a.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := a.NewOperator.Accept(v); err != nil {
		return fmt.Errorf("new operator: %w", err)
	}
	if err := a.Type.Accept(v); err != nil {
		return fmt.Errorf("type: %w", err)
	}
	if err := a.Open.Accept(v); err != nil {
		return fmt.Errorf("open: %w", err)
	}
	if err := a.Size.Accept(v); err != nil {
		return fmt.Errorf("size: %w", err)
	}
	if err := a.Close.Accept(v); err != nil {
		return fmt.Errorf("close: %w", err)
	}
	for _, c := range a.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range a.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitArrayCreation(a); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitAssignment visits the [Assignment] node then all children nodes and
// returns an error if any call returns an error.
func (v PostorderVisitor) VisitAssignment(a *Assignment) error {
	for _, c := range a.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range a.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := a.Assignee.Accept(v); err != nil {
		return fmt.Errorf("assignee: %w", err)
	}
	if err := a.Operator.Accept(v); err != nil {
		return fmt.Errorf("operator: %w", err)
	}
	if err := a.Accept(v); err != nil {
		return fmt.Errorf("value: %w", err)
	}
	for _, c := range a.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range a.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitAssignment(a); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitBinary visits the [Binary] node then all children nodes and
// returns an error if any call returns an error.
func (v PostorderVisitor) VisitBinary(b *Binary) error {
	for _, c := range b.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range b.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := b.LeftOperand.Accept(v); err != nil {
		return fmt.Errorf("left operand: %w", err)
	}
	if err := b.Operator.Accept(v); err != nil {
		return fmt.Errorf("operator: %w", err)
	}
	if err := b.RightOperand.Accept(v); err != nil {
		return fmt.Errorf("right operand: %w", err)
	}
	for _, c := range b.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range b.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitBinary(b); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitCall visits the [Call] node then all children nodes and returns an
// error if any call returns an error.
func (v PostorderVisitor) VisitCall(c *Call) error {
	for _, c := range c.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range c.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := c.Reciever.Accept(v); err != nil {
		return fmt.Errorf("reciever: %w", err)
	}
	if err := c.Open.Accept(v); err != nil {
		return fmt.Errorf("open: %w", err)
	}
	for _, a := range c.Arguments {
		if err := a.Accept(v); err != nil {
			return fmt.Errorf("argument: %w", err)
		}
	}
	if err := c.Close.Accept(v); err != nil {
		return fmt.Errorf("close: %w", err)
	}
	for _, c := range c.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range c.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitCall(c); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitCast visits the [Cast] node then all children nodes and returns an error
// if any call returns an error.
func (v PostorderVisitor) VisitCast(c *Cast) error {
	for _, c := range c.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range c.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := c.Accept(v); err != nil {
		return fmt.Errorf("value: %w", err)
	}
	if err := c.Operator.Accept(v); err != nil {
		return fmt.Errorf("operator: %w", err)
	}
	if err := c.Type.Accept(v); err != nil {
		return fmt.Errorf("type: %w", err)
	}
	for _, c := range c.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range c.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitCast(c); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitDocComment visits the [DocComment] node then all children nodes and
// returns an error if any call returns an error.
func (v PostorderVisitor) VisitDocComment(c *DocComment) error {
	if err := c.Open.Accept(v); err != nil {
		return fmt.Errorf("open: %w", err)
	}
	if err := c.Text.Accept(v); err != nil {
		return fmt.Errorf("text: %w", err)
	}
	if err := c.Close.Accept(v); err != nil {
		return fmt.Errorf("type: %w", err)
	}
	if err := v.Delegate.VisitDocComment(c); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitBlockComment visits the [BlockComment] node then all children nodes and
// returns an error if any call returns an error.
func (v PostorderVisitor) VisitBlockComment(c *BlockComment) error {
	if err := c.Open.Accept(v); err != nil {
		return fmt.Errorf("open: %w", err)
	}
	if err := c.Text.Accept(v); err != nil {
		return fmt.Errorf("text: %w", err)
	}
	if err := c.Close.Accept(v); err != nil {
		return fmt.Errorf("type: %w", err)
	}
	if err := v.Delegate.VisitBlockComment(c); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitLineComment visits the [LineComment] node then all children nodes and
// returns an error if any call returns an error.
func (v PostorderVisitor) VisitLineComment(c *LineComment) error {
	if err := c.Open.Accept(v); err != nil {
		return fmt.Errorf("open: %w", err)
	}
	if err := c.Text.Accept(v); err != nil {
		return fmt.Errorf("text: %w", err)
	}
	if err := v.Delegate.VisitLineComment(c); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitEvent visits the [Event] node then all children nodes and returns an
// error if any call returns an error.
func (v PostorderVisitor) VisitEvent(e *Event) error {
	for _, c := range e.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range e.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := e.Keyword.Accept(v); err != nil {
		return fmt.Errorf("keyword: %w", err)
	}
	if err := e.Name.Accept(v); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if err := e.Open.Accept(v); err != nil {
		return fmt.Errorf("open: %w", err)
	}
	for _, p := range e.Parameters {
		if err := p.Accept(v); err != nil {
			return fmt.Errorf("parameter: %w", err)
		}
	}
	if err := e.Close.Accept(v); err != nil {
		return fmt.Errorf("close: %w", err)
	}
	for _, t := range e.Native {
		if err := t.Accept(v); err != nil {
			return fmt.Errorf("native: %w", err)
		}
	}
	if err := e.Comment.Accept(v); err != nil {
		return fmt.Errorf("doc comment: %w", err)
	}
	for _, s := range e.Statements {
		if err := s.Accept(v); err != nil {
			return fmt.Errorf("statement: %w", err)
		}
	}
	if err := e.EndKeyword.Accept(v); err != nil {
		return fmt.Errorf("end keyword: %w", err)
	}
	for _, c := range e.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range e.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitEvent(e); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitFunction visits the [Function] node then all children nodes and returns
// an error if any call returns an error.
func (v PostorderVisitor) VisitFunction(f *Function) error {
	for _, c := range f.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range f.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if f.ReturnType != nil {
		if err := f.ReturnType.Accept(v); err != nil {
			return fmt.Errorf("return type: %w", err)
		}
	}
	if err := f.Keyword.Accept(v); err != nil {
		return fmt.Errorf("keyword: %w", err)
	}
	if err := f.Name.Accept(v); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if err := f.Open.Accept(v); err != nil {
		return fmt.Errorf("open: %w", err)
	}
	for _, p := range f.Parameters {
		if err := p.Accept(v); err != nil {
			return fmt.Errorf("parameter: %w", err)
		}
	}
	if err := f.Close.Accept(v); err != nil {
		return fmt.Errorf("close: %w", err)
	}
	flags := append(f.Global, f.Native...)
	slices.SortFunc(flags, func(a, b *Token) int { return a.Location.ByteOffset - b.Location.ByteOffset })
	for _, t := range flags {
		if err := t.Accept(v); err != nil {
			return fmt.Errorf("flag: %w", err)
		}
	}
	if err := f.Comment.Accept(v); err != nil {
		return fmt.Errorf("doc comment: %w", err)
	}
	for _, s := range f.Statements {
		if err := s.Accept(v); err != nil {
			return fmt.Errorf("statement: %w", err)
		}
	}
	if err := f.EndKeyword.Accept(v); err != nil {
		return fmt.Errorf("end keyword: %w", err)
	}
	for _, c := range f.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range f.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitFunction(f); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitIdentifier visits the [Identifier] node then all children nodes and
// returns an error if any call returns an error.
func (v PostorderVisitor) VisitIdentifier(i *Identifier) error {
	for _, c := range i.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range i.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := i.Text.Accept(v); err != nil {
		return fmt.Errorf("text: %w", err)
	}
	for _, c := range i.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range i.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitIdentifier(i); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitIf visits the [If] node then all children nodes and returns an error if
// any call returns an error.
func (v PostorderVisitor) VisitIf(i *If) error {
	for _, c := range i.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range i.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := i.Keyword.Accept(v); err != nil {
		return fmt.Errorf("keyword: %w", err)
	}
	if err := i.Condition.Accept(v); err != nil {
		return fmt.Errorf("condition: %w", err)
	}
	for _, s := range i.Statements {
		if err := s.Accept(v); err != nil {
			return fmt.Errorf("statement: %w", err)
		}
	}
	for _, e := range i.ElseIfs {
		if err := e.Accept(v); err != nil {
			return fmt.Errorf("else if: %w", err)
		}
	}
	if i.Else != nil {
		if err := i.Else.Accept(v); err != nil {
			return fmt.Errorf("else: %w", err)
		}
	}
	if err := i.EndKeyword.Accept(v); err != nil {
		return fmt.Errorf("end keyword: %w", err)
	}
	for _, c := range i.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range i.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitIf(i); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitElseIf visits the [ElseIf] node then all children nodes and returns an
// error if any call returns an error.
func (v PostorderVisitor) VisitElseIf(e *ElseIf) error {
	for _, c := range e.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range e.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := e.Keyword.Accept(v); err != nil {
		return fmt.Errorf("keyword: %w", err)
	}
	if err := e.Condition.Accept(v); err != nil {
		return fmt.Errorf("condition: %w", err)
	}
	for _, s := range e.Statements {
		if err := s.Accept(v); err != nil {
			return fmt.Errorf("statement: %w", err)
		}
	}
	for _, c := range e.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range e.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitElseIf(e); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitElse visits the [Else] node then all children nodes and returns an error
// if any call returns an error.
func (v PostorderVisitor) VisitElse(e *Else) error {
	for _, c := range e.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range e.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := e.Keyword.Accept(v); err != nil {
		return fmt.Errorf("keyword: %w", err)
	}
	for _, s := range e.Statements {
		if err := s.Accept(v); err != nil {
			return fmt.Errorf("statement: %w", err)
		}
	}
	for _, c := range e.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range e.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitElse(e); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitImport visits the [Import] node then all children nodes and returns an
// error if any call returns an error.
func (v PostorderVisitor) VisitImport(i *Import) error {
	for _, c := range i.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range i.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := i.Keyword.Accept(v); err != nil {
		return fmt.Errorf("keyword: %w", err)
	}
	if err := i.Name.Accept(v); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	for _, c := range i.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range i.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitImport(i); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitIndex visits the [Index] node then all children nodes and returns an
// error if any call returns an error.
func (v PostorderVisitor) VisitIndex(i *Index) error {
	for _, c := range i.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range i.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := i.Accept(v); err != nil {
		return fmt.Errorf("value: %w", err)
	}
	if err := i.Open.Accept(v); err != nil {
		return fmt.Errorf("open: %w", err)
	}
	if err := i.Index.Accept(v); err != nil {
		return fmt.Errorf("index: %w", err)
	}
	if err := i.Close.Accept(v); err != nil {
		return fmt.Errorf("close: %w", err)
	}
	for _, c := range i.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range i.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitIndex(i); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitBoolLiteral visits the [BoolLiteral] node then all children nodes and
// returns an error if any call returns an error.
func (v PostorderVisitor) VisitBoolLiteral(l *BoolLiteral) error {
	for _, c := range l.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range l.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := l.Text.Accept(v); err != nil {
		return fmt.Errorf("text: %w", err)
	}
	for _, c := range l.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range l.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitBoolLiteral(l); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitIntLiteral visits the [IntLiteral] node then all children nodes and
// returns an error if any call returns an error.
func (v PostorderVisitor) VisitIntLiteral(l *IntLiteral) error {
	for _, c := range l.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range l.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := l.Text.Accept(v); err != nil {
		return fmt.Errorf("text: %w", err)
	}
	for _, c := range l.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range l.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitIntLiteral(l); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitFloatLiteral visits the [FloatLiteral] node then all children nodes and
// returns an error if any call returns an error.
func (v PostorderVisitor) VisitFloatLiteral(l *FloatLiteral) error {
	for _, c := range l.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range l.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := l.Text.Accept(v); err != nil {
		return fmt.Errorf("text: %w", err)
	}
	for _, c := range l.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range l.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitFloatLiteral(l); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitStringLiteral visits the [StringLiteral] node then all children nodes
// and returns an error if any call returns an error.
func (v PostorderVisitor) VisitStringLiteral(l *StringLiteral) error {
	for _, c := range l.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range l.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := l.Text.Accept(v); err != nil {
		return fmt.Errorf("text: %w", err)
	}
	for _, c := range l.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range l.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitStringLiteral(l); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitNoneLiteral visits the [NoneLiteral] node then all children nodes and
// returns an error if any call returns an error.
func (v PostorderVisitor) VisitNoneLiteral(l *NoneLiteral) error {
	for _, c := range l.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range l.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := l.Text.Accept(v); err != nil {
		return fmt.Errorf("text: %w", err)
	}
	for _, c := range l.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range l.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitNoneLiteral(l); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitParameter visits the [Parameter] node then all children nodes and
// returns an error if any call returns an error.
func (v PostorderVisitor) VisitParameter(p *Parameter) error {
	for _, c := range p.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range p.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := p.Type.Accept(v); err != nil {
		return fmt.Errorf("type: %w", err)
	}
	if err := p.Name.Accept(v); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if p.Operator != nil {
		if err := p.Operator.Accept(v); err != nil {
			return fmt.Errorf("operator: %w", err)
		}
	}
	if p.Value != nil {
		if err := p.Value.Accept(v); err != nil {
			return fmt.Errorf("value: %w", err)
		}
	}
	for _, c := range p.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range p.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitParameter(p); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitParenthetical visits the [Parenthetical] node then all children nodes
// and returns an error if any call returns an error.
func (v PostorderVisitor) VisitParenthetical(p *Parenthetical) error {
	for _, c := range p.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range p.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := p.Open.Accept(v); err != nil {
		return fmt.Errorf("open: %w", err)
	}
	if err := p.Accept(v); err != nil {
		return fmt.Errorf("value: %w", err)
	}
	if err := p.Close.Accept(v); err != nil {
		return fmt.Errorf("close: %w", err)
	}
	for _, c := range p.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range p.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitParenthetical(p); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitProperty visits the [Property] node then all children nodes and returns
// an error if any call returns an error.
func (v PostorderVisitor) VisitProperty(p *Property) error {
	for _, c := range p.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range p.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := p.Type.Accept(v); err != nil {
		return fmt.Errorf("type: %w", err)
	}
	if err := p.Keyword.Accept(v); err != nil {
		return fmt.Errorf("keyword: %w", err)
	}
	if err := p.Name.Accept(v); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if p.Operator != nil {
		if err := p.Operator.Accept(v); err != nil {
			return fmt.Errorf("operator: %w", err)
		}
	}
	if p.Value != nil {
		if err := p.Accept(v); err != nil {
			return fmt.Errorf("value: %w", err)
		}
	}
	if p.Auto != nil {
		if err := p.Auto.Accept(v); err != nil {
			return fmt.Errorf("auto: %w", err)
		}
	}
	if p.AutoReadOnly != nil {
		if err := p.AutoReadOnly.Accept(v); err != nil {
			return fmt.Errorf("autoreadonly: %w", err)
		}
	}
	flags := append(p.Hidden, p.Conditional...)
	slices.SortFunc(flags, func(a, b *Token) int { return a.Location.ByteOffset - b.Location.ByteOffset })
	for _, t := range flags {
		if err := t.Accept(v); err != nil {
			return fmt.Errorf("flag: %w", err)
		}
	}
	if err := p.Comment.Accept(v); err != nil {
		return fmt.Errorf("doc comment: %w", err)
	}
	first := p.Get
	second := p.Set
	if p.Get != nil && p.Set != nil && second.Location.ByteOffset < first.Location.ByteOffset {
		first = p.Set
		second = p.Get
	}
	if first != nil {
		if err := first.Accept(v); err != nil {
			return fmt.Errorf("function: %w", err)
		}
	}
	if second != nil {
		if err := second.Accept(v); err != nil {
			return fmt.Errorf("function: %w", err)
		}
	}
	if err := p.EndKeyword.Accept(v); err != nil {
		return fmt.Errorf("end keyword: %w", err)
	}
	for _, c := range p.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range p.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitProperty(p); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitReturn visits the [Return] node then all children nodes and returns an
// error if any call returns an error.
func (v PostorderVisitor) VisitReturn(r *Return) error {
	for _, c := range r.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range r.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := r.Keyword.Accept(v); err != nil {
		return fmt.Errorf("keyword: %w", err)
	}
	if r.Value != nil {
		if err := r.Accept(v); err != nil {
			return fmt.Errorf("value: %w", err)
		}
	}
	if err := v.Delegate.VisitReturn(r); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitScript visits the [Script] node then all children nodes and returns an
// error if any call returns an error.
func (v PostorderVisitor) VisitScript(s *Script) error {
	for _, c := range s.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range s.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := s.Keyword.Accept(v); err != nil {
		return fmt.Errorf("keyword: %w", err)
	}
	if err := s.Name.Accept(v); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if s.Extends != nil {
		if err := s.Extends.Accept(v); err != nil {
			return fmt.Errorf("extends: %w", err)
		}
	}
	if s.Parent != nil {
		if err := s.Parent.Accept(v); err != nil {
			return fmt.Errorf("parent: %w", err)
		}
	}
	flags := append(s.Hidden, s.Conditional...)
	slices.SortFunc(flags, func(a, b *Token) int { return a.Location.ByteOffset - b.Location.ByteOffset })
	for _, t := range flags {
		if err := t.Accept(v); err != nil {
			return fmt.Errorf("flag: %w", err)
		}
	}
	if err := s.Comment.Accept(v); err != nil {
		return fmt.Errorf("doc comment: %w", err)
	}
	for _, s := range s.Statements {
		if err := s.Accept(v); err != nil {
			return fmt.Errorf("statement: %w", err)
		}
	}
	for _, c := range s.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range s.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitScript(s); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitState visits the [State] node then all children nodes and returns an
// error if any call returns an error.
func (v PostorderVisitor) VisitState(s *State) error {
	for _, c := range s.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range s.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if s.Auto != nil {
		if err := s.Auto.Accept(v); err != nil {
			return fmt.Errorf("auto: %w", err)
		}
	}
	if err := s.Keyword.Accept(v); err != nil {
		return fmt.Errorf("keyword: %w", err)
	}
	if err := s.Name.Accept(v); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	for _, s := range s.Invokables {
		if err := s.Accept(v); err != nil {
			return fmt.Errorf("invokable: %w", err)
		}
	}
	if err := s.EndKeyword.Accept(v); err != nil {
		return fmt.Errorf("end keyword: %w", err)
	}
	for _, c := range s.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range s.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitState(s); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitToken visits the [Token] node and returns an error if the delegate call
// returns an error.
func (v PostorderVisitor) VisitToken(t *Token) error {
	if err := v.Delegate.VisitToken(t); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitTypeLiteral visits the [TypeLiteral] node then all children nodes and
// returns an error if any call returns an error.
func (v PostorderVisitor) VisitTypeLiteral(t *TypeLiteral) error {
	for _, c := range t.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range t.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := t.Text.Accept(v); err != nil {
		return fmt.Errorf("text: %w", err)
	}
	if t.Open != nil {
		if err := t.Open.Accept(v); err != nil {
			return fmt.Errorf("open: %w", err)
		}
	}
	if t.Close != nil {
		if err := t.Close.Accept(v); err != nil {
			return fmt.Errorf("close: %w", err)
		}
	}
	if err := v.Delegate.VisitTypeLiteral(t); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitUnary visits the [Unary] node then all children nodes and returns an
// error if any call returns an error.
func (v PostorderVisitor) VisitUnary(u *Unary) error {
	for _, c := range u.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range u.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := u.Operator.Accept(v); err != nil {
		return fmt.Errorf("operator: %w", err)
	}
	if err := u.Operand.Accept(v); err != nil {
		return fmt.Errorf("operand: %w", err)
	}
	if err := v.Delegate.VisitUnary(u); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitScriptVariable visits the [ScriptVariable] node then all children nodes
// and returns an error if any call returns an error.
func (v PostorderVisitor) VisitScriptVariable(s *ScriptVariable) error {
	for _, c := range s.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range s.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := s.Type.Accept(v); err != nil {
		return fmt.Errorf("type: %w", err)
	}
	if err := s.Name.Accept(v); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if s.Operator != nil {
		if err := s.Operator.Accept(v); err != nil {
			return fmt.Errorf("operator: %w", err)
		}
	}
	if s.Value != nil {
		if err := s.Value.Accept(v); err != nil {
			return fmt.Errorf("value: %w", err)
		}
	}
	for _, c := range s.Conditional {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("conditional: %w", err)
		}
	}
	if err := v.Delegate.VisitScriptVariable(s); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitFunctionVariable visits the [FunctionVariable] node then all children
// nodes and returns an error if any call returns an error.
func (v PostorderVisitor) VisitFunctionVariable(f *FunctionVariable) error {
	for _, c := range f.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range f.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := f.Type.Accept(v); err != nil {
		return fmt.Errorf("type: %w", err)
	}
	if err := f.Name.Accept(v); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if f.Operator != nil {
		if err := f.Operator.Accept(v); err != nil {
			return fmt.Errorf("operator: %w", err)
		}
	}
	if f.Value != nil {
		if err := f.Accept(v); err != nil {
			return fmt.Errorf("value: %w", err)
		}
	}
	if err := v.Delegate.VisitFunctionVariable(f); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitWhile visits the [While] node then all children nodes and returns an
// error if any call returns an error.
func (v PostorderVisitor) VisitWhile(w *While) error {
	for _, c := range w.LeadingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("leading comment: %w", err)
		}
	}
	for _, c := range w.PrefixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("prefix comment: %w", err)
		}
	}
	if err := w.Keyword.Accept(v); err != nil {
		return fmt.Errorf("keyword: %w", err)
	}
	if err := w.Condition.Accept(v); err != nil {
		return fmt.Errorf("condition: %w", err)
	}
	for _, s := range w.Statements {
		if err := s.Accept(v); err != nil {
			return fmt.Errorf("statement: %w", err)
		}
	}
	if err := w.EndKeyword.Accept(v); err != nil {
		return fmt.Errorf("end keyword: %w", err)
	}
	for _, c := range w.SuffixComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("suffix comment: %w", err)
		}
	}
	for _, c := range w.TrailingComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("trailing comment: %w", err)
		}
	}
	if err := v.Delegate.VisitWhile(w); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitErrorScriptStatement visits the [ErrorScriptStatement] node and returns
// an error if the delegate call returns an error.
func (v PostorderVisitor) VisitErrorScriptStatement(e *ErrorScriptStatement) error {
	if err := v.Delegate.VisitErrorScriptStatement(e); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitErrorFunctionStatement visits the [ErrorFunctionStatement] node and
// returns an error if the delegate call returns an error.
func (v PostorderVisitor) VisitErrorFunctionStatement(e *ErrorFunctionStatement) error {
	if err := v.Delegate.VisitErrorFunctionStatement(e); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}
