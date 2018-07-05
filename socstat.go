package socstat

import (
	"sync"
	"time"
)

type SocStat interface {
	Duration(d time.Duration)
	IncConn()
	CntConn() int
}

type node struct {
	st  time.Time
	nxt *node
}

type socStat struct {
	mu   sync.Mutex
	d    time.Duration
	head *node
	tail *node
	cnt  int
}

func (s *socStat) Duration(dd time.Duration) {
	s.d = dd
}

func (s *socStat) IncConn() {
	s.mu.Lock()
	n := node{
		st:  time.Now(),
		nxt: nil,
	}
	s.cnt++

	if s.tail != nil {
		s.tail.nxt = &n
	}

	if s.head == nil {
		s.head = &n
	}

	s.tail = &n
	s.mu.Unlock()
}

func (s *socStat) rmExpired() {
	s.mu.Lock()
	now := time.Now()
	for {
		if s.head == nil || now.Sub(s.head.st) <= s.d {
			break
		}
		s.head = s.head.nxt
		s.cnt--
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
// Inititally duration is set to 5 minutes.
func NewSocStat() SocStat {
	s := socStat{
		cnt:  0,
		head: nil,
		tail: nil,
		d:    time.Minute * 10,
	}
	return &s
}
