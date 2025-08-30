package parser

import (
	"github.com/TLBuf/papyrus/ast"
)

// attachLooseComments attaches comments to nodes in place.
func (p *parser) attachLooseComments(script *ast.Script, comments []ast.Comment) {
	if len(comments) == 0 {
		return
	}

	var visitor nodes
	_ = script.Accept(&ast.PreorderVisitor{Delegate: &visitor})
	preorder := visitor.nodes
	visitor.nodes = make([]commentable, 0, len(visitor.nodes))
	_ = script.Accept(&ast.PostorderVisitor{Delegate: &visitor})
	postorder := visitor.nodes

	prefixCursor := 0
	suffixCursor := 0
	for _, comment := range comments {
		cl := comment.Location()
		switch {
		case comment.Prefix():
			last := postorder[len(postorder)-1]
			ll := last.Location()
			if ll.End() < cl.Start() {
				p.attachSuffixComments(last, comment)
				continue
			}
			node := preorder[prefixCursor]
			for prefixCursor < len(preorder) {
				if node.Location().Start() > cl.Start() || prefixCursor == len(preorder)-1 {
					break
				}
				prefixCursor++
				node = preorder[prefixCursor]
			}
			p.attachPrefixComments(last, comment)
		case comment.Suffix():
			node := postorder[suffixCursor]
			for suffixCursor < len(postorder) {
				if suffixCursor == len(postorder)-1 {
					break
				}
				curr := postorder[suffixCursor+1]
				nl := curr.Location()
				if nl.End() > cl.Start() {
					break
				}
				node = curr
				suffixCursor++
			}
			p.attachSuffixComments(node, comment)
		default:
			p.failWithDetail(intenalInvalidState, comment.Location(), "Unexpected standalone commnet.")
		}
	}
}

type commentable interface {
	ast.Node

	Comments() *ast.Comments
}

func (p *parser) attachPrefixComments(node commentable, comments ...ast.Comment) {
	dst := nodeComments(node)
	if dst == nil {
		p.failWithDetail(intenalInvalidState, node.Location(), "Cannot attach comments to %T", node)
	}
	dst.PrefixComments = append(dst.PrefixComments, comments...)
}

func (p *parser) attachSuffixComments(node commentable, comments ...ast.Comment) {
	dst := nodeComments(node)
	if dst == nil {
		p.failWithDetail(intenalInvalidState, node.Location(), "Cannot attach comments to %T", node)
	}
	dst.SuffixComments = append(dst.SuffixComments, comments...)
}

func nodeComments(node commentable) *ast.Comments {
	var ptr **ast.Comments
	switch node := node.(type) {
	case *ast.Access:
		ptr = &node.NodeComments
	case *ast.Argument:
		ptr = &node.NodeComments
	case *ast.ArrayCreation:
		ptr = &node.NodeComments
	case *ast.Assignment:
		ptr = &node.NodeComments
	case *ast.Binary:
		ptr = &node.NodeComments
	case *ast.BoolLiteral:
		ptr = &node.NodeComments
	case *ast.Call:
		ptr = &node.NodeComments
	case *ast.Cast:
		ptr = &node.NodeComments
	case *ast.Else:
		ptr = &node.NodeComments
	case *ast.ElseIf:
		ptr = &node.NodeComments
	case *ast.Event:
		ptr = &node.NodeComments
	case *ast.ExpressionStatement:
		ptr = &node.NodeComments
	case *ast.FloatLiteral:
		ptr = &node.NodeComments
	case *ast.Function:
		ptr = &node.NodeComments
	case *ast.Identifier:
		ptr = &node.NodeComments
	case *ast.If:
		ptr = &node.NodeComments
	case *ast.Import:
		ptr = &node.NodeComments
	case *ast.Index:
		ptr = &node.NodeComments
	case *ast.IntLiteral:
		ptr = &node.NodeComments
	case *ast.NoneLiteral:
		ptr = &node.NodeComments
	case *ast.Parameter:
		ptr = &node.NodeComments
	case *ast.Parenthetical:
		ptr = &node.NodeComments
	case *ast.Property:
		ptr = &node.NodeComments
	case *ast.Return:
		ptr = &node.NodeComments
	case *ast.State:
		ptr = &node.NodeComments
	case *ast.StringLiteral:
		ptr = &node.NodeComments
	case *ast.TypeLiteral:
		ptr = &node.NodeComments
	case *ast.Unary:
		ptr = &node.NodeComments
	case *ast.Variable:
		ptr = &node.NodeComments
	case *ast.While:
		ptr = &node.NodeComments
	default:
		return nil
	}
	if *ptr == nil {
		*ptr = &ast.Comments{}
	}
	return node.Comments()
}

// nodes is an [ast.Visitor] that builds and ordererd
// list of nodes that can have comments attached.
type nodes struct {
	nodes []commentable
}

func (v *nodes) VisitAccess(node *ast.Access) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitArgument(node *ast.Argument) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitArrayCreation(node *ast.ArrayCreation) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitAssignment(node *ast.Assignment) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitBinary(node *ast.Binary) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitCall(node *ast.Call) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitCast(node *ast.Cast) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (*nodes) VisitDocumentation(*ast.Documentation) error {
	return nil
}

func (*nodes) VisitBlockComment(*ast.BlockComment) error {
	return nil
}

func (*nodes) VisitCommentStatement(*ast.CommentStatement) error {
	return nil
}

func (*nodes) VisitLineComment(*ast.LineComment) error {
	return nil
}

func (v *nodes) VisitEvent(node *ast.Event) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitFunction(node *ast.Function) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitIdentifier(node *ast.Identifier) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitIf(node *ast.If) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitElseIf(node *ast.ElseIf) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitElse(node *ast.Else) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitExpressionStatement(node *ast.ExpressionStatement) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitImport(node *ast.Import) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitIndex(node *ast.Index) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitBoolLiteral(node *ast.BoolLiteral) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitIntLiteral(node *ast.IntLiteral) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitFloatLiteral(node *ast.FloatLiteral) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitStringLiteral(node *ast.StringLiteral) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitNoneLiteral(node *ast.NoneLiteral) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitParameter(node *ast.Parameter) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitParenthetical(node *ast.Parenthetical) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitProperty(node *ast.Property) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitReturn(node *ast.Return) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (*nodes) VisitScript(*ast.Script) error {
	return nil
}

func (v *nodes) VisitState(node *ast.State) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitTypeLiteral(node *ast.TypeLiteral) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitUnary(node *ast.Unary) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitVariable(node *ast.Variable) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (v *nodes) VisitWhile(node *ast.While) error {
	v.nodes = append(v.nodes, node)
	return nil
}

func (*nodes) VisitErrorStatement(*ast.ErrorStatement) error {
	return nil
}
