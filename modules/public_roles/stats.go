package public_roles

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"github.com/diraven/sugo"
	"sort"
	"strconv"
)

type sStat struct {
	role  *discordgo.Role
	count int
}

type tStats []sStat

func (ss *tStats) increment(role *discordgo.Role) {
	for i := range *ss {
		if (*ss)[i].role.ID == role.ID {
			(*ss)[i].count = (*ss)[i].count + 1
		}
	}
}

func (ss *tStats) Len() int {
	return len(*ss)
}
func (ss *tStats) Less(i, j int) bool {
	return (*ss)[i].count < (*ss)[j].count
}

func (ss *tStats) Swap(i, j int) {
	(*ss)[i], (*ss)[j] = (*ss)[j], (*ss)[i]
}

var statsCmd = &sugo.Command{
	Trigger:            "stats",
	Description:        "Lists public roles with the highest/lowest count of people.",
	Usage:              "top/bottom",
	PermittedByDefault: true,
	ParamsAllowed: true,
	Execute: func(ctx context.Context, sg *sugo.Instance, c *sugo.Command, m *discordgo.Message, q string) error {
		var err error

		// Get a guild.
		guild, err := sg.GuildFromMessage(m)
		if err != nil {
			_, err = sg.RespondDanger(m, err.Error())
			return err
		}

		// Get all public roles.
		roles, err := publicRoles.findGuildPublicRole(sg, m, "")

		// Make a storage for stats we are about to gather.
		stats := &tStats{}

		// Fill stats with zero values.
		for _, role := range roles {
			*stats = append(*stats, sStat{
				role,
				0,
			})
		}

		// Make members array we will be working with.
		for _, member := range guild.Members {
			for _, roleID := range member.Roles {
				for _, role := range roles {
					// Check if user has role
					if role.ID == roleID {
						stats.increment(role)
					}
				}
			}
		}

		// Sort people.
		sort.Sort(stats)

		// Reverse results if we want bottom side of the chart
		if q != "bottom" {
			sort.Sort(sort.Reverse(stats))
		}

		if len(*stats) > 0 {
			// Start building response.
			var response string
			response = response + "```\n"
			for i, stat := range *stats {
				response = response + strconv.Itoa(i+1) + ". " + stat.role.Name + " (" + strconv.Itoa(stat.count) + ")\n"
				if i > 9 {
					break
				}
			}
			response = response + "```"
			_, err = sg.Respond(m, response)

		} else {
			_, err = sg.RespondDanger(m, "no data available")
		}

		return err
	},
}
