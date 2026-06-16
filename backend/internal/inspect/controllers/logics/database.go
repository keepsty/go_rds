package logics

import (
	"github.com/keepsty/go_rds/internal/inspect/controllers"
	"github.com/keepsty/go_rds/internal/inspect/controllers/dao"
	"github.com/keepsty/go_rds/internal/inspect/controllers/traverses"
)

// LogicCreateDatabaseIsExist
func LogicCreateDatabaseIsExist(v *traverses.TraverseCreateDatabaseIsExist, r *controllers.RuleHint) {
	if msg, err := dao.CheckIfDatabaseExists(v.Name, r.DB); err == nil {
		r.Warn("数据库已存在：" + msg)
		r.IsBreak = true
	}
}
