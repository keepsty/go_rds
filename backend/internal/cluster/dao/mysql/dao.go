package mysql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/keepsty/go_rds/internal/config"
)

var db *sqlx.DB

// Init 初始化 sqlx 数据库连接（独立于主项目的 GORM 连接）
func Init(cfg *config.Database) (err error) {
	// 先连接不指定数据库，确保目标库存在且字符集正确
	initDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=%s&parseTime=True&loc=Local",
		cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.Charset)
	initDB, err := sqlx.Connect("mysql", initDSN)
	if err != nil {
		return fmt.Errorf("连接MySQL失败: %w", err)
	}

	// 创建数据库（如不存在）并指定字符集
	createSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET %s COLLATE %s_general_ci",
		cfg.Database, cfg.Charset, cfg.Charset)
	if _, err = initDB.Exec(createSQL); err != nil {
		initDB.Close()
		return fmt.Errorf("创建数据库失败: %w", err)
	}
	initDB.Close()

	// 正式连接目标数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.UserName, cfg.Password, cfg.Host, cfg.Port, cfg.Database, cfg.Charset)
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		return fmt.Errorf("连接目标数据库失败: %w", err)
	}
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	return
}

// Close 关闭数据库连接
func Close() {
	if db != nil {
		_ = db.Close()
	}
}
