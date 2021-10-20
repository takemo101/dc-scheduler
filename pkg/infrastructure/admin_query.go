package infrastructure

import (
	"github.com/takemo101/dc-scheduler/core"
	"github.com/takemo101/dc-scheduler/pkg/application"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- AdminQuery ---

// AdminQuery application.AdminQueryの実装
type AdminQuery struct {
	db core.Database
}

// NewAdminQuery コンストラクタ
func NewAdminQuery(db core.Database) application.AdminQuery {
	return AdminQuery{
		db,
	}
}

// Search Adminの一覧取得
func (query AdminQuery) Search(parameter application.AdminSearchParameterDTO) (dto application.AdminSearchPaginatorDTO, err error) {
	var models []Admin

	paging := NewGormPaging(
		query.db.GormDB,
		parameter.Page,
		parameter.Limit,
		[]string{parameter.OrderBy},
	)

	paginator, err := paging.Paging(&models)
	if err != nil {
		return dto, err
	}

	dto.Pagination = paginator

	admins := make([]application.AdminDetailDTO, len(models))
	for i, m := range models {
		admins[i] = CreateAdminDetailDTOFromModel(m)
	}

	dto.Admins = admins

	return dto, err
}

// FindByID Adminの詳細取得
func (query AdminQuery) FindByID(id domain.AdminID) (dto application.AdminDetailDTO, err error) {
	model := Admin{}

	if err = query.db.GormDB.Where("id = ?", id.Value()).First(&model).Error; err != nil {
		return dto, err
	}

	return CreateAdminDetailDTOFromModel(model), err
}

// CreateAdminDetailDTOFromModel AdminからAdminDetailDTOを生成する
func CreateAdminDetailDTOFromModel(model Admin) application.AdminDetailDTO {
	return application.AdminDetailDTO{
		ID:        model.ID,
		Name:      model.Name,
		Email:     model.Email,
		Role:      model.Role,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}
