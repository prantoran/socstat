package socstat

import (
	"sync"
	"time"
)

// SocStat defines interface
type SocStat interface {
	Duration(time.Duration)
	IncConn()
	CntConn() int
}

type node struct {
	st  time.Time
	nxt *node
	cnt int
}

// socStat
// d = duration in which to count
// nd = period range of each node
type socStat struct {
	mu   sync.Mutex
	d    time.Duration
	nd   time.Duration
	head *node
	tail *node
	cnt  int
}

func (s *socStat) Duration(dd time.Duration) {
	s.d = dd
	s.nd = s.d / 10
}

func (s *socStat) IncConn() {
	s.mu.Lock()
	now := time.Now()

	if s.tail == nil {
		s.tail = &node{
			st:  now,
			nxt: nil,
			cnt: 0,
		}
	}
	if now.Sub(s.tail.st) > s.nd {
		n := node{
			st:  now,
			nxt: nil,
			cnt: 0,
		}

		s.tail.nxt = &n
		s.tail = &n
	}

	s.tail.cnt++
	s.cnt++

	if s.head == nil {
		s.head = s.tail
	}

	s.mu.Unlock()
}

func (s *socStat) rmExpired() {
	s.mu.Lock()
	now := time.Now()
	for {
		if s.head == nil || now.Sub(s.head.st) <= s.d {
			break
		}
		s.cnt -= s.head.cnt
		s.head = s.head.nxt
	}

	if s.head == nil {
		s.tail = nil
	}
	s.mu.Unlock()
}

func (s *socStat) CntConn() int {
	s.rmExpired()
	return s.cnt
}

// NewSocStat returns a socStat that implements the SocStat interface.
// Duration is set to 10 minutes and duration of each node is set to 1 minute.
func NewSocStat() SocStat {
	s := socStat{
		cnt:  0,
		head: nil,
		tail: nil,
		d:    time.Minute * 10,
		nd:   time.Minute,
	}
	return &s
}
