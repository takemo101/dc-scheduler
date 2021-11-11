package application

import (
	"time"

	common "github.com/takemo101/dc-scheduler/pkg/application/common"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- RegularPostQuery ---

// RegularPostQuery RegularPost参照
type RegularPostQuery interface {
	Search(RegularPostSearchParameterDTO) (RegularPostSearchPaginatorDTO, error)
	FindByID(domain.PostMessageID) (RegularPostDetailDTO, error)
}

// RegularPostDetailDTO RegularPost詳細DTO
type RegularPostDetailDTO struct {
	ID             uint               `json:"id"`
	Message        string             `json:"message"`
	Bot            BotDetailDTO       `json:"bot"`
	Active         bool               `json:"active"`
	RegularTimings []RegularTimingDTO `json:"timings"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

// RegularTimingDTO RegularTiming詳細DTO
type RegularTimingDTO struct {
	ID        uint             `json:"id"`
	DayOfWeek domain.DayOfWeek `json:"day_of_week"`
	HourTime  time.Time        `json:"hour_time"`
}

// RegularPostSearchParameterDTO RegularPost一覧取得パラメータDTO
type RegularPostSearchParameterDTO struct {
	Page        int
	Limit       int
	OrderByKey  string
	OrderByType common.OrderByType
}

// RegularPostSearchPaginatorDTO RegularPost一覧DTO
type RegularPostSearchPaginatorDTO struct {
	RegularPosts []RegularPostDetailDTO `json:"messages"`
	Pagination   common.Paginator       `json:"pagination"`
}
