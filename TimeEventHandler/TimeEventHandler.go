package TimeEventHandler

import "time"

const (
	eveningTimeAsNumber = 1140 //19:00
)

func TimeToNumber() int16 {
	return int16(time.Now().Hour()*60 + time.Now().Minute())
}

func IsEveningTime() bool {
	if TimeToNumber() <= eveningTimeAsNumber {
		return false
	}
return true
}
