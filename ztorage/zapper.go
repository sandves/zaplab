package ztorage

import (
	"github.com/zaplab/chzap"
)

type Zapper interface {
	StoreZap(chzap.ChZap)
	TopTenChannels() []string
	ComputeViewers(string) int
	ComputeZaps() int
}
