package vm

import (
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/pkg/application"
)

// ToPostMessageIndexMap PostMessageのIndexデータ
func ToPostMessageIndexMap(
	dto application.PostMessageSearchPaginatorDTO,
) helper.DataMap {
	messages := make([]helper.DataMap, len(dto.PostMessages))

	for i, message := range dto.PostMessages {
		messages[i] = ToPostMessageDetailMap(message)
	}

	return helper.DataMap{
		"post_messages": messages,
		"pagination":    helper.StructToJsonMap(&dto.Pagination),
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

// ToSentMessageHistoryMap SentMessageのHistoryデータ
func ToSentMessageHistoryMap(
	dto application.SentMessageSearchPaginatorDTO,
) helper.DataMap {
	messages := make([]helper.DataMap, len(dto.SentMessages))

	for i, message := range dto.SentMessages {
		messages[i] = ToSentMessageDetailMap(message)
	}

	return helper.DataMap{
		"sent_messages": messages,
		"pagination":    helper.StructToJsonMap(&dto.Pagination),
	}
}

// ToSentMessageDetailMap SentMessageのDetailデータ
func ToSentMessageDetailMap(
	dto application.SentMessageDetailDTO,
) helper.DataMap {
	return helper.DataMap{
		"id":           dto.ID,
		"post_message": ToPostMessageDetailMap(dto.PostMessage),
		"message":      dto.Message,
		"sended_at":    dto.SendedAt,
	}
}
