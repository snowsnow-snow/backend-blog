package utility

import (
	"log"
	"sync"

	"github.com/bwmarrin/snowflake"
)

var (
	node *snowflake.Node
	once sync.Once
)

// InitSnowflake initializes the snowflake node.
// nodeID should be unique per instance (0-1023)
func InitSnowflake(nodeID int64) {
	once.Do(func() {
		var err error
		node, err = snowflake.NewNode(nodeID)
		if err != nil {
			log.Fatalf("Failed to initialize snowflake node: %v", err)
		}
	})
}

// GenID generates a new snowflake ID (int64)
func GenID() int64 {
	if node == nil {
		// Default initialization if not called explicitly, though Init is better
		InitSnowflake(1)
	}
	return node.Generate().Int64()
}

// GenIDString generates a new snowflake ID as string
func GenIDString() string {
	if node == nil {
		InitSnowflake(1)
	}
	return node.Generate().String()
}
