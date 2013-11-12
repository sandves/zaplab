package ztorage

import (
	"github.com/sandves/zaplab/chzap"
)

type SliceZaps []chzap.ChZap

func NewZapStore() *Zaps {
	zs := make(Zaps, 0)
	return &zs
}

func (zs *SliceZaps) StoreZap(z chzap.ChZap) {
	*zs = append(*zs, z)
}

func (zs *SliceZaps) ComputeViewers(chName string) int {
	viewers := 0
	for _, v := range *zs {
		if v.ToChan == chName {
			viewers++
		}
		if v.FromChan == chName {
			viewers--
		}
	}
	return viewers
}

func (zs *SliceZaps) ComputeZaps() int {
	return len(zs)
}

func (zs *SliceZaps) TopTenChannels() []string {
	top10 := make([string]int)

	for  _, v := range *zs {
		if _, ok := top10[v.ToChan]; !ok
				s = s[v.ToChan] = zs.ComputeViewers(v.ToChan))
			}
		}
	}

	return Top10(top10)
}
