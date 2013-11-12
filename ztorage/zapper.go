import (
	"github.com/sandves/zaplab/chzap"
)

type Zapper interface {
   StoreZap(chzap.ChZap)
   TopTenChannels() []string
   ComputeViewers(string) int
   CoputeZaps() int
}