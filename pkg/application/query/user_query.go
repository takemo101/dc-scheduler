package application

import (
	"time"

	common "github.com/takemo101/dc-scheduler/pkg/application/common"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- UserQuery ---

// UserQuery User参照
type UserQuery interface {
	Search(UserSearchParameterDTO) (UserSearchPaginatorDTO, error)
	FindByID(domain.UserID) (UserDetailDTO, error)
}

// UserDetailDTO User詳細DTO
type UserDetailDTO struct {
	ID            uint      `json:"id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	ActivationKey string    `json:"activation_key"`
	Active        bool      `json:"active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// UserSearchParameterDTO User一覧取得パラメータDTO
type UserSearchParameterDTO struct {
	Page        int
	Limit       int
	OrderByKey  string
	OrderByType common.OrderByType
}

// UserSearchPaginatorDTO User一覧DTO
type UserSearchPaginatorDTO struct {
	Users      []UserDetailDTO  `json:"admins"`
	Pagination common.Paginator `json:"pagination"`
}
