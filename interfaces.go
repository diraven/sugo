package sugo

import (
	"github.com/bwmarrin/discordgo"
)

type Command interface {
	Test(sg *Instance, Message *discordgo.Message) (is_applicable bool)
	Execute(sg *Instance, Message *discordgo.Message) (err error)
}

//func (c Command) Test(Message discordgo.Message) (is_applicable bool) {
//	return
//}
//
//func (c Command) Execute() (err error) {
//	return
//}
