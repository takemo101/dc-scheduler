package application

import (
	"time"

	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- PostMessageQuery ---

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
	PostMessages []PostMessageDetailDTO `json:"messages"`
	Pagination   Paginator              `json:"pagination"`
}

// --- SentMessageQuery --

// SentMessageQuery
type SentMessageQuery interface {
	Search(SentMessageSearchParameterDTO) (SentMessageSearchPaginatorDTO, error)
}

// SentMessageDetailDTO PostMessage詳細DTO
type SentMessageDetailDTO struct {
	ID          uint                 `json:"id"`
	PostMessage PostMessageDetailDTO `json:"post_message"`
	Message     string               `json:"message"`
	SendedAt    time.Time            `json:"sended_at"`
}

// SentMessageSearchParameterDTO SentMessage一覧取得パラメータDTO
type SentMessageSearchParameterDTO struct {
	Page    int
	Limit   int
	OrderBy string
}

// SentMessageSearchPaginatorDTO SentMessage一覧DTO
type SentMessageSearchPaginatorDTO struct {
	SentMessages []SentMessageDetailDTO `json:"messages"`
	Pagination   Paginator              `json:"pagination"`
}
