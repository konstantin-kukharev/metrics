package roundtripper

import "time"

var DefaultRetryDurations = []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
