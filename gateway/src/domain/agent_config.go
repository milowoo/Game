package domain

import (
	"time"
)

type AgentConfig struct {
	EnableLogSend          bool
	EnableLogRecv          bool
	EnableCheckPing        bool
	EnableCheckLoginParams bool

	EnableCachedMsg   bool
	EnableCompressMsg bool

	CachedMsgMaxCount        int
	CompressMsgSizeThreshold int

	TimestampExpireDuration time.Duration
}
