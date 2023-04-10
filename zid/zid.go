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

// InitSnowflake
// @Description:
// @param date 2006-01-02
// @param machineID
func InitSnowflake(date string, machineID int64) {
	var st time.Time
	st, _ = time.Parse("2006-01-02", date)
	snowflake.Epoch = st.UnixNano() / 1000000
	node, _ = snowflake.NewNode(machineID)
}

// NewSnowflake
// @Auth: oak  2021-10-15 18:36:13
// @Description:  雪花ID
// @return string
func NewSnowflake() string {
	if node == nil {
		panic("Please initialize snowflake first")
	}
	return node.Generate().String()
}
