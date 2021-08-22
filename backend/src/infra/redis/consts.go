package redis

import (
	"fmt"
	"os"
)

const (
	EVENT_SET     = "set"
	EVENT_DEL     = "del"
	EVENT_EXPIRED = "expired"
	EVENT_EXPIRE  = "expire"
)

var (
	KEYSPACE = fmt.Sprintf("__keyspace@%s__", os.Getenv("REDIS_DB"))
	KEYEVENT = fmt.Sprintf("__keyevent@%s__", os.Getenv("REDIS_DB"))
)
