package vm

import (
	"github.com/takemo101/dc-scheduler/app/helper"
	application "github.com/takemo101/dc-scheduler/pkg/application/query"
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

// ToRegularPostIndexMap RegularPostのIndexデータ
func ToRegularPostIndexMap(
	dto application.RegularPostSearchPaginatorDTO,
) helper.DataMap {
	messages := make([]helper.DataMap, len(dto.RegularPosts))

	for i, message := range dto.RegularPosts {
		messages[i] = ToRegularPostDetailMap(message)
	}

	return helper.DataMap{
		"regular_posts": messages,
		"pagination":    helper.StructToJsonMap(&dto.Pagination),
	}
}

// ToRegularPostDetailMap RegularPostのDetailデータ
func ToRegularPostDetailMap(
	dto application.RegularPostDetailDTO,
) helper.DataMap {

	timings := make([]helper.DataMap, len(dto.RegularTimings))
	for i, tm := range dto.RegularTimings {
		timings[i] = ToRegularTimingMap(tm)
	}

	return helper.DataMap{
		"id":              dto.ID,
		"message":         dto.Message,
		"bot":             ToBotDetailMap(dto.Bot),
		"active":          dto.Active,
		"regular_timings": timings,
		"created_at":      dto.CreatedAt,
		"updated_at":      dto.UpdatedAt,
	}
}

// ToRegularTimingMap RegularTimingのデータ
func ToRegularTimingMap(
	dto application.RegularTimingDTO,
) helper.DataMap {
	return helper.DataMap{
		"id":               dto.ID,
		"day_of_week":      dto.DayOfWeek,
		"day_of_week_name": dto.DayOfWeek.Name(),
		"hour_time":        dto.HourTime,
		"hour_time_text":   dto.HourTime.Format("15:04"),
	}
}

// ToApiPostIndexMap ApiPostのIndexデータ
func ToApiPostIndexMap(
	dto application.ApiPostSearchPaginatorDTO,
) helper.DataMap {
	messages := make([]helper.DataMap, len(dto.ApiPosts))

	for i, message := range dto.ApiPosts {
		messages[i] = ToApiPostDetailMap(message)
	}

	return helper.DataMap{
		"api_posts":  messages,
		"pagination": helper.StructToJsonMap(&dto.Pagination),
	}
}

// ToApiPostDetailMap ApiPostのDetailデータ
func ToApiPostDetailMap(
	dto application.ApiPostDetailDTO,
) helper.DataMap {
	return helper.DataMap{
		"id":         dto.ID,
		"message":    dto.Message,
		"api_key":    dto.ApiKey,
		"bot":        ToBotDetailMap(dto.Bot),
		"created_at": dto.CreatedAt,
		"updated_at": dto.UpdatedAt,
	}
}
