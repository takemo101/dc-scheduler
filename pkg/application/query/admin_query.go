package application

import (
	"time"

	common "github.com/takemo101/dc-scheduler/pkg/application/common"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- AdminQuery ---

// AdminQuery Admin参照
type AdminQuery interface {
	Search(AdminSearchParameterDTO) (AdminSearchPaginatorDTO, error)
	FindByID(domain.AdminID) (AdminDetailDTO, error)
}

// AdminDetailDTO Admin詳細DTO
type AdminDetailDTO struct {
	ID        uint             `json:"id"`
	Name      string           `json:"name"`
	Email     string           `json:"email"`
	Role      domain.AdminRole `json:"role"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

// AdminSearchParameterDTO Admin一覧取得パラメータDTO
type AdminSearchParameterDTO struct {
	Page        int
	Limit       int
	OrderByKey  string
	OrderByType common.OrderByType
}

// AdminSearchPaginatorDTO Admin一覧DTO
type AdminSearchPaginatorDTO struct {
	Admins     []AdminDetailDTO `json:"admins"`
	Pagination common.Paginator `json:"pagination"`
}
