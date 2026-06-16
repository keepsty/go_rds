package forms

// DBReadQueryExecute SQL 查询执行参数
type DBReadQueryExecute struct {
	DBId      int64  `json:"db_id" binding:"required"`
	Sql       string `json:"sql" binding:"required"`
	Character string `json:"character"`
	QueryHash string `json:"query_hash"`
	Username  string `json:"username"`
}

// RequestGetUserHistoryQueryHandler 查询历史参数
type RequestGetUserHistoryQueryHandler struct {
	PageSize int64  `form:"page_size"`
	Page     int64  `form:"page"`
	DBID     int64  `form:"db_id"`
	Username string `form:"username"`
	Table    string `form:"table"`
}

// RequestGetDBTableInfo 表信息查询参数
type RequestGetDBTableInfo struct {
	DBID  int64  `json:"db_id" form:"db_id" binding:"required"`
	Table string `json:"table" form:"table"`
	Type  string `json:"type" form:"type"`
}
