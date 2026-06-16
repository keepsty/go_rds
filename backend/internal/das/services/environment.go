package services

import (
	"github.com/keepsty/go_rds/internal/global"

	"github.com/keepsty/go_rds/internal/common/models"

	"github.com/gin-gonic/gin"
)

type GetEnvironmentsService struct {
	C *gin.Context
}

func (s *GetEnvironmentsService) Run() ([]models.InsightInstanceEnvironments, error) {
	var results []models.InsightInstanceEnvironments
	global.App.DB.Table("`insight_instance_environments` a").
		Select("a.`name`").
		Scan(&results)
	return results, nil
}
