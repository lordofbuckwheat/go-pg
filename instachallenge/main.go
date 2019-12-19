package main

import (
	"fmt"
	"github.com/lordofbuckwheat/goinsta/v2"
)

func main() {
	insta := goinsta.New("tvbittest5", "password1234")
	err := insta.Login()
	if err != nil {
		fmt.Println("err", err)
		switch v := err.(type) {
		case goinsta.ChallengeError:
			fmt.Println("challenge error", v)
			err := insta.Challenge.Process(v.Challenge.APIPath)
			if err != nil {
				panic(err)
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
	fmt.Println("12")
	fmt.Println("3")
}