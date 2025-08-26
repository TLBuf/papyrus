package types

import (
	"fmt"

	"github.com/TLBuf/papyrus/ast"
)

// Resolver creates types for [ast.Node] instances.
//
// Note: Resolver is stateful! In order to resolve script (object) types
// correctly, the script they extend must have been resolved first.
type Resolver struct {
	objects map[string]*Object
}

// Resolve returns the type for a given [ast.Node] or an error if the node is
// not one of the following typed nodes or in the case of an [ast.Script], the
// parent has not previously been resolved.
//
//   - [ast.Script]
//   - [ast.Function]
//   - [ast.Event]
//   - [ast.Property]
//   - [ast.Variable]
//   - [ast.Parameter]
func (r *Resolver) Resolve(node ast.Node) (Type, error) {
	switch node := node.(type) {
	case *ast.Script:
		obj := &Object{
			name:       node.Name.Text,
			normalized: normalize(node.Name.Text),
		}
		if node.Parent != nil {
			var ok bool
			if obj.parent, ok = r.objects[normalize(node.Parent.Text)]; !ok {
				return nil, fmt.Errorf("%v extends unknown parent %q", node, node.Parent.Text)
			}
		}
		return r.record(obj), nil
	case *ast.Function:
		rt, err := r.resolveTypeLiteral(node.ReturnType)
		if err != nil {
			return nil, fmt.Errorf("%v return type: %w", node, err)
		}
		pts := make([]Value, 0, len(node.ParameterList))
		for _, p := range node.ParameterList {
			pt, err := r.resolveTypeLiteral(p.Type)
			if err != nil {
				return nil, fmt.Errorf("%v parameter: %w", node, err)
			}
			pts = append(pts, pt)
		}
		return NewFunction(node.Name.Text, rt, pts...), nil
	case *ast.Event:
		pts := make([]Value, 0, len(node.ParameterList))
		for _, p := range node.ParameterList {
			pt, err := r.resolveTypeLiteral(p.Type)
			if err != nil {
				return nil, fmt.Errorf("%v parameter: %w", node, err)
			}
			pts = append(pts, pt)
		}
		return NewFunction(node.Name.Text, nil, pts...), nil
	case *ast.Property:
		typ, err := r.resolveTypeLiteral(node.Type)
		if err != nil {
			return nil, fmt.Errorf("%v type: %w", node, err)
		}
		return typ, nil
	case *ast.Variable:
		typ, err := r.resolveTypeLiteral(node.Type)
		if err != nil {
			return nil, fmt.Errorf("%v type: %w", node, err)
		}
		return typ, nil
	case *ast.Parameter:
		typ, err := r.resolveTypeLiteral(node.Type)
		if err != nil {
			return nil, fmt.Errorf("%v type: %w", node, err)
		}
		return typ, nil
	default:
		return nil, fmt.Errorf("%v is not typed", node)
	}
}

func (r *Resolver) record(obj *Object) *Object {
	if r.objects == nil {
		r.objects = make(map[string]*Object)
	}
	r.objects[obj.normalized] = obj
	return obj
}

func (r *Resolver) resolveTypeLiteral(literal *ast.TypeLiteral) (Value, error) {
	if literal == nil {
		return nil, nil
	}
	switch name := normalize(literal.Name.Text); name {
	case Bool.normalized:
		if literal.IsArray {
			return BoolArray, nil
		}
		return Bool, nil
	case Int.normalized:
		if literal.IsArray {
			return IntArray, nil
		}
		return Int, nil
	case Float.normalized:
		if literal.IsArray {
			return FloatArray, nil
		}
		return Float, nil
	case String.normalized:
		if literal.IsArray {
			return StringArray, nil
		}
		return String, nil
	default:
		obj, ok := r.objects[name]
		if !ok {
			return nil, fmt.Errorf("%v refers to an unknown object type %q", literal, literal.Name.Text)
		}
		if literal.IsArray {
			return &Array{obj}, nil
		}
		return obj, nil
	}
}
