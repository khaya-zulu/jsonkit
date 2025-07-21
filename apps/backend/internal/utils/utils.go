package utils

import "encoding/json"

func AddFieldToRawMessage(rawMsg json.RawMessage, key string, value interface{}) (json.RawMessage, error) {
    var inputMap map[string]interface{}
    if err := json.Unmarshal(rawMsg, &inputMap); err != nil {
        return nil, err
    }
    
    inputMap[key] = value
    
    return json.Marshal(inputMap)
}