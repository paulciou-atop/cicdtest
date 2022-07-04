/*
  This file consist of inventory timestamp related functions
*/
package inventory

import (
	"nms/lib/pgutils"
	"time"
)

func pipe(funs ...func(Inventory) Inventory) func(Inventory) Inventory {
	return func(inv Inventory) Inventory {
		for _, f := range funs {
			inv = f(inv)
		}
		return inv
	}
}

func now() string {
	return time.Now().Format(pgutils.TIME_LAYOUT)
}

// checkInvRecoveringTS, check incoming inventory is recovering, update timestamp
func checkInvRecoveringTS(inv Inventory) Inventory {
	ret := inv
	lastMissing, err := time.Parse(pgutils.TIME_LAYOUT, ret.LastMissing)
	if err != nil {
		// no missing ts, inventory wasn't recovering
		return ret
	}

	lastRecover, err := time.Parse(pgutils.TIME_LAYOUT, ret.LastRecovered)
	if err != nil {
		// no recovered ts, assign the zero time
		lastRecover = time.Time{}
	}
	if lastRecover.Before(lastMissing) {
		// lastMissing > lastRecorvered and inventory show againg, so it's recorver
		ret.LastRecovered = now()
	}
	return ret
}

func missingInvTS(inv Inventory) Inventory {
	ret := inv
	ret.LastMissing = now()
	return ret
}

func newInvTS(inv Inventory) Inventory {
	ret := inv
	ret.CreatedAt = now()
	ret.LastSeen = now()
	return ret
}

func updateInvTS(inv Inventory) Inventory {
	ret := inv
	ret.LastSeen = now()
	return ret
}
