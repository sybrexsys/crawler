package manager

import (
	"time"
)

type limiter struct {
	counts []time.Time
	header int
	footer int
	size   int
	period time.Duration
	empty  bool
}

func newLimiter(size int, period time.Duration) *limiter {
	return &limiter{
		counts: make([]time.Time, size),
		size:   size,
		period: period,
		empty:  true,
	}
}

func (l *limiter) getSize() int {
	if l.empty {
		return 0
	}
	if l.header == l.footer {
		return l.size
	}
	if l.header > l.footer {
		return l.header - l.footer
	}
	return l.header + l.size - l.footer
}

func (l *limiter) check() time.Duration {
	if l.header != l.footer {
		return 0
	}
	k := l.header
	if l.header <= l.footer {
		k += len(l.counts)
	}
	currentTime := time.Now()
	var i int
	for i = l.footer; i < k; i++ {
		sub := currentTime.Sub(l.counts[i%l.size])
		if sub < l.period {
			if i == l.footer {
				return l.period - sub
			}
			break
		}
		l.footer++
		if l.footer == l.size {
			l.footer = 0
		}
	}
	if i == k {
		l.empty = true
	}
	return time.Duration(0)
}

func (l *limiter) addNewTime() bool {
	if l.getSize() == l.size {
		return false
	}
	l.counts[l.header] = time.Now()
	l.header++
	if l.header == l.size {
		l.header = 0
	}
	l.empty = false
	return true
}
