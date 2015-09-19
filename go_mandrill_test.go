package mandrill

import (
	"log"
	"testing"
)

func TestUsersInfo(t *testing.T) {
	u := UsersAPI{}
	info, err := u.GetInfo("NjixlbCzdB14TazGCnYyEQ")
	if err != nil {
		log.Println("UsersAPI GetInfo Error")
		log.Print(err)
	}
	log.Println("UsersAPI GetInfo Successful")
	log.Printf("UsersAPI GetInfo Results: %#v\n\n", info)
}

func TestUsersPing(t *testing.T) {
	u := UsersAPI{}
	ping, err := u.Ping("NjixlbCzdB14TazGCnYyEQ")
	if err != nil {
		log.Println("UsersAPI Ping Error")
		log.Print(err)
	}
	log.Println("UsersAPI Ping Successful")
	log.Printf("UsersAPI Ping Results: %#v\n\n", ping)
}
