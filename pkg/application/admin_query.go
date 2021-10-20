package application

import (
	"time"

	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- AdminSearchQuery ---

// AdminQuery Admin参照
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
	Page    int
	Limit   int
	OrderBy string
}

// AdminSearchPaginatorDTO Admin一覧DTO
type AdminSearchPaginatorDTO struct {
	Admins     []AdminDetailDTO `json:"admins"`
	Pagination Paginator        `json:"pagination"`
}
