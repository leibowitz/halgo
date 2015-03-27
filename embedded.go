package halgo

import (
	"encoding/json"
)

// Embedded represents a collection of objects.
type Embedded struct {
	Items map[string]EmbeddedItems `json:"_embedded,omitempty"`
}

// Add creates multiple objects with the same relation.
//
//     Add("abc", halgo.Object{}, halgo.Object{})
func (l Embedded) Add(rel string, objects ...interface{}) Embedded {
	if l.Items == nil {
		l.Items = make(map[string]EmbeddedItems)
	}

	set, exists := l.Items[rel]

	if exists {
		set = append(set, objects...)
	} else {
		set = make([]interface{}, len(objects))
		copy(set, objects)
	}

	l.Items[rel] = set

	return l
}

type EmbeddedItems []interface{}

func (l EmbeddedItems) MarshalJSON() ([]byte, error) {
	if len(l) == 1 {
		return json.Marshal(l[0])
	}

	other := make([]interface{}, len(l))
	copy(other, l)

	return json.Marshal(other)
}

func (l *EmbeddedItems) UnmarshalJSON(d []byte) error {
	var single interface{}
	err := json.Unmarshal(d, &single)
	if err == nil {
		*l = []interface{}{single}
		return nil
	}

	if _, ok := err.(*json.UnmarshalTypeError); !ok {
		return err
	}

	multiple := []interface{}{}
	err = json.Unmarshal(d, &multiple)

	if err == nil {
		*l = multiple
		return nil
	}

	return err
}
