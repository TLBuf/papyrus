package parser

import "github.com/TLBuf/papyrus/pkg/ast"

// attachLooseComments updates the [ast.Trivia] on nodes in place with
// references to the appropriate comments.
func attachLooseComments(script *ast.Script, comments []ast.LooseComment) error {
	return newError(script.Location, "attaching loose comments is not yet implemented: failed to attach %d comments", len(comments))
}
