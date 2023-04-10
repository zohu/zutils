package zid

import "testing"

func TestNewSnowflake(t *testing.T) {
	InitSnowflake("2023-04-10", 1)
	t.Logf("snowflake -> %s", NewSnowflake())
}
