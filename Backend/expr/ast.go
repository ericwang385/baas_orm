package expr

type AstType int

const (
	ASTBinaryOP = iota
	ASTUnaryOP
	ASTValueInt
	ASTValueFloat
	ASTValueBool
	ASTValueText
	ASTFuncCall
	ASTVariable
	ASTColumn
)

type AstNode struct {
	Type     AstType
	Value    string
	Children []*AstNode
	Offset   int
}

func newAst(t AstType, value string, offset int, children ...*AstNode) *AstNode {
	return &AstNode{
		Type:     t,
		Value:    value,
		Children: children,
		Offset:   offset,
	}
}
