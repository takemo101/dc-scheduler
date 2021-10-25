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
	messages := ToSentMessagesMap(dto.SentMessages)

	return helper.DataMap{
		"sent_messages": messages,
		"pagination":    helper.StructToJsonMap(&dto.Pagination),
	}
}

func ToSentMessagesMap(list []application.SentMessageDetailDTO) []helper.DataMap {
	messages := make([]helper.DataMap, len(list))

	for i, message := range list {
		messages[i] = ToSentMessageDetailMap(message)
	}

	return messages
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

// ToImmediatePostIndexMap ImmediatePostのIndexデータ
func ToImmediatePostIndexMap(
	dto application.ImmediatePostSearchPaginatorDTO,
) helper.DataMap {
	messages := make([]helper.DataMap, len(dto.ImmediatePosts))

	for i, message := range dto.ImmediatePosts {
		messages[i] = ToImmediatePostDetailMap(message)
	}

	return helper.DataMap{
		"immediate_posts": messages,
		"pagination":      helper.StructToJsonMap(&dto.Pagination),
	}
}

// ToImmediatePostDetailMap ImmediatePostのDetailデータ
func ToImmediatePostDetailMap(
	dto application.ImmediatePostDetailDTO,
) helper.DataMap {
	return helper.DataMap{
		"id":         dto.ID,
		"message":    dto.Message,
		"bot":        ToBotDetailMap(dto.Bot),
		"created_at": dto.CreatedAt,
		"updated_at": dto.UpdatedAt,
	}
}

// ToSchedulePostIndexMap SchedulePostのIndexデータ
func ToSchedulePostIndexMap(
	dto application.SchedulePostSearchPaginatorDTO,
) helper.DataMap {
	messages := make([]helper.DataMap, len(dto.SchedulePosts))

	for i, message := range dto.SchedulePosts {
		messages[i] = ToSchedulePostDetailMap(message)
	}

	return helper.DataMap{
		"schedule_posts": messages,
		"pagination":     helper.StructToJsonMap(&dto.Pagination),
	}
}

// ToSchedulePostDetailMap SchedulePostのDetailデータ
func ToSchedulePostDetailMap(
	dto application.SchedulePostDetailDTO,
) helper.DataMap {
	return helper.DataMap{
		"id":             dto.ID,
		"message":        dto.Message,
		"bot":            ToBotDetailMap(dto.Bot),
		"is_sended":      dto.IsSended,
		"reservation_at": dto.ReservationAt,
		"created_at":     dto.CreatedAt,
		"updated_at":     dto.UpdatedAt,
	}
}
