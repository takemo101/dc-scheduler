package vm

import (
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/pkg/application"
)

// ToBotIndexMap BotのIndexデータ
func ToBotIndexMap(
	dto application.BotSearchPaginatorDTO,
) helper.DataMap {
	bots := make([]helper.DataMap, len(dto.Bots))

	for i, bot := range dto.Bots {
		bots[i] = ToBotDetailMap(bot)
	}

	return helper.DataMap{
		"bots":       bots,
		"pagination": helper.StructToJsonMap(&dto.Pagination),
	}
}

// ToBotDetailMap BotのDetailデータ
func ToBotDetailMap(
	dto application.BotDetailDTO,
) helper.DataMap {
	return helper.DataMap{
		"id":          dto.ID,
		"name":        dto.Name,
		"avator_url":  dto.AvatorURL,
		"avator_path": dto.AvatorPath,
		"webhook":     dto.Webhook,
		"active":      dto.Active,
		"created_at":  dto.CreatedAt,
		"updated_at":  dto.UpdatedAt,
	}
}
