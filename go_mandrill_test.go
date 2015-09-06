package mandrill

import (
	"log"
	"testing"
)

func TestUsersInfo(t *testing.T) {
	u := User{}
	if err := u.GetInfo("EiTvzstNtLqKKoVYIQjg"); err.Code != 0 {
		t.Errorf("Error Occured: %v\n", err)
	}
	log.Printf("%#v\n", u)
}
