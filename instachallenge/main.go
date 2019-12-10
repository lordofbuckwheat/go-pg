package main

import (
	"fmt"
	"github.com/ahmdrz/goinsta/v2"
	"log"
)

func main() {
	//username := flag.String("session", "session.json", "path to config file")
	//flag.Parse()
	//session := fmt.Sprintf("sessions/instagram/%s.json", *username)
	//_, err := os.Stat(session)
	//if err != nil {
	// panic(err)
	//}
	insta := goinsta.New("palemaltd", "Palema19.")
	err := insta.Login()
	if err != nil {
		switch v := err.(type) {
		case goinsta.ChallengeError:
			fmt.Println("err", err)
			err := insta.Challenge.Process(v.Challenge.APIPath)
			if err != nil {
				log.Fatalln(err)
			}
			fmt.Println("123", v, err)
		}
		panic(err)
	}
	user, err := insta.Profiles.ByName("eurohoops_official")
	if err != nil {
		panic(err)
	}
	fmt.Println("before feeds")
	media := user.Feed()
	fmt.Println("after feeds")
	var i = 0
	for media.Next() && i < 100 {
		fmt.Println("item", len(media.Items), media.Error())
		i++
	}
}