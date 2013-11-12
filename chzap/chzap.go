package chzap

import (
	"fmt"
	"net"
	"strings"
	"time"
)

const timeLayout = "2006/01/02, 15:04:05"

type ChZap struct {
	Date     time.Time
	IP       net.IP
	FromChan string
	ToChan   string
}

func NewChZap(chzap string) *ChZap {
	chzapSlice := strings.Split(chzap, ", ")

	//The Zap Event consists of 5 fields
	if len(chzapSlice) == 5 {
		unParsedTime := strings.Join(chzapSlice[0:2], ", ")
		date, err := time.Parse(timeLayout, unParsedTime)

		if err != nil {
			fmt.Println(err)
		}

		ip := net.ParseIP(chzapSlice[2])
		fromCh := chzapSlice[3]
		toCh := chzapSlice[4]

		return &ChZap{date, ip, fromCh, toCh}
	} else {
		/*Return a zero initialized ChZap if
		the given string is a StatusChange
		and not a Zap Event*/
		return new(ChZap)
	}
}

func (ch *ChZap) String() string {
	s := []string{
		ch.Date.Format(timeLayout),
		ch.IP.String(),
		ch.FromChan,
		ch.ToChan,
	}
	return strings.Join(s, ", ")
}

func (ch *ChZap) Duration(provided ChZap) time.Duration {
	duration := ch.Date.Sub(provided.Date)
	if duration < 0 {
		return provided.Date.Sub(ch.Date)
	}
	return duration
}
