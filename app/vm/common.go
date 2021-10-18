package vm

import (
	"github.com/takemo101/dc-scheduler/app/helper"
	"github.com/takemo101/dc-scheduler/pkg/domain"
)

func ToKeyValueMap(
	kvs []domain.KeyValue,
) []helper.DataMap {
	mapData := make([]helper.DataMap, len(kvs))

	for i, kv := range kvs {
		mapData[i] = helper.DataMap{
			"key":   kv.Key,
			"value": kv.Value,
		}
	}

	return mapData
}
