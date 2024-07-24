package utils

import (
	"fmt"
	"time"

	"github.com/xeonx/timeago"
)

type PlayerDungeonState struct {
	Max   int
	Left  int
	Reset time.Time
}

func (s *PlayerDungeonState) RenderUser() string {
	if s.Reset.Before(time.Now()) && s.Left == 0 {
		return "🕓 - `Dungeons` (soon)"
	}
	if s.Reset.Before(time.Now()) {
		return fmt.Sprintf("✅ - `Dungeons` (**%d/%d**)", s.Left, s.Max)
	}
	return fmt.Sprintf("🕓 - `Dungeons` (%s)", timeago.English.Format(s.Reset))
}
