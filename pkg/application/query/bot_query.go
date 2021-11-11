package application

import (
	"time"

	common "github.com/takemo101/dc-scheduler/pkg/application/common"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- BotQuery ---

// BotQuery Bot参照
type BotQuery interface {
	Search(BotSearchParameterDTO) (BotSearchPaginatorDTO, error)
	FindByID(domain.BotID) (BotDetailDTO, error)
}

// BotDetailDTO Bot詳細DTO
type BotDetailDTO struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	AtatarURL  string    `json:"avatar_url"`
	AtatarPath string    `json:"avatar_path"`
	Webhook    string    `json:"webhook"`
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// BotSearchParameterDTO Bot一覧取得パラメータDTO
type BotSearchParameterDTO struct {
	Page        int
	Limit       int
	OrderByKey  string
	OrderByType common.OrderByType
}

// BotSearchPaginatorDTO Bot一覧DTO
type BotSearchPaginatorDTO struct {
	Bots       []BotDetailDTO   `json:"bots"`
	Pagination common.Paginator `json:"pagination"`
}
