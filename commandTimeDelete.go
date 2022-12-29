package main

import (
	"database/sql"

	"github.com/bwmarrin/discordgo"
)

func commandTimeDelete(m *discordgo.MessageCreate, args []string) {
	// Check to see if this is a race channel (and get the race from the database).
	var race *Race
	if v, err := getRace(m.ChannelID); err == sql.ErrNoRows {
		discordSend(m.ChannelID, "You can only use that command in a race channel.")
		return
	} else if err != nil {
		msg := "Failed to get the race from the database: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	} else {
		race = v
	}

	// Check to see if this person is one of the two racers.
	if m.Author.ID != race.Racer1.DiscordID && m.Author.ID != race.Racer2.DiscordID {
		discordSend(m.ChannelID, "Only \""+race.Racer1.Username+"\" and \""+race.Racer2.Username+"\" can reschedule this match.")
		return
	}

	// Check to see if this race has already been scheduled.
	if race.State == RaceStateInitial {
		discordSend(m.ChannelID, "There is no need to reschedule until both racers have already agreed to a time.")
		return
	}

	// Set the scheduled time to null.
	if err := modals.Races.UnsetDatetimeScheduled(m.ChannelID); err != nil {
		msg := "Failed to unset the scheduled time: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	// Set the state.
	race.State = RaceStateInitial
	if err := modals.Races.SetState(m.ChannelID, race.State); err != nil {
		msg := "Failed to set the state: " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	discordSend(m.ChannelID, "The currently scheduled time has been deleted. Please suggest a new time with the `!time` command.")
}
