package rules

import (
	"github.com/keepsty/go_rds/internal/inspect/controllers/logics"
	"github.com/keepsty/go_rds/internal/inspect/controllers/traverses"

	"github.com/pingcap/tidb/pkg/parser/ast"
)

func CreateDatabaseRules() []Rule {
	return []Rule{
		{
			Hint:      "CreateDatabase#检查DB是否存在",
			CheckFunc: (*Rule).RuleCreateDatabaseIsExist,
		},
	}
}

// RuleCreateDatabaseIsExist
func (r *Rule) RuleCreateDatabaseIsExist(tistmt *ast.StmtNode) {
	v := &traverses.TraverseCreateDatabaseIsExist{}
	(*tistmt).Accept(v)
	logics.LogicCreateDatabaseIsExist(v, r.RuleHint)
}
