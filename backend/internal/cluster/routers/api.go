package routers

import (
	"github.com/gin-gonic/gin"

	"github.com/keepsty/go_rds/internal/cluster/views"
)

func RegisterApiRoutes(cluster *gin.RouterGroup) {
	// 集群列表与详情
	cluster.GET("/clusters", views.ClustersHandler)
	cluster.GET("/clusters/:id", views.GetClusterByIDHandler)
	cluster.GET("/clusters/instances/:id", views.GetClusterInsByIDHandler)
	cluster.POST("/clusters/instances/status", views.SetInstanceStatusByIDHandler)
	cluster.GET("/clusters/databases/:id", views.GetClusterDBsByIDHandler)
	cluster.GET("/clusters/database-detail/:id", views.GetClusterDBDetailByIDHandler)
	cluster.GET("/clusters/proxysql/:id", views.GetClusterProxyListByIDHandler)

	// SaltStack 自动化部署
	cluster.POST("/clusters/addmysqlcluster", views.AddMySQLClusterHandler)

	// 数据库查询与数据字典
	cluster.GET("/clusters/dboptions/:username", views.GetDBsByClusterNameHandler)
	cluster.GET("/databases/:id", views.GetDBsHandler)
	cluster.GET("/clusters/dbtables/:sg_id/:db_id", views.GetTbsByClusterNameHandler)
	cluster.POST("/clusters/db/query/execute", views.ClusterDBReadQueryExecuteHandler)
	cluster.POST("/clusters/db/query/datadict", views.ClusterGetDBDataDictHandler)
	cluster.GET("/clusters/db/query/tableinfo", views.ClusterGetDBTableInfoHandler)
	cluster.GET("/clusters/db/query/history", views.GetUserHistorySqlHandler)
}
