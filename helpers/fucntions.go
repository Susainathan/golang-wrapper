package helpers

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func ConvertStrToMap(c *gin.Context, value string) (map[string]interface{}, bool) {
	var data map[string]interface{}

	if err := json.Unmarshal([]byte(value), &data); err != nil {
		ApiException(c, ExceptionStruct{ExceptionKey: "ERR-INPUT", Message: fmt.Sprintf("Invalid input given, Please provide the valid input.")}, 422)
		return data, true
	}
	return data, false
}

func GenerateOutputData(data map[string]interface{}) OutputData {
	input := OutputData{
		Event:       data["ev"].(string),
		EventType:   data["et"].(string),
		AppId:       data["id"].(string),
		UserId:      data["uid"].(string),
		MessageId:   data["mid"].(string),
		PageTitle:   data["t"].(string),
		PageURL:     data["p"].(string),
		BrowserLang: data["l"].(string),
		ScreenSize:  data["sc"].(string),
		Attributes:  make(map[string]interface{}),
		UserTraits:  make(map[string]interface{}),
	}

	for key, _ := range data {
		if strings.HasPrefix(key, "atrk") {
			atrkIndexNumber := key[4:]
			input.Attributes[data["atrk"+atrkIndexNumber].(string)] = map[string]interface{}{
				"value": data["atrv"+atrkIndexNumber].(string),
				"type":  data["atrt"+atrkIndexNumber].(string),
			}
		} else if strings.HasPrefix(key, "uatrk") {
			uatrkIndexNumber := key[5:]
			input.UserTraits[data["uatrk"+uatrkIndexNumber].(string)] = map[string]interface{}{
				"value": data["uatrv"+uatrkIndexNumber].(string),
				"type":  data["uatrt"+uatrkIndexNumber].(string),
			}
		}
	}
	return input
}
