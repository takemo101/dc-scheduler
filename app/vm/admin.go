package vm

import (
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/pkg/application"
)

// ToAdminIndexMap AdminのIndexデータ
func ToAdminIndexMap(
	dto application.AdminSearchPaginatorDTO,
) helper.DataMap {
	admins := make([]helper.DataMap, len(dto.Admins))

	for i, admin := range dto.Admins {
		admins[i] = ToAdminDetailMap(admin)
	}

	return helper.DataMap{
		"admins":     admins,
		"pagination": helper.StructToJsonMap(&dto.Pagination),
	}
}

// ToAdminDetailMap AdminのDetailデータ
func ToAdminDetailMap(
	dto application.AdminDetailDTO,
) helper.DataMap {
	return helper.DataMap{
		"id":         dto.ID,
		"name":       dto.Name,
		"email":      dto.Email,
		"role_type":  dto.Role.Value(),
		"role_name":  dto.Role.Name(),
		"created_at": dto.CreatedAt,
		"updated_at": dto.UpdatedAt,
	}
}
