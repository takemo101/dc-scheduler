package vm

import (
	"github.com/takemo101/dc-scheduler/app/helper"
	application "github.com/takemo101/dc-scheduler/pkg/application/query"
)

// ToUserIndexMap UserのIndexデータ
func ToUserIndexMap(
	dto application.UserSearchPaginatorDTO,
) helper.DataMap {
	users := make([]helper.DataMap, len(dto.Users))

	for i, user := range dto.Users {
		users[i] = ToUserDetailMap(user)
	}

	return helper.DataMap{
		"users":      users,
		"pagination": helper.StructToJsonMap(&dto.Pagination),
	}
}

// ToUserDetailMap UserのDetailデータ
func ToUserDetailMap(
	dto application.UserDetailDTO,
) helper.DataMap {
	return helper.DataMap{
		"id":             dto.ID,
		"name":           dto.Name,
		"email":          dto.Email,
		"activation_key": dto.ActivationKey,
		"active":         dto.Active,
		"created_at":     dto.CreatedAt,
		"updated_at":     dto.UpdatedAt,
	}
}
