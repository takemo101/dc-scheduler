package vm

import (
	"github.com/takemo101/dc-scheduler/app/helper"
	application "github.com/takemo101/dc-scheduler/pkg/application/query"
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
		"avatar_url":  dto.AtatarURL,
		"avatar_path": dto.AtatarPath,
		"webhook":     dto.Webhook,
		"active":      dto.Active,
		"user":        ToUserDetailMap(dto.User),
		"created_at":  dto.CreatedAt,
		"updated_at":  dto.UpdatedAt,
	}
}
