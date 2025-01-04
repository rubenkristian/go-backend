package commons

import (
	"github.com/rubenkristian/backend/configs"
	"github.com/rubenkristian/backend/utils"
	"gorm.io/gorm"
)

type AppConfig struct {
	Db        *gorm.DB
	Mailer    *utils.Emailer
	Env       *configs.EnvConfig
	S3Service *utils.S3Service
}

type PaginationParams struct {
	Take   int    `query:"take"`
	Skip   int    `query:"skip"`
	Search string `query:"search"`
	Sort   string `query:"sort"`
	SortBy string `query:"sortBy"`
}

func (paginationParams *PaginationParams) SetParams(take int, sort, sortBy string) {
	if paginationParams.Take == 0 {
		paginationParams.Take = 10
	}

	if paginationParams.Sort == "" {
		paginationParams.Sort = sort
	}

	if paginationParams.SortBy == "" {
		paginationParams.SortBy = sortBy
	}
}
