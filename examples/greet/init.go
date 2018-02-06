package greet

import (
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
)

// Init initializes module on the given bot.
func Init(sg *sugo.Instance, message string) {
	if message != "" {
		sg.AddStartupHandler(func(sg *sugo.Instance) error {
			sg.Session.AddHandler(func(s *discordgo.Session, ma *discordgo.GuildMemberAdd) {
				channel, err := s.UserChannelCreate(ma.Member.User.ID)
				if err != nil {
					sg.HandleError(err)
				}
				_, err = s.ChannelMessageSend(channel.ID, message)
				if err != nil {
					sg.HandleError(err)
				}
			})
			return nil
		})
	}
}
