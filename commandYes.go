package main

import (
	"database/sql"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

func commandYes(m *discordgo.MessageCreate, args []string) {
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
	var racerNum int
	if m.Author.ID == race.Racer1.DiscordID {
		racerNum = 1
	} else if m.Author.ID == race.Racer2.DiscordID {
		racerNum = 2
	} else {
		discordSend(m.ChannelID, "Only `"+race.Racer1.Username+"` and `"+race.Racer2.Username+"` can veto a build.")
		return
	}

	// Check to see if this race is in the vetoing phase.
	if race.State != RaceStateVetoCharacters && race.State != RaceStateVetoBuilds {
		discordSend(m.ChannelID, "You can only veto something once the characters have been chosen.")
		return
	}

	// Check to see if it is their turn.
	if race.ActiveRacer != racerNum {
		discordSend(m.ChannelID, "It is not your turn.")
		return
	}

	// Check to see if they are out of vetos.
	if (racerNum == 1 && race.Racer1Vetos == 0) ||
		(racerNum == 2 && race.Racer2Vetos == 0) {

		discordSend(m.ChannelID, "You have already used all of your vetos for the match.")
		return
	}

	// The character/build was already added, so remove it.
	var veto string
	if race.State == RaceStateVetoCharacters {
		veto = race.Characters[len(race.Characters)-1]
		race.Characters = race.Characters[:len(race.Characters)-1] // Delete the last element.
		if err := modals.Races.SetCharacters(race.ChannelID, race.Characters); err != nil {
			msg := "Failed to set the characters for race \"" + race.Name() + "\": " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		}
	} else if race.State == RaceStateVetoBuilds {
		veto = race.Builds[len(race.Builds)-1]
		race.Builds = race.Builds[:len(race.Builds)-1] // Delete the last element.
		if err := modals.Races.SetBuilds(race.ChannelID, race.Builds); err != nil {
			msg := "Failed to set the builds for race \"" + race.Name() + "\": " + err.Error()
			log.Error(msg)
			discordSend(m.ChannelID, msg)
			return
		}
	}

	// Decrement the vetos.
	var vetosLeft int
	if racerNum == 1 {
		race.Racer1Vetos--
		vetosLeft = race.Racer1Vetos
	} else if racerNum == 2 {
		race.Racer2Vetos--
		vetosLeft = race.Racer2Vetos
	}
	if err := modals.Races.SetVetos(race.ChannelID, racerNum, vetosLeft); err != nil {
		msg := "Failed to set the vetos for racer " + strconv.Itoa(racerNum) + " on race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	// Set the number of people who have voted on this build.
	race.NumVoted = 2 // If this person is vetoing, then the other person does not get a say.
	if err := modals.Races.SetNumVoted(race.ChannelID, race.NumVoted); err != nil {
		msg := "Failed to set the NumVoted for race \"" + race.Name() + "\": " + err.Error()
		log.Error(msg)
		discordSend(m.ChannelID, msg)
		return
	}

	incrementActiveRacer(race)
	msg := m.Author.Mention() + " vetoed: *" + veto + "*\n\n"
	if race.State == RaceStateVetoCharacters {
		charactersRound(race, msg)
	} else if race.State == RaceStateVetoBuilds {
		buildsRound(race, msg)
	}
}
