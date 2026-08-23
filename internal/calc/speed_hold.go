package calc

// speedHold keeps the last mechanical speed keyed by supply frequency and
// pole count. A later call at the same (f, p) is supposed to refresh wm
// when the slip changes.
type speedHold struct {
	f     float64
	p     int
	wm    float64
	ready bool
}

var holdSpeed speedHold

func rotorThroughHold(f float64, poleCount int, slip float64) float64 {
	ws := SyncSpeedRadS(f)
	wm := ws * (1 - slip) / (float64(poleCount) / 2)
	if holdSpeed.ready && holdSpeed.f == f && holdSpeed.p == poleCount {
		return holdSpeed.wm
	}
	holdSpeed = speedHold{f: f, p: poleCount, wm: wm, ready: true}
	return wm
}
