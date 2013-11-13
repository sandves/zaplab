package ztorage

import (
	"github.com/sandves/zaplab/chzap"
	"time"
)

type Zaps map[string]int
type Stats map[string]chzap.ChZap

var totalZapDuration time.Duration
var totalZapEvents int
var stats Stats

func NewZapStore() Zaps {
	stats = make(Stats)
	zs := make(Zaps)
	return zs
}

func (zs Zaps) StoreZap(z chzap.ChZap) {

	totalZapEvents++
	if _, ok := stats[(z.IP).String()]; ok {
		dur := z.Duration(stats[(z.IP).String()])
		totalZapDuration += dur
	}
	stats[(z.IP).String()] = z

	/*If the channel doesn't exist in the map,
	put the key (channelname) in the map and
	assign its value (number of viewers) to zero*/
	if _, ok := zs[z.ToChan]; !ok {
		zs[z.ToChan] = 0
	}
	if _, ok := zs[z.FromChan]; !ok {
		zs[z.FromChan] = 0
	}

	/*if a viewer zaps to a channel, increment the
	number of viewers by one.
	if a viewer leaves a channel, decrement the number
	of viewers by one*/
	for key, _ := range zs {
		if z.ToChan == key {
			zs[key]++
		}
		if z.FromChan == key {
			zs[key]--
		}
	}
}

func (zs Zaps) ComputeViewers(chName string) int {
	viewers := 0
	for k, v := range zs {
		if k == chName {
			viewers = v
		}
	}
	return viewers
}

func (zs Zaps) ComputeZaps() int {
	zaps := 0
	for _, v := range zs {
		zaps += v
	}
	return zaps
}

func (zs Zaps) AverageZapDuration() time.Duration {
	if totalZapEvents != 0 {
		return (totalZapDuration) / (time.Duration(totalZapEvents))
	} else {
		return time.Duration(0)
	}
}

func (zs Zaps) TopTenChannels() []string {
	return Top10(zs)
}
