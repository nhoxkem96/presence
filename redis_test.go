package presence

import (
	"errors"
	"fmt"
	"testing"
	"time"
)

func initRedisence() (*Session, error) {
	backend, err := NewRedis("192.168.59.103:6381", 10, time.Second*1)
	// backend, err := NewRedis("localhost:6379", 10, time.Second*1)
	if err != nil {
		return nil, err
	}

	ses, err := New(backend)
	if err != nil {
		return nil, err
	}

	return ses, nil
}

func withConn(f func(s *Session)) error {
	s, err := initRedisence()
	if err != nil {
		return err
	}

	f(s)

	return s.Close()
}

func TestSinglePing(t *testing.T) {
	err := withConn(func(s *Session) {
		if err := s.Online("id"); err != nil {
			t.Fatal(err)
		}
	})

	if err != nil {
		t.Fatal(err)
	}

}

func TestMultiPing(t *testing.T) {
	err := withConn(func(s *Session) {
		if err := s.Online("id", "id2"); err != nil {
			t.Fatal(err)
		}
	})

	if err != nil {
		t.Fatal(err)
	}

}

func TestOnlineStatus(t *testing.T) {
	err := withConn(func(s *Session) {
		id := "id3"
		if err := s.Online(id); err != nil {
			t.Fatal(err)
		}

		status, err := s.Status(id)
		if err != nil {
			t.Fatal(err)
		}

		res := status[0]

		if res.Status != Online {
			t.Fatalf("%s should be %s, but it is %s", res.ID, Online, res.Status)
		}
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestOfflineStatus(t *testing.T) {
	err := withConn(func(s *Session) {

		id := "id4"
		if err := s.Online(id); err != nil {
			t.Fatal(err)
		}

		status, err := s.Status("id5")
		if err != nil {
			t.Fatal(err)
		}

		res := status[0]
		if res.Status != Offline {
			t.Fatalf("%s should be %s, but it is %s", res.ID, Offline, res.Status)
		}
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestMultiStatusAllOnline(t *testing.T) {
	err := withConn(func(s *Session) {
		if err := s.Online("id6", "id7"); err != nil {
			t.Fatal(err)
		}

		status, err := s.Status([]string{"id6", "id7"}...)
		if err != nil {
			t.Fatal(err)
		}
		for _, res := range status {
			if res.Status != Online {
				t.Fatalf("%s should be %s, but it is %s", res.ID, Online, res.Status)
			}
		}
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestMultiStatusAllOffline(t *testing.T) {
	err := withConn(func(s *Session) {
		if err := s.Online("id8", "id9"); err != nil {
			t.Fatal(err)
		}

		status, err := s.Status([]string{"id10", "id11"}...)
		if err != nil {
			t.Fatal(err)
		}

		for _, res := range status {
			if res.Status != Offline {
				t.Fatalf("%s should be %s, but it is %s", res.ID, Offline, res.Status)
			}
		}
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestStatusWithTimeout(t *testing.T) {
	err := withConn(func(s *Session) {
		id := "12"
		if err := s.Online(id); err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Second * 2)
		status, err := s.Status(id)
		if err != nil {
			t.Fatal(err)
		}

		res := status[0]
		if res.Status == Online {
			t.Fatalf("%s should be %s, but it is %s", res.ID, Online, res.Status)
		}
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestSubscriptions(t *testing.T) {
	t.Skip("Skipped to travis")
	s, err := initRedisence()
	if err != nil {
		t.Fatal(err)
	}

	// wait for all keys to expire
	time.Sleep(time.Second * 1)

	id1 := "13"
	id2 := "14"
	id3 := "15"

	time.AfterFunc(time.Second*5, func() {
		err := s.Close()
		if err != nil {
			t.Fatal(err)
		}
	})

	time.AfterFunc(time.Second*1, func() {
		err := s.Online(id1, id2, id3)
		if err != nil {
			t.Fatal(err)
		}
		// err = s.Offline(id1, id2, id3)
		// if err != nil {
		//  t.Fatal(err)
		// }
	})

	onlineCount := 0
	offlineCount := 0
	for event := range s.ListenStatusChanges() {
		switch event.Status {
		case Online:
			onlineCount++
		case Offline:
			offlineCount++
		}
	}

	if onlineCount != 3 {
		t.Fatal(
			fmt.Errorf("online count should be 3 it is %d", onlineCount),
		)
	}

	if offlineCount != 3 {
		t.Fatal(
			fmt.Errorf("offline count should be 3 it is %d", offlineCount),
		)
	}
}

func TestJustMultiOffline(t *testing.T) {
	err := withConn(func(s *Session) {
		if err := s.Offline("id16", "id17"); err != nil {
			t.Fatal(err)
		}
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestMultiOnlineAndOfflineTogether(t *testing.T) {
	err := withConn(func(s *Session) {
		if err := s.Online("id18", "id19"); err != nil {
			t.Fatal(err)
		}
		if err := s.Offline("id18", "id19"); err != nil {
			t.Fatal(err)
		}
	})

	if err != nil {
		t.Fatal(err)
	}
}

func TestMultiOfflineWithMultiStatus(t *testing.T) {
	err := withConn(func(s *Session) {
		if err := s.Online("id20", "id21"); err != nil {
			t.Fatal(err)
		}
		if err := s.Offline("id20", "id21"); err != nil {
			t.Fatal(err)
		}
		status, err := s.Status([]string{"id20", "id21"}...)
		if err != nil {
			t.Fatal(err)
		}

		for _, st := range status {
			if st.Status != Offline {
				t.Fatal(errors.New("user should be offline"))
			}
		}
	})

	if err != nil {
		t.Fatal(err)
	}
}
