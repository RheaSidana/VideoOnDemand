package videoEncryption

import (
	"encoding/json"
	"vod/model"
)

type VideoEncryptRequest struct {
	VideoID int `json:"videoID" binding:"required"`
}

type ErrorResponse struct {
	Message string
}

type VideoEncryptResponse struct {
	Message            string
	VideoMD            model.VideoMetaData
	EncodedToEncrypted map[model.VideoLinks]model.VideoLinks
}

func (ver VideoEncryptResponse) MarshalJSON() ([]byte, error) {
	// Create a temporary map with string keys for marshaling
	tempMap := make(map[string]model.VideoLinks)
	for k, v := range ver.EncodedToEncrypted {
		tempMap[k.ToString()] = v
	}

	// Create a custom JSON structure
	customJSON := struct {
		Message            string                    `json:"message"`
		VideoMD            model.VideoMetaData       `json:"video_metadata"`
		EncodedToEncrypted map[string]model.VideoLinks `json:"encoded_to_encrypted"`
	}{
		Message:            ver.Message,
		VideoMD:            ver.VideoMD,
		EncodedToEncrypted: tempMap,
	}

	return json.Marshal(customJSON)
}

func (ver *VideoEncryptResponse) UnmarshalJSON(data []byte) error {
	// Define a custom JSON structure for unmarshaling
	customJSON := struct {
		Message            string                    `json:"message"`
		VideoMD            model.VideoMetaData       `json:"video"`
		EncodedToEncrypted map[string]model.VideoLinks `json:"encoded_to_encrypted"`
	}{}

	// Unmarshal JSON data into customJSON
	if err := json.Unmarshal(data, &customJSON); err != nil {
		return err
	}

	// Create the VideoEncryptResponse struct from the custom JSON structure
	ver.Message = customJSON.Message
	ver.VideoMD = customJSON.VideoMD

	// Create the EncodedToEncrypted map and convert string keys back to model.VideoLinks
	ver.EncodedToEncrypted = make(map[model.VideoLinks]model.VideoLinks)
	for k, v := range customJSON.EncodedToEncrypted {
		key, err := model.StringToVideoLinks(k)
		if err != nil {
			return err
		}
		ver.EncodedToEncrypted[key] = v
	}

	return nil
}