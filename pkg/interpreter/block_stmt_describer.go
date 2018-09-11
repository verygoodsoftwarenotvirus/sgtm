package interpret

import (
	"go/ast"
)

type BlockStmt struct {
	//original  *ast.BlockStmt
	original []ast.Stmt
	Verbosity verbosity
}

func NewBlockStmt(bs []ast.Stmt, v verbosity) *BlockStmt {
	blockStmt := &BlockStmt{
		original:  bs,
		Verbosity: v,
	}

	return blockStmt
}

func (bs *BlockStmt) Describe() (string, error) {
	return "", nil
}