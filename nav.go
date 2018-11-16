package query

type Nav interface {
	//Current() Node
	NodeType() NodeType
	Value() string
	Copy() Nav
	LocalName() string
	Prefix() string
	MoveToRoot()
	MoveToParent() bool
	MoveToNextAttribute() bool
	MoveToChild() bool
	MoveToFirst() bool
	MoveToNext() bool
	MoveToPrevious() bool
	MoveTo(Nav) bool
	//String() string
}
