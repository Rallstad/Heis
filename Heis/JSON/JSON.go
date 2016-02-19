package JSON

import "encoding/json"

func Encode(v interface{}) ([]byte, error){
	return json.Marshal(v)
}

func decode(data []byte, v interface{}) (error){
	return json.Unmarshal(data, v)
}