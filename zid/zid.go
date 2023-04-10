package zid

import (
	"github.com/bwmarrin/snowflake"
	"github.com/google/uuid"
	"strings"
	"time"
)

// NewUuid
// @Auth: oak  2021-10-15 18:36:29
// @Description:  UUIDv4
// @return string
func NewUuid() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

var node *snowflake.Node

func InitSnowflake(machineID int64) {
	var err error
	snowflake.Epoch = time.Now().UnixMilli()
	node, err = snowflake.NewNode(machineID)
	if err != nil {
		panic(err)
	}
}

// NewSnowflake
// @Auth: oak  2021-10-15 18:36:13
// @Description:  雪花ID
// @return string
func NewSnowflake() int64 {
	if node == nil {
		panic("Please initialize snowflake first")
	}
	return node.Generate().Int64()
}
