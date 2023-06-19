package main

import (
	"github.com/bwmarrin/discordgo"
)

var (
	commandHandlerMap = make(map[string]func(*discordgo.MessageCreate, []string))
)

func commandHelpGetMsg() string {
	msg := "General commands (all channels):\n"
	msg += "```\n"
	msg += "Command                  Description\n"
	msg += "-----------------------------------------------------------------------\n"
	msg += "!help                    Get a list of all of the commands\n"
	msg += "!r+                      Get info about the Racing+ mod\n"
	msg += "!bracket                 Get the link to the bracket\n"
	msg += "!timezone                Get your stored timezone\n"
	msg += "!timezone [timezone]     Set your stored timezone\n"
	msg += "!gettimezone [username]  Get the timezone of the specified person\n"
	msg += "!stream                  Get your stored stream URL\n"
	msg += "!stream [url]            Set your stored stream URL\n"
	msg += "!getstream [username]    Get the stream of the specified person\n"
	msg += "!random [min] [max]      Get a random number\n"
	msg += "!randomchar              Get a random character\n"
	msg += "!randombuild             Get a random build\n"
	msg += "!getnext                 Get the time of the next scheduled match\n"
	msg += "!schedule                Get a list of all of the currently scheduled matches\n"
	msg += "```\n"
	msg += "Match commands (in a match channel):\n"
	msg += "```\n"
	msg += "Command                  Description\n"
	msg += "-----------------------------------------------------------------------\n"
	msg += "!time                    Get the currently scheduled time for the match\n"
	msg += "!time [date & time]      Suggest a time for the match to your opponent\n"
	msg += "!timeok                  Confirm that the suggested time is good\n"
	msg += "!timedelete              Delete the currently scheduled time\n"
	msg += "!cast                    Volunteer to be the caster for the match\n"
	msg += "!castcancel              Unvolunteer to be the caster\n"
	msg += "!caster                  Get the person who volunteered to cast\n"
	msg += "!casterok                Confirm that you are okay with the caster\n"
	msg += "!casternotok             Reject the current caster\n"
	msg += "!ban [num]               Ban something\n"
	msg += "!pick [num]              Pick something\n"
	msg += "!yes                     Veto the selected thing\n"
	msg += "!no                      Do not veto the selected thing\n"
	msg += "!score                   Report the score after the match has completed\n"
	msg += "                         (with your number first)\n"
	msg += "```"
	/*
		msg += "Admin-only commands:\n"
		msg += "```\n"
		msg += "Command               Description\n"
		msg += "-------------------------------------------------------------------------\n"
		msg += "!settimezone             Set a user's timezone for them\n"
		msg += "!setstream               Set a user's stream for them\n"
		msg += "!setcasterok             Give permission on behalf of a racer\n"
		msg += "!setcasternotok          Deny permission on behalf of a racer\n"
		msg += "!setcasteralwaysok       Enable a user's default caster approval for them\n"
		msg += "!setcasteralwaysnotok    Disable a user's default caster approval for them\n"
		msg += "!startround              Start the current round of the tournament\n"
		msg += "!endround                Delete all of the channels for this round\n"
		msg += "!checkround              Do a dry run of "!startround"\n"
		msg += "!forcetime               Force a scheduled time\n"
		msg += "!forcetimeok             Force the scheduled time to be ok\n"
		msg += "!forcetimedelete         Force the currently scheduled time to be deleted\n"
		msg += "!forceban [num]          Force the current racer to ban\n"
		msg += "!forcepick [num]         Force the current racer to pick\n"
		msg += "!forceyes                Force the current racer to veto\n"
		msg += "!forceno                 Force the current racer to not veto\n"
		msg += "!join                    Print out the URL to join another server\n"
		msg += "!getstate                Get the current state of the match\n"
		msg += "!getchannelid [name]     Get the ID of the specified Discord channel\n"
		msg += "!debug                   Execute the debug function\n"
		msg += "```"
	*/

	return msg
}

