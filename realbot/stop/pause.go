package stop

type Stop struct {
	SinglePauses map[string]*SinglePause
}

type SinglePause struct {
	Name           string
	FromTime       string
	ToTime         string
	ToTimeTimstamp int64
	ActivePause    bool
}

func (stop *Stop) Init(formations ...string) {
	stop.SinglePauses = make(map[string]*SinglePause)
	for _, f := range formations {
		sp := &SinglePause{
			Name:           f,
			FromTime:       "",
			ToTime:         "",
			ToTimeTimstamp: 0,
			ActivePause:    false,
		}
		stop.SinglePauses[f] = sp
	}
}

func (stop *Stop) SetPause(formation string, currentCandleTimestamp int64, minutes int64) {
	stop.SinglePauses[formation].FromTime = TimestampToDate(currentCandleTimestamp)
	stop.SinglePauses[formation].ToTimeTimstamp = currentCandleTimestamp + minutes*60
	stop.SinglePauses[formation].ToTime = TimestampToDate(currentCandleTimestamp + minutes*60)
	stop.SinglePauses[formation].ActivePause = true
}

func (stop *Stop) CancelPause(formation string) bool {
	if stop.SinglePauses[formation].ActivePause {
		stop.SinglePauses[formation].FromTime = ""
		stop.SinglePauses[formation].ToTimeTimstamp = 0
		stop.SinglePauses[formation].ToTime = ""
		stop.SinglePauses[formation].ActivePause = true
		return true
	}
	return false
}
func (stop *Stop) CheckPause(formation string, currentCandleTimestamp int64) bool {
	if stop.SinglePauses[formation].ActivePause {
		if stop.SinglePauses[formation].ToTimeTimstamp > currentCandleTimestamp {
			return true
		} else {
			stop.CancelPause(formation)
			return false
		}
	} else {
		return false
	}
}
