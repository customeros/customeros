package openai

import (
	"encoding/json"
	"errors"
)

type StrArray []string

var ErrStrArrayUnsupportedType = errors.New("unsupported type, must be string or []string")

func (sa *StrArray) UnmarshalJSON(data []byte) error {
	var jsonObj interface{}
	err := json.Unmarshal(data, &jsonObj)
	if err != nil {
		return err
	}

	switch obj := jsonObj.(type) {

	case string:
		*sa = StrArray([]string{obj})
		return nil

	case []interface{}:
		s := make([]string, 0, len(obj))
		i := 0
		for _, v := range obj {
			value, ok := v.(string)
			if !ok {
				return ErrStrArrayUnsupportedType
			}
			s[i] = value
			i++
		}
		*sa = StrArray(s)
		return nil

	default:
		return ErrStrArrayUnsupportedType
	}
}