func commandInit() {
	// General commands
	commandHandlerMap["help"] = commandHelp
	commandHandlerMap["commands"] = commandHelp
	commandHandlerMap["r+"] = commandRacingPlus
	commandHandlerMap["racing+"] = commandRacingPlus
	commandHandlerMap["racingplus"] = commandRacingPlus
	commandHandlerMap["bracket"] = commandBracket
	commandHandlerMap["timezone"] = commandTimezone
	commandHandlerMap["gettimezone"] = commandGetTimezone
	commandHandlerMap["stream"] = commandStream
	commandHandlerMap["getstream"] = commandGetStream
	commandHandlerMap["random"] = commandRandom
	commandHandlerMap["rand"] = commandRandom
	commandHandlerMap["roll"] = commandRandom
	commandHandlerMap["randomchar"] = commandRandomChar
	commandHandlerMap["randomcharacter"] = commandRandomChar
	commandHandlerMap["randchar"] = commandRandomChar
	commandHandlerMap["randcharacter"] = commandRandomChar
	commandHandlerMap["randombuild"] = commandRandomBuild
	commandHandlerMap["randomitem"] = commandRandomBuild
	commandHandlerMap["randbuild"] = commandRandomBuild
	commandHandlerMap["randitem"] = commandRandomBuild
	commandHandlerMap["getnext"] = commandGetNext
	commandHandlerMap["schedule"] = commandSchedule

	// Match commands
	commandHandlerMap["time"] = commandTime
	commandHandlerMap["timeok"] = commandTimeOk
	commandHandlerMap["timedelete"] = commandTimeDelete
	commandHandlerMap["cast"] = commandCast
	commandHandlerMap["castcancel"] = commandCastCancel
	commandHandlerMap["cancelcast"] = commandCastCancel
	commandHandlerMap["castdelete"] = commandCastCancel
	commandHandlerMap["deletecast"] = commandCastCancel
	commandHandlerMap["caster"] = commandCaster
	commandHandlerMap["casterok"] = commandCasterOk
	commandHandlerMap["casternotok"] = commandCasterNotOk
	commandHandlerMap["casteralwaysok"] = commandCasterAlwaysOk
	commandHandlerMap["casteralwaysnotok"] = commandCasterAlwaysNotOk
	commandHandlerMap["ban"] = commandBan
	commandHandlerMap["pick"] = commandPick
	commandHandlerMap["yes"] = commandYes
	commandHandlerMap["no"] = commandNo
	commandHandlerMap["score"] = commandScore
	commandHandlerMap["status"] = commandStatus

	// Admin-only commands
	commandHandlerMap["settimezone"] = commandSetTimezone
	commandHandlerMap["timezoneset"] = commandSetTimezone
	commandHandlerMap["setstream"] = commandSetStream
	commandHandlerMap["streamset"] = commandSetStream
	commandHandlerMap["setcasterok"] = commandSetCasterOk
	commandHandlerMap["casterokset"] = commandSetCasterOk
	commandHandlerMap["setcasternotok"] = commandSetCasterNotOk
	commandHandlerMap["casternotokset"] = commandSetCasterNotOk
	commandHandlerMap["setcasteralwaysok"] = commandSetCasterAlwaysOk
	commandHandlerMap["casteralwaysokset"] = commandSetCasterAlwaysOk
	commandHandlerMap["setcasteralwaysnotok"] = commandSetCasterAlwaysNotOk
	commandHandlerMap["casteralwaysnotokset"] = commandSetCasterAlwaysNotOk
	commandHandlerMap["checkround"] = commandCheckRound
	commandHandlerMap["startround"] = commandStartRound
	commandHandlerMap["roundstart"] = commandStartRound
	commandHandlerMap["start"] = commandStartRound
	commandHandlerMap["beginround"] = commandStartRound
	commandHandlerMap["roundbegin"] = commandStartRound
	commandHandlerMap["begin"] = commandStartRound
	commandHandlerMap["endround"] = commandEndRound
	commandHandlerMap["roundend"] = commandEndRound
	commandHandlerMap["end"] = commandEndRound
	commandHandlerMap["forcetime"] = commandForceTime
	commandHandlerMap["timeforce"] = commandForceTime
	commandHandlerMap["forcetimeok"] = commandForceTimeOk
	commandHandlerMap["timeokforce"] = commandForceTimeOk
	commandHandlerMap["forcetimedelete"] = commandForceTimeDelete
	commandHandlerMap["timedeleteforce"] = commandForceTimeDelete
	commandHandlerMap["forceban"] = commandForceBan
	commandHandlerMap["banforce"] = commandForceBan
	commandHandlerMap["forcepick"] = commandForcePick
	commandHandlerMap["pickforce"] = commandForcePick
	commandHandlerMap["forceyes"] = commandForceYes
	commandHandlerMap["yesforce"] = commandForceYes
	commandHandlerMap["forceno"] = commandForceNo
	commandHandlerMap["noforce"] = commandForceNo
	commandHandlerMap["join"] = commandJoin
	commandHandlerMap["getstate"] = commandGetState
	commandHandlerMap["getchannelid"] = commandGetChannelID
	commandHandlerMap["debug"] = commandDebug
}
