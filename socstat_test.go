package socstat

import (
	"testing"
	"time"
)

func TestCntConnEmpty(t *testing.T) {
	s := NewSocStat()

	cnt := s.CntConn()
	if cnt != 0 {
		t.Errorf("cnt not zero when nothing added")
	}
}

func TestCntConnOneEntry(t *testing.T) {
	s := NewSocStat()
	s.IncConn()
	cnt := s.CntConn()
	if cnt != 1 {
		t.Errorf("cnt not one after one increment")
	}
}

func TestCntConnTwoEntry(t *testing.T) {
	s := NewSocStat()
	s.IncConn()
	s.IncConn()
	cnt := s.CntConn()
	if cnt != 2 {
		t.Errorf("cnt not one after one increment")
	}
}

func TestCntConnTwoEntryTwoNodes(t *testing.T) {
	s := NewSocStat()
	s.Duration(10 * time.Second)
	s.IncConn()
	time.Sleep(time.Second)
	s.IncConn()
	cnt := s.CntConn()
	if cnt != 2 {
		t.Errorf("cnt not one after one increment")
	}
}

func TestCntConnAfterDurationOneValid(t *testing.T) {
	s := NewSocStat()
	s.Duration(2 * time.Second)
	s.IncConn()
	time.Sleep(2 * time.Second)
	s.IncConn()
	cnt := s.CntConn()
	if cnt == 2 {
		t.Errorf("first inc not removed after duration passed")
	} else if cnt == 0 {
		t.Errorf("both inc removed")
	} else if cnt != 1 {
		t.Errorf("cnt not valid")
	}
}

func TestCntConnAfterDurationNoneValid(t *testing.T) {
	s := NewSocStat()
	s.Duration(2 * time.Second)
	s.IncConn()
	s.IncConn()
	time.Sleep(2 * time.Second)

	cnt := s.CntConn()
	if cnt == 2 {
		t.Errorf("nothing removed after duration passed")
	} else if cnt != 0 {
		t.Errorf("cnt not valid")
	}
}

func TestCntConnAfterDurationNoneValidThenOneInc(t *testing.T) {
	s := NewSocStat()
	s.Duration(2 * time.Second)
	s.IncConn()
	s.IncConn()
	time.Sleep(2 * time.Second)
	s.IncConn()
	cnt := s.CntConn()

	switch cnt {
	case 3:
		t.Errorf("none inc removed")
	case 2:
		t.Errorf("one inc not removed")
	case 0:
		t.Errorf("last inc not counted")
	}
}

func Test12IncAt1SecPause(t *testing.T) {
	s := NewSocStat()
	s.Duration(10 * time.Second)
	for i := 0; i < 12; i++ {
		time.Sleep(time.Second)
		s.IncConn()
	}

	cnt := s.CntConn()

	if cnt != 10 {
		t.Errorf("cnt not 10, cnt: %v", cnt)
	}
}
