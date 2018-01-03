package sugo

import "github.com/bwmarrin/discordgo"

// ChannelFromMessage returns a *discordgo.Channel struct from given *discordgo.Message struct.
func (sg *Instance) ChannelFromMessage(m *discordgo.Message) (*discordgo.Channel, error) {
	return sg.State.Channel(m.ChannelID)
}

// GuildFromMessage returns a *discordgo.Guild struct from given *discordgo.Message struct.
func (sg *Instance) GuildFromMessage(m *discordgo.Message) (*discordgo.Guild, error) {
	c, err := sg.ChannelFromMessage(m)
	if err != nil {
		return nil, err
	}
	return sg.State.Guild(c.GuildID)
}

// MemberFromMessage returns a *discordgo.Member struct from given *discordgo.Message struct.
func (sg *Instance) MemberFromMessage(m *discordgo.Message) (*discordgo.Member, error) {
	g, err := sg.GuildFromMessage(m)
	if err != nil {
		return nil, err
	}
	return sg.State.Member(g.ID, m.Author.ID)
}
