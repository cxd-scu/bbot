package utils

import "fmt"

func GenerateSessionID(agentID int64) string {
	// 这里可以根据实际需求生成唯一的 SessionID，例如使用 UUID
	return "unique-session-id-" + fmt.Sprint(agentID)
}
