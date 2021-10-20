package application

import (
	"time"

	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- PostMessageSearchQuery ---

// PostMessageQuery PostMessage参照
type PostMessageQuery interface {
	Search(PostMessageSearchParameterDTO) (PostMessageSearchPaginatorDTO, error)
	FindByID(domain.PostMessageID) (PostMessageDetailDTO, error)
}

// PostMessageDetailDTO PostMessage詳細DTO
type PostMessageDetailDTO struct {
	ID          uint               `json:"id"`
	Message     string             `json:"message"`
	MessageType domain.MessageType `json:"message_type"`
	Bot         BotDetailDTO       `json:"bot"`
	CreatedAt   time.Time          `json:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at"`
}

// PostMessageSearchParameterDTO PostMessage一覧取得パラメータDTO
type PostMessageSearchParameterDTO struct {
	Page    int
	Limit   int
	OrderBy string
}

// PostMessageSearchPaginatorDTO PostMessage一覧DTO
type PostMessageSearchPaginatorDTO struct {
	PostMessages []PostMessageDetailDTO `json:"bots"`
	Pagination   Paginator              `json:"pagination"`
}
