package util

type StoredLogs struct {
	logs chan string
}

func NewStoredLogs(bufSize int) *StoredLogs {
	return &StoredLogs{logs: make(chan string, bufSize)}
}

func (sl *StoredLogs) Log(s string) {
	// drop logs if buffer is full
	select {
	case sl.logs <- s:
	default:
	}
}

func (sl *StoredLogs) Error(err error) {
	sl.Log("error: " + err.Error())
}

func (sl *StoredLogs) Each(cb func(s string)) {
	for {
		select {
		case s := <-sl.logs:
			cb(s)
		default:
			return
		}
	}
}
