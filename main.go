package main

import (
	"fmt"
	"log"
	"regexp"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

type user struct {
	user   string
	expire int64
}

var users []user

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token:  "YOUR_TOKEN",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal("Error starting Bot")
	}
	b.Handle(tb.OnUserJoined, func(m *tb.Message) {
		var user user
		user.expire = m.Unixtime + 60
		user.user = m.UserJoined.Username
		users = append(users, user)
		fmt.Println(m.UserJoined.Username + " joined " + m.Chat.Title)
	})
	b.Handle(tb.OnText, func(m *tb.Message) {
		i := 0
		if m.Text == "/start" {
			b.Reply(m, "I'm still alive ☺")
			b.Reply(m, "if you want me to ban newcomers who try to send links, add me to a group and make me admin.")
		}
		u := m.Sender.Username
		for i < len(users) {
			if u == users[i].user {
				if containslink(m.Text) {
					member, _ := b.ChatMemberOf(m.Chat, m.Sender)
					b.Ban(m.Chat, member)
					b.Send(m.Chat, u+" just started spamming. They are now banned for 60s.")
					fmt.Println(u + " was banned on " + m.Chat.Title)
					time.Sleep(60 * time.Second)
					b.Unban(m.Chat, m.Sender)
					fmt.Println(u + " was unbanned on " + m.Chat.Title)
				}
			}
			i++
		}
	})
	go usercleaner()
	b.Start()
}

func usercleaner() {
	for true {
		i := 0
		curtime := time.Now().Unix()
		var tmpusr []user
		for i < len(users) {
			if users[i].expire < curtime {
				tmpusr = append(tmpusr, users[i])
			}
			i++
		}
		users = tmpusr
		time.Sleep(60 * time.Second)
	}
}

func containslink(s string) bool {
	b, _ := regexp.MatchString(`(https?)?(\:\/\/)?(www\.)?\w+(\.\w+)`, s)
	return b
}
