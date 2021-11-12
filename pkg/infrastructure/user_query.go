package infrastructure

import (
	"github.com/takemo101/dc-scheduler/core"
	application "github.com/takemo101/dc-scheduler/pkg/application/query"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

// --- UserQuery ---

// UserQuery application.UserQueryの実装
type UserQuery struct {
	db core.Database
}

// NewUserQuery コンストラクタ
func NewUserQuery(db core.Database) application.UserQuery {
	return UserQuery{
		db,
	}
}

// Search Userの一覧取得
func (query UserQuery) Search(parameter application.UserSearchParameterDTO) (dto application.UserSearchPaginatorDTO, err error) {
	var models []User

	paging := NewGormPaging(
		query.db.GormDB,
		parameter.Page,
		parameter.Limit,
		[]string{parameter.OrderByType.ToQuery(parameter.OrderByKey)},
	)

	paginator, err := paging.Paging(&models)
	if err != nil {
		return dto, err
	}

	dto.Pagination = paginator

	admins := make([]application.UserDetailDTO, len(models))
	for i, m := range models {
		admins[i] = CreateUserDetailDTOFromModel(m)
	}

	dto.Users = admins

	return dto, err
}

// FindByID Userの詳細取得
func (query UserQuery) FindByID(id domain.UserID) (dto application.UserDetailDTO, err error) {
	model := User{}

	if err = query.db.GormDB.Where("id = ?", id.Value()).First(&model).Error; err != nil {
		return dto, err
	}

	return CreateUserDetailDTOFromModel(model), err
}

// CreateUserDetailDTOFromModel UserからUserDetailDTOを生成する
func CreateUserDetailDTOFromModel(model User) application.UserDetailDTO {
	return application.UserDetailDTO{
		ID:            model.ID,
		Name:          model.Name,
		Email:         model.Email,
		ActivationKey: model.ActivationKey,
		Active:        model.Active.Bool,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}
}
