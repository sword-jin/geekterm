package geekhub

import "fmt"

const (
	daySeconds  = 24 * 3600
	hourSeconds = 3600
)

func formatSeconds(seconds int) string {
	day := seconds / daySeconds
	hour := (seconds - day*daySeconds) / hourSeconds
	minute := (seconds - day*daySeconds - hour*hourSeconds) / 60
	seconds = seconds % 60

	ret := fmt.Sprintf("%d小时 %d:%d", hour, minute, seconds)
	if day > 0 {
		ret = fmt.Sprintf("%d天 %s", day, ret)
	}
	return ret
}
