package forms

// ParamClusterList 集群列表查询参数
type ParamClusterList struct {
	Page         int `json:"page" form:"page"`
	Size         int `json:"size" form:"size"`
	ProductionID int `json:"production_id"`
}

// InstanceDetail 实例状态更新参数
type InstanceDetail struct {
	ID             int64  `json:"id" form:"id"`
	Status         int64  `json:"status" form:"status"`
	HostId         int64  `json:"host_id" form:"host_id"`
	ServiceGroupID int64  `json:"servicegroup_id" form:"servicegroup_id"`
	Port           int64  `json:"port" form:"port"`
	Role           int64  `json:"role" form:"role"`
	Purpose        string `json:"purpose" form:"purpose"`
	MysqlVersion   string `json:"mysql_version" form:"mysql_version"`
}
