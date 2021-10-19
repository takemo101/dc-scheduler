package vm

import (
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/pkg/application"
)

// ToPostMessageIndexMap PostMessageのIndexデータ
func ToPostMessageIndexMap(
	dto application.PostMessageSearchPaginatorDTO,
) helper.DataMap {
	bots := make([]helper.DataMap, len(dto.PostMessages))

	for i, bot := range dto.PostMessages {
		bots[i] = ToPostMessageDetailMap(bot)
	}

	return helper.DataMap{
		"bots":       bots,
		"pagination": helper.StructToJsonMap(&dto.Pagination),
	}
}

// ToPostMessageDetailMap PostMessageのDetailデータ
func ToPostMessageDetailMap(
	dto application.PostMessageDetailDTO,
) helper.DataMap {
	return helper.DataMap{
		"id":                dto.ID,
		"message":           dto.Message,
		"message_type":      dto.MessageType,
		"message_type_name": dto.MessageType.Name(),
		"bot":               ToBotDetailMap(dto.Bot),
		"created_at":        dto.CreatedAt,
		"updated_at":        dto.UpdatedAt,
	}
}
