package types

// NewObject returns a new Object type with a
// given declared name and optional parent object.
func NewObject(name string, parent *Object) *Object {
	return &Object{
		name:       name,
		normalized: normalize(name),
		parent:     parent,
	}
}

// Object represents a named typed (i.e. a script).
type Object struct {
	name, normalized string
	parent           *Object
}

// Name returns the declared name for the type.
func (o *Object) Name() string {
	return o.name
}

// Parent is the object this object extends or nil if it extends nothing.
func (o *Object) Parent() *Object {
	return o.parent
}

func (o *Object) String() string {
	return o.name
}

func (*Object) types() {}

func (*Object) scalar() {}

var _ Scalar = (*Object)(nil)
