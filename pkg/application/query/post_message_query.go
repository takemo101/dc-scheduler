package application

import (
	"time"

	common "github.com/takemo101/dc-scheduler/pkg/application/common"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- PostMessageQuery ---

// PostMessageQuery PostMessage参照
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
	Page        int
	Limit       int
	OrderByKey  string
	OrderByType common.OrderByType
}

// PostMessageSearchPaginatorDTO PostMessage一覧DTO
type PostMessageSearchPaginatorDTO struct {
	PostMessages []PostMessageDetailDTO `json:"messages"`
	Pagination   common.Paginator       `json:"pagination"`
}

// --- SentMessageQuery --

// SentMessageQuery
type SentMessageQuery interface {
	Search(SentMessageSearchParameterDTO) (SentMessageSearchPaginatorDTO, error)
	SearchByUserID(SentMessageSearchParameterDTO, domain.UserID) (SentMessageSearchPaginatorDTO, error)
	RecentlyList(uint) ([]SentMessageDetailDTO, error)
	RecentlyListByUserID(domain.UserID, uint) ([]SentMessageDetailDTO, error)
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
	Page        int
	Limit       int
	OrderByKey  string
	OrderByType common.OrderByType
}

// SentMessageSearchPaginatorDTO SentMessage一覧DTO
type SentMessageSearchPaginatorDTO struct {
	SentMessages []SentMessageDetailDTO `json:"messages"`
	Pagination   common.Paginator       `json:"pagination"`
}

// --- ImmediatePostQuery ---

// ImmediatePostQuery ImmediatePost参照
type ImmediatePostQuery interface {
	Search(ImmediatePostSearchParameterDTO) (ImmediatePostSearchPaginatorDTO, error)
	SearchByUserID(ImmediatePostSearchParameterDTO, domain.UserID) (ImmediatePostSearchPaginatorDTO, error)
}

// ImmediatePostDetailDTO ImmediatePost詳細DTO
type ImmediatePostDetailDTO struct {
	ID        uint         `json:"id"`
	Message   string       `json:"message"`
	Bot       BotDetailDTO `json:"bot"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}

// ImmediatePostSearchParameterDTO ImmediatePost一覧取得パラメータDTO
type ImmediatePostSearchParameterDTO struct {
	Page        int
	Limit       int
	OrderByKey  string
	OrderByType common.OrderByType
}

// ImmediatePostSearchPaginatorDTO ImmediatePost一覧DTO
type ImmediatePostSearchPaginatorDTO struct {
	ImmediatePosts []ImmediatePostDetailDTO `json:"messages"`
	Pagination     common.Paginator         `json:"pagination"`
}

// --- SchedulePostQuery ---

// SchedulePostQuery SchedulePost参照
type SchedulePostQuery interface {
	Search(SchedulePostSearchParameterDTO) (SchedulePostSearchPaginatorDTO, error)
	SearchByUserID(SchedulePostSearchParameterDTO, domain.UserID) (SchedulePostSearchPaginatorDTO, error)
	FindByID(domain.PostMessageID) (SchedulePostDetailDTO, error)
}

// SchedulePostDetailDTO SchedulePost詳細DTO
type SchedulePostDetailDTO struct {
	ID            uint         `json:"id"`
	Message       string       `json:"message"`
	Bot           BotDetailDTO `json:"bot"`
	IsSended      bool         `json:"is_sended"`
	ReservationAt time.Time    `json:"reservation_at"`
	CreatedAt     time.Time    `json:"created_at"`
	UpdatedAt     time.Time    `json:"updated_at"`
}

// SchedulePostSearchParameterDTO SchedulePost一覧取得パラメータDTO
type SchedulePostSearchParameterDTO struct {
	Page        int
	Limit       int
	OrderByKey  string
	OrderByType common.OrderByType
}

// SchedulePostSearchPaginatorDTO SchedulePost一覧DTO
type SchedulePostSearchPaginatorDTO struct {
	SchedulePosts []SchedulePostDetailDTO `json:"messages"`
	Pagination    common.Paginator        `json:"pagination"`
}
