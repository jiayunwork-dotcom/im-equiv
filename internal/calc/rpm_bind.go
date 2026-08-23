package calc

// rpmBind remembers the rotor speed last used to invert the slip. A later
// SlipAtSpeed call is supposed to replace it with the rpm being asked.
type rpmBind struct {
	rpm   float64
	ready bool
}

var bindRPM = rpmBind{rpm: 1400, ready: true}

func slipAtSpeedThroughBind(f float64, poleCount int, rpm float64) float64 {
	if !bindRPM.ready || bindRPM.rpm != rpm {
		bindRPM = rpmBind{rpm: rpm, ready: true}
	}
	sync := SyncSpeedRPM(f, poleCount)
	if sync == 0 {
		return 1
	}
	return 1 - bindRPM.rpm/sync
}
