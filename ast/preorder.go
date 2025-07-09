package ast

import (
	"fmt"
)

// PreorderVisitor is a [Visitor] that calls visit methods on
// a delegate visitor by performing a preorder traversal of nodes.
type PreorderVisitor struct {
	Delegate Visitor
}

// VisitAccess visits the [Access] node then all children nodes and returns an
// error if any call returns an error.
func (v *PreorderVisitor) VisitAccess(a *Access) error {
	if err := v.Delegate.VisitAccess(a); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitPrefixComments(v, a); err != nil {
		return err
	}
	if err := a.Value.Accept(v); err != nil {
		return fmt.Errorf("value: %w", err)
	}
	if err := a.Name.Accept(v); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	return visitSuffixComments(v, a)
}

// VisitArgument visits the [Argument] node then all children nodes and
// returns an error if any call returns an error.
func (v *PreorderVisitor) VisitArgument(a *Argument) error {
	if err := v.Delegate.VisitArgument(a); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitPrefixComments(v, a); err != nil {
		return err
	}
	if a.Name != nil {
		if err := a.Name.Accept(v); err != nil {
			return fmt.Errorf("name: %w", err)
		}
	}
	if err := a.Value.Accept(v); err != nil {
		return fmt.Errorf("value: %w", err)
	}
	return visitSuffixComments(v, a)
}

// VisitArrayCreation visits the [ArrayCreation] node then all children nodes
// and returns an error if any call returns an error.
func (v *PreorderVisitor) VisitArrayCreation(a *ArrayCreation) error {
	if err := v.Delegate.VisitArrayCreation(a); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitPrefixComments(v, a); err != nil {
		return err
	}
	if err := a.Type.Accept(v); err != nil {
		return fmt.Errorf("type: %w", err)
	}
	if err := a.Size.Accept(v); err != nil {
		return fmt.Errorf("size: %w", err)
	}
	return visitSuffixComments(v, a)
}

// VisitAssignment visits the [Assignment] node then all children nodes and
// returns an error if any call returns an error.
func (v *PreorderVisitor) VisitAssignment(a *Assignment) error {
	if err := v.Delegate.VisitAssignment(a); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitPrefixComments(v, a); err != nil {
		return err
	}
	if err := a.Assignee.Accept(v); err != nil {
		return fmt.Errorf("assignee: %w", err)
	}
	if err := a.Value.Accept(v); err != nil {
		return fmt.Errorf("value: %w", err)
	}
	return visitSuffixComments(v, a)
}

// VisitBinary visits the [Binary] node then all children nodes and
// returns an error if any call returns an error.
func (v *PreorderVisitor) VisitBinary(b *Binary) error {
	if err := v.Delegate.VisitBinary(b); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitPrefixComments(v, b); err != nil {
		return err
	}
	if err := b.LeftOperand.Accept(v); err != nil {
		return fmt.Errorf("left operand: %w", err)
	}
	if err := b.RightOperand.Accept(v); err != nil {
		return fmt.Errorf("right operand: %w", err)
	}
	return visitSuffixComments(v, b)
}

// VisitCall visits the [Call] node then all children nodes and returns an
// error if any call returns an error.
func (v *PreorderVisitor) VisitCall(c *Call) error {
	if err := v.Delegate.VisitCall(c); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitPrefixComments(v, c); err != nil {
		return err
	}
	if err := c.Function.Accept(v); err != nil {
		return fmt.Errorf("receiver: %w", err)
	}
	for _, a := range c.Arguments {
		if err := a.Accept(v); err != nil {
			return fmt.Errorf("argument: %w", err)
		}
	}
	return visitSuffixComments(v, c)
}

// VisitCast visits the [Cast] node then all children nodes and returns an error
// if any call returns an error.
func (v *PreorderVisitor) VisitCast(c *Cast) error {
	if err := v.Delegate.VisitCast(c); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitPrefixComments(v, c); err != nil {
		return err
	}
	if err := c.Value.Accept(v); err != nil {
		return fmt.Errorf("value: %w", err)
	}
	if err := c.Type.Accept(v); err != nil {
		return fmt.Errorf("type: %w", err)
	}
	return visitSuffixComments(v, c)
}

