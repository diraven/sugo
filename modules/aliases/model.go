package aliases

import (
	"context"
	"github.com/bwmarrin/discordgo"
	"database/sql"
	"errors"
	"github.com/diraven/sugo"
)

type tAliases map[string]string // map[alias]commandPath

type tAliasesStorage map[string]*tAliases // map[guildID]*tAliases

var aliases = tAliasesStorage{}

func (as *tAliasesStorage) all(g *discordgo.Guild) *tAliases {
	aliases, ok := (*as)[g.ID]
	if ok {
		return aliases
	}
	return &tAliases{}
}

func (as *tAliasesStorage) get(sg *sugo.Instance, g *discordgo.Guild, a string) string {
	commandPath, ok := (*as.all(g))[a]
	if ok {
		return commandPath
	}
	return ""
}

func (as *tAliasesStorage) set(sg *sugo.Instance, g *discordgo.Guild, a string, commandPath string) error {
	_, err := sg.DB.Exec(`
		INSERT OR REPLACE INTO aliases (
			guild_id, alias, command_path
		) VALUES (
			?, ?, ?
		);
	`, g.ID, a, commandPath)

	if err != nil {
		return err
	}

	err = as.reload(sg)
	return err
}

func (as *tAliasesStorage) swap(sg *sugo.Instance, g *discordgo.Guild, a1 string, a2 string) error {
	var err error

	var rowid1 string
	err = sg.DB.QueryRow("SELECT rowid FROM aliases WHERE alias=?;", a1).Scan(&rowid1)
	switch {
	case err == sql.ErrNoRows:
		return errors.New("alias not found: " + a1)
	case err != nil:
		return err
	}

	var rowid2 string
	err = sg.DB.QueryRow("SELECT rowid FROM aliases WHERE alias=?;", a2).Scan(&rowid2)
	switch {
	case err == sql.ErrNoRows:
		return errors.New("alias not found: " + a2)
	case err != nil:
		return err
	}

	tx, err := sg.DB.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE aliases SET rowid=? WHERE alias=?", -1, a2)
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE aliases SET rowid=? WHERE alias=?", rowid2, a1)
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE aliases SET rowid=? WHERE alias=?", rowid1, a2)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	err = as.reload(sg)
	return err
}

func (as *tAliasesStorage) del(sg *sugo.Instance, g *discordgo.Guild, a string) error {
	var err error

	_, err = sg.DB.Exec(`
		DELETE FROM aliases
			WHERE alias=? AND guild_id=?;
	`, a, g.ID)
	if err != nil {
		return err
	}

	err = as.reload(sg)
	return err
}

func (as *tAliasesStorage) reload(sg *sugo.Instance) error {
	// Variable to store errors if any.
	var err error

	// Initialize storage.
	*as = tAliasesStorage{}

	// Get rows from DB.
	rows, err := sg.DB.Query("SELECT guild_id, alias, command_path FROM aliases")
	if err != nil {
		return err
	}

	// Put rows into the in-memory struct.
	for rows.Next() {
		var guildID string
		var alias string
		var commandPath string
		if err := rows.Scan(&guildID, &alias, &commandPath); err != nil {
			return err
		}
		_, ok := (*as)[guildID]
		if !ok {
			(*as)[guildID] = &tAliases{}
		}
		(*(*as)[guildID])[alias] = commandPath
	}

	return err
}