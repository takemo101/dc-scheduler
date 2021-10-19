package application

import (
	"time"

	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- BotSearchQuery ---

// BotQuery Bot参照
type BotQuery interface {
	Search(BotSearchParameterDTO) (BotSearchPaginatorDTO, error)
	FindByID(domain.BotID) (BotDetailDTO, error)
}

// BotDetailDTO Bot詳細DTO
type BotDetailDTO struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	AvatorURL  string    `json:"avator_url"`
	AvatorPath string    `json:"avator_path"`
	Webhook    string    `json:"webhook"`
	Active     bool      `json:"active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// BotSearchParameterDTO Bot一覧取得パラメータDTO
type BotSearchParameterDTO struct {
	Page    int
	Limit   int
	OrderBy string
}

// BotSearchPaginatorDTO Bot一覧DTO
type BotSearchPaginatorDTO struct {
	Bots       []BotDetailDTO `json:"bots"`
	Pagination Paginator      `json:"pagination"`
}