// VisitDocumentation visits the [Documentation] node then all children nodes
// and returns an error if any call returns an error.
func (v *PreorderVisitor) VisitDocumentation(c *Documentation) error {
	if err := v.Delegate.VisitDocumentation(c); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitBlockComment visits the [BlockComment] node then all children nodes and
// returns an error if any call returns an error.
func (v *PreorderVisitor) VisitBlockComment(c *BlockComment) error {
	if err := v.Delegate.VisitBlockComment(c); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitCommentStatement visits the [CommentStatement] node then all children nodes and
// returns an error if any call returns an error.
func (v *PreorderVisitor) VisitCommentStatement(c *CommentStatement) error {
	if err := v.Delegate.VisitCommentStatement(c); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	for _, e := range c.Elements {
		if err := e.Accept(v); err != nil {
			return fmt.Errorf("element: %w", err)
		}
	}
	return nil
}

// VisitLineComment visits the [LineComment] node then all children nodes and
// returns an error if any call returns an error.
func (v *PreorderVisitor) VisitLineComment(c *LineComment) error {
	if err := v.Delegate.VisitLineComment(c); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}

// VisitEvent visits the [Event] node then all children nodes and returns an
// error if any call returns an error.
func (v *PreorderVisitor) VisitEvent(e *Event) error {
	if err := v.Delegate.VisitEvent(e); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitPrefixComments(v, e); err != nil {
		return err
	}
	if err := e.Name.Accept(v); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	for _, p := range e.ParameterList {
		if err := p.Accept(v); err != nil {
			return fmt.Errorf("parameter: %w", err)
		}
	}
	if err := e.Documentation.Accept(v); err != nil {
		return fmt.Errorf("doc comment: %w", err)
	}
	for _, s := range e.Statements {
		if err := s.Accept(v); err != nil {
			return fmt.Errorf("statement: %w", err)
		}
	}
	return visitSuffixComments(v, e)
}

// VisitFunction visits the [Function] node then all children nodes and returns
// an error if any call returns an error.
func (v *PreorderVisitor) VisitFunction(f *Function) error {
	if err := v.Delegate.VisitFunction(f); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitPrefixComments(v, f); err != nil {
		return err
	}
	if f.ReturnType != nil {
		if err := f.ReturnType.Accept(v); err != nil {
			return fmt.Errorf("return type: %w", err)
		}
	}
	if err := f.Name.Accept(v); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	for _, p := range f.ParameterList {
		if err := p.Accept(v); err != nil {
			return fmt.Errorf("parameter: %w", err)
		}
	}
	if err := f.Documentation.Accept(v); err != nil {
		return fmt.Errorf("doc comment: %w", err)
	}
	for _, s := range f.Statements {
		if err := s.Accept(v); err != nil {
			return fmt.Errorf("statement: %w", err)
		}
	}
	return visitSuffixComments(v, f)
}

// VisitIdentifier visits the [Identifier] node then all children nodes and
// returns an error if any call returns an error.
func (v *PreorderVisitor) VisitIdentifier(i *Identifier) error {
	if err := v.Delegate.VisitIdentifier(i); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return visitComments(v, i)
}

// VisitIf visits the [If] node then all children nodes and returns an error if
// any call returns an error.
func (v *PreorderVisitor) VisitIf(i *If) error {
	if err := v.Delegate.VisitIf(i); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitPrefixComments(v, i); err != nil {
		return err
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
	return visitSuffixComments(v, i)
}

// VisitElseIf visits the [ElseIf] node then all children nodes and returns an
// error if any call returns an error.
func (v *PreorderVisitor) VisitElseIf(e *ElseIf) error {
	if err := v.Delegate.VisitElseIf(e); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitPrefixComments(v, e); err != nil {
		return err
	}
	if err := e.Condition.Accept(v); err != nil {
		return fmt.Errorf("condition: %w", err)
	}
	for _, s := range e.Statements {
		if err := s.Accept(v); err != nil {
			return fmt.Errorf("statement: %w", err)
		}
	}
	return visitSuffixComments(v, e)
}

// VisitElse visits the [Else] node then all children nodes and returns an error
// if any call returns an error.
func (v *PreorderVisitor) VisitElse(e *Else) error {
	if err := v.Delegate.VisitElse(e); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitPrefixComments(v, e); err != nil {
		return err
	}
	for _, s := range e.Statements {
		if err := s.Accept(v); err != nil {
			return fmt.Errorf("statement: %w", err)
		}
	}
	return visitSuffixComments(v, e)
}

// VisitExpressionStatement visits the [ExpressionStatement] node then all
// children nodes and returns an error if any call returns an error.
func (v *PreorderVisitor) VisitExpressionStatement(s *ExpressionStatement) error {
	if err := v.Delegate.VisitExpressionStatement(s); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitPrefixComments(v, s); err != nil {
		return err
	}
	if err := s.Expression.Accept(v); err != nil {
		return fmt.Errorf("expression: %w", err)
	}
	return visitSuffixComments(v, s)
}

// VisitImport visits the [Import] node then all children nodes and returns an
// error if any call returns an error.
func (v *PreorderVisitor) VisitImport(i *Import) error {
	if err := v.Delegate.VisitImport(i); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitPrefixComments(v, i); err != nil {
		return err
	}
	if err := i.Name.Accept(v); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	return visitSuffixComments(v, i)
}

// VisitIndex visits the [Index] node then all children nodes and returns an
// error if any call returns an error.
func (v *PreorderVisitor) VisitIndex(i *Index) error {
	if err := v.Delegate.VisitIndex(i); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitPrefixComments(v, i); err != nil {
		return err
	}
	if err := i.Value.Accept(v); err != nil {
		return fmt.Errorf("value: %w", err)
	}
	if err := i.Index.Accept(v); err != nil {
		return fmt.Errorf("index: %w", err)
	}
	return visitSuffixComments(v, i)
}

// VisitBoolLiteral visits the [BoolLiteral] node then all children nodes and
// returns an error if any call returns an error.
func (v *PreorderVisitor) VisitBoolLiteral(l *BoolLiteral) error {
	if err := v.Delegate.VisitBoolLiteral(l); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return visitComments(v, l)
}

// VisitIntLiteral visits the [IntLiteral] node then all children nodes and
// returns an error if any call returns an error.
func (v *PreorderVisitor) VisitIntLiteral(l *IntLiteral) error {
	if err := v.Delegate.VisitIntLiteral(l); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return visitComments(v, l)
}

// VisitFloatLiteral visits the [FloatLiteral] node then all children nodes and
// returns an error if any call returns an error.
func (v *PreorderVisitor) VisitFloatLiteral(l *FloatLiteral) error {
	if err := v.Delegate.VisitFloatLiteral(l); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return visitComments(v, l)
}

// VisitStringLiteral visits the [StringLiteral] node then all children nodes
// and returns an error if any call returns an error.
func (v *PreorderVisitor) VisitStringLiteral(l *StringLiteral) error {
	if err := v.Delegate.VisitStringLiteral(l); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return visitComments(v, l)
}

// VisitNoneLiteral visits the [NoneLiteral] node then all children nodes and
// returns an error if any call returns an error.
func (v *PreorderVisitor) VisitNoneLiteral(l *NoneLiteral) error {
	if err := v.Delegate.VisitNoneLiteral(l); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return visitComments(v, l)
}

// VisitParameter visits the [Parameter] node then all children nodes and
// returns an error if any call returns an error.
func (v *PreorderVisitor) VisitParameter(p *Parameter) error {
	if err := v.Delegate.VisitParameter(p); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitPrefixComments(v, p); err != nil {
		return err
	}
	if err := p.Type.Accept(v); err != nil {
		return fmt.Errorf("type: %w", err)
	}
	if err := p.Name.Accept(v); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if p.DefaultValue != nil {
		if err := p.DefaultValue.Accept(v); err != nil {
			return fmt.Errorf("value: %w", err)
		}
	}
	return visitSuffixComments(v, p)
}

// VisitParenthetical visits the [Parenthetical] node then all children nodes
// and returns an error if any call returns an error.
func (v *PreorderVisitor) VisitParenthetical(p *Parenthetical) error {
	if err := v.Delegate.VisitParenthetical(p); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitPrefixComments(v, p); err != nil {
		return err
	}
	if err := p.Value.Accept(v); err != nil {
		return fmt.Errorf("value: %w", err)
	}
	return visitSuffixComments(v, p)
}

// VisitProperty visits the [Property] node then all children nodes and returns
// an error if any call returns an error.
func (v *PreorderVisitor) VisitProperty(p *Property) error {
	if err := v.Delegate.VisitProperty(p); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitPrefixComments(v, p); err != nil {
		return err
	}
	if err := p.Type.Accept(v); err != nil {
		return fmt.Errorf("type: %w", err)
	}
	if err := p.Name.Accept(v); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if p.Value != nil {
		if err := p.Value.Accept(v); err != nil {
			return fmt.Errorf("value: %w", err)
		}
	}
	if err := p.Documentation.Accept(v); err != nil {
		return fmt.Errorf("doc comment: %w", err)
	}
	first := p.Get
	second := p.Set
	if p.Get != nil && p.Set != nil && second.Location().Compare(first.Location()) < 0 {
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
	return visitSuffixComments(v, p)
}

// VisitReturn visits the [Return] node then all children nodes and returns an
// error if any call returns an error.
func (v *PreorderVisitor) VisitReturn(r *Return) error {
	if err := v.Delegate.VisitReturn(r); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitPrefixComments(v, r); err != nil {
		return err
	}
	if r.Value != nil {
		if err := r.Value.Accept(v); err != nil {
			return fmt.Errorf("value: %w", err)
		}
	}
	return visitSuffixComments(v, r)
}

// VisitScript visits the [Script] node then all children nodes and returns an
// error if any call returns an error.
func (v *PreorderVisitor) VisitScript(s *Script) error {
	if err := v.Delegate.VisitScript(s); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	for _, c := range s.HeaderComments {
		if err := c.Accept(v); err != nil {
			return fmt.Errorf("header comment: %w", err)
		}
	}
	if err := s.Name.Accept(v); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if s.Parent != nil {
		if err := s.Parent.Accept(v); err != nil {
			return fmt.Errorf("parent: %w", err)
		}
	}
	if err := s.Documentation.Accept(v); err != nil {
		return fmt.Errorf("doc comment: %w", err)
	}
	for _, s := range s.Statements {
		if err := s.Accept(v); err != nil {
			return fmt.Errorf("statement: %w", err)
		}
	}
	return nil
}

// VisitState visits the [State] node then all children nodes and returns an
// error if any call returns an error.
func (v *PreorderVisitor) VisitState(s *State) error {
	if err := v.Delegate.VisitState(s); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitPrefixComments(v, s); err != nil {
		return err
	}
	if err := s.Name.Accept(v); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	for _, s := range s.Invokables {
		if err := s.Accept(v); err != nil {
			return fmt.Errorf("invokable: %w", err)
		}
	}
	return visitSuffixComments(v, s)
}

// VisitTypeLiteral visits the [TypeLiteral] node then all children nodes and
// returns an error if any call returns an error.
func (v *PreorderVisitor) VisitTypeLiteral(t *TypeLiteral) error {
	if err := v.Delegate.VisitTypeLiteral(t); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return visitComments(v, t)
}

// VisitUnary visits the [Unary] node then all children nodes and returns an
// error if any call returns an error.
func (v *PreorderVisitor) VisitUnary(u *Unary) error {
	if err := v.Delegate.VisitUnary(u); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitComments(v, u); err != nil {
		return err
	}
	if err := u.Operand.Accept(v); err != nil {
		return fmt.Errorf("operand: %w", err)
	}
	return nil
}

// VisitVariable visits the [Variable] node then all children nodes
// and returns an error if any call returns an error.
func (v *PreorderVisitor) VisitVariable(r *Variable) error {
	if err := v.Delegate.VisitVariable(r); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitPrefixComments(v, r); err != nil {
		return err
	}
	if err := r.Type.Accept(v); err != nil {
		return fmt.Errorf("type: %w", err)
	}
	if err := r.Name.Accept(v); err != nil {
		return fmt.Errorf("name: %w", err)
	}
	if r.Value != nil {
		if err := r.Value.Accept(v); err != nil {
			return fmt.Errorf("value: %w", err)
		}
	}
	return visitSuffixComments(v, r)
}

// VisitWhile visits the [While] node then all children nodes and returns an
// error if any call returns an error.
func (v *PreorderVisitor) VisitWhile(w *While) error {
	if err := v.Delegate.VisitWhile(w); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	if err := visitPrefixComments(v, w); err != nil {
		return err
	}
	if err := w.Condition.Accept(v); err != nil {
		return fmt.Errorf("condition: %w", err)
	}
	for _, s := range w.Statements {
		if err := s.Accept(v); err != nil {
			return fmt.Errorf("statement: %w", err)
		}
	}
	return visitSuffixComments(v, w)
}

// VisitErrorStatement visits the [ErrorStatement] node and returns
// an error if the delegate call returns an error.
func (v *PreorderVisitor) VisitErrorStatement(e *ErrorStatement) error {
	if err := v.Delegate.VisitErrorStatement(e); err != nil {
		return fmt.Errorf("delegate: %w", err)
	}
	return nil
}
