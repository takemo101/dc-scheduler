package application

import (
	"time"

	common "github.com/takemo101/dc-scheduler/pkg/application/common"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- ApiPostQuery ---

// ApiPostQuery ApiPost参照
type ApiPostQuery interface {
	Search(ApiPostSearchParameterDTO) (ApiPostSearchPaginatorDTO, error)
	SearchByUserID(ApiPostSearchParameterDTO, domain.UserID) (ApiPostSearchPaginatorDTO, error)
}

// ApiPostDetailDTO ApiPost詳細DTO
type ApiPostDetailDTO struct {
	ID        uint         `json:"id"`
	Message   string       `json:"message"`
	ApiKey    string       `json:"api_key"`
	Bot       BotDetailDTO `json:"bot"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// ApiPostSearchParameterDTO ApiPost一覧取得パラメータDTO
type ApiPostSearchParameterDTO struct {
	Page        int
	Limit       int
	OrderByKey  string
	OrderByType common.OrderByType
}

// ApiPostSearchPaginatorDTO ApiPost一覧DTO
type ApiPostSearchPaginatorDTO struct {
	ApiPosts   []ApiPostDetailDTO `json:"messages"`
	Pagination common.Paginator   `json:"pagination"`
}
