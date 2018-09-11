package interpret

import (
	"go/ast"
)

type BlockStmt struct {
	//original  *ast.BlockStmt
	original []ast.Stmt
}

func NewBlockStmt(bs []ast.Stmt) *BlockStmt {
	blockStmt := &BlockStmt{
		original: bs,
	}

	return blockStmt
}

func (bs *BlockStmt) Describe() (string, error) {
	return "", nil
}
