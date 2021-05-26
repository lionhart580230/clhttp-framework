package skylog

import "time"

func GetTimeFormat(stamp uint32, fmt string) (string, error) {
	var utc time.Time
	if stamp == 0 {
		utc = time.Now().UTC()
	} else {
		utc = time.Unix(int64(stamp), 0).UTC()
	}

	//ntime := utc.In(loc)
	return utc.Add(8 * time.Hour).Format(fmt), nil
}
