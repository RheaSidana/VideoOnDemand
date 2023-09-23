package videoMetadata

import (
	"encoding/json"
	"vod/model"
)

type RedisData struct {
	VideoMetaData model.VideoMetaData
}

func (rd RedisData) MarshalBinary() ([]byte, error) {
    return json.Marshal(rd) // You can choose any binary serialization format you prefer
}

func (rd *RedisData) UnmarshalBinary(data []byte) error {
    return json.Unmarshal(data, rd) // You should use the same serialization format as MarshalBinary
} 