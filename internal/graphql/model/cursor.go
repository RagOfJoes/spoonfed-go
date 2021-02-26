package model

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const (
	cursorPrefix = "cursor"
)

// EncodeCursor uses base64 to encode a Node's field
func EncodeCursor(v interface{}) string {
	switch v.(type) {
	case int:
		intCursor := v.(int)
		return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s%d", cursorPrefix, intCursor+1)))
	case string:
		stringCursor := v.(string)
		return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s%s", cursorPrefix, stringCursor)))
	default:
		return ""
	}
}

// DecodeCursor into raw value
// TODO: add support for multiple types
func DecodeCursor(s string) (interface{}, error) {
	decoded, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	prefixed := strings.TrimPrefix(string(decoded), cursorPrefix)
	if err != nil {
		return nil, err
	}
	toInt, err := strconv.ParseInt(prefixed, 10, 64)
	if err == nil {
		unix := time.Unix(toInt, 0)
		return &unix, nil
	}
	return &prefixed, nil
}

// CursorToBson decodes a cursor and returns a valid mongo
// interface for filtering/sorting
func CursorToBson(cursor *string, sortKey string, sortOrder int) bson.E {
	if cursor == nil {
		return bson.E{}
	}

	decoded, err := DecodeCursor(*cursor)
	if err != nil {
		return bson.E{}
	}

	var filter bson.E
	switch decoded.(type) {
	case *time.Time:
		op := "$lt"
		if sortOrder != -1 {
			op = "$gt"
		}
		decodedTime := decoded.(*time.Time)
		filter = bson.E{
			Key: sortKey, Value: bson.D{
				{
					Key: op, Value: decodedTime.UTC(),
				},
			},
		}
		break
	case *string:
		op := "$lte"
		if sortOrder != -1 {
			op = "$gte"
		}
		decodedString := decoded.(*string)
		filter = bson.E{
			Key: sortKey, Value: bson.D{
				{
					Key: op, Value: decodedString,
				},
			},
		}
		break
	default:
		filter = bson.E{}
	}
	return filter
}
