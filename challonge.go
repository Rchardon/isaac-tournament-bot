package main

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type Tournament struct {
	Name              string
	ChallongeURL      string
	ChallongeID       float64
	Ruleset           Ruleset
	DiscordCategoryID string
	BestOf            int
}

var (
	challongeUsername string
	challongeAPIKey   string
	tournaments       = make(map[string]Tournament) // Indexed by Challonge URL suffix.

	// We don't want to use the default http.Client structure because it has no default timeout set.
	myHTTPClient = &http.Client{
		Timeout: 10 * time.Second,
	}
)

func challongeInit() {
	// Read the Challonge configuration from the environment variables.
	challongeUsername = os.Getenv("CHALLONGE_USERNAME")
	if len(challongeUsername) == 0 {
		log.Fatal("The \"CHALLONGE_USERNAME\" environment variable is blank. Set it in the \".env\" file.")
		return
	}

	challongeAPIKey = os.Getenv("CHALLONGE_API_KEY")
	if len(challongeAPIKey) == 0 {
		log.Fatal("The \"CHALLONGE_API_KEY\" environment variable is blank. Set it in the \".env\" file.")
		return
	}

	tournamentURLsString := os.Getenv("TOURNAMENT_CHALLONGE_URLS")
	if len(tournamentURLsString) == 0 {
		log.Fatal("The \"TOURNAMENT_CHALLONGE_URLS\" environment variable is blank. Set it in the \".env\" file.")
		return
	}
	tournamentURLs := strings.Split(tournamentURLsString, ",")

	tournamentRulesetsString := os.Getenv("TOURNAMENT_RULESETS")
	if len(tournamentRulesetsString) == 0 {
		log.Fatal("The \"TOURNAMENT_RULESETS\" environment variable is blank. Set it in the \".env\" file.")
		return
	}
	tournamentRulesetsStringSlice := strings.Split(tournamentRulesetsString, ",")
	tournamentRulesets := make([]Ruleset, 0)
	for _, ruleset := range tournamentRulesetsStringSlice {
		if ruleset != "seeded" && ruleset != "unseeded" && ruleset != "team" {
			log.Fatal("The \"TOURNAMENT_RULESETS\" environment variable is set to \"" + ruleset + "\", which is an invalid value.")
			return
		}
		tournamentRulesets = append(tournamentRulesets, Ruleset(ruleset))
	}

	tournamentDiscordCategoryIDsString := os.Getenv("TOURNAMENT_DISCORD_CATEGORY_IDS")
	if len(tournamentDiscordCategoryIDsString) == 0 {
		log.Fatal("The \"TOURNAMENT_DISCORD_CATEGORY_IDS\" environment variable is blank. Set it in the \".env\" file.")
		return
	}
	tournamentDiscordCategoryIDs := strings.Split(tournamentDiscordCategoryIDsString, ",")

	tournamentBestOfString := os.Getenv("TOURNAMENT_BEST_OF")
	if len(tournamentBestOfString) == 0 {
		log.Fatal("The \"TOURNAMENT_BEST_OF\" environment variable is blank. Set it in the \".env\" file.")
		return
	}
	tournamentBestOfStrings := strings.Split(tournamentBestOfString, ",")

	// Validate that all of the "best of" values are numbers.
	tournamentBestOf := make([]int, 0)
	for _, bestOfString := range tournamentBestOfStrings {
		if v, err := strconv.Atoi(bestOfString); err != nil {
			log.Fatal("One of the values in the \"TOURNAMENT_BEST_OF\" environment variable is not a number.")
			return
		} else {
			tournamentBestOf = append(tournamentBestOf, v)
		}
	}

	// Get all of the Challonge user's tournaments.
	apiURL := "https://api.challonge.com/v1/tournaments.json?"
	apiURL += "api_key=" + challongeAPIKey
	var raw []byte
	if v, err := challongeGetJSON("GET", apiURL, nil); err != nil {
		log.Fatal("Failed to get the tournament from Challonge:", err)
		return
	} else {
		raw = v
	}

	jsonTournaments := make([]interface{}, 0)
	if err := json.Unmarshal(raw, &jsonTournaments); err != nil {
		log.Fatal("Failed to unmarshal the Challonge JSON:", err)
	}

	// Figure out the ID for all the tournaments listed in the environment variable.
	for i, tournamentURL := range tournamentURLs {
		found := false
		for _, v := range jsonTournaments {
			vMap := v.(map[string]interface{})
			jsonTournament := vMap["tournament"].(map[string]interface{})
			if jsonTournament["url"] == tournamentURL {
				found = true
				tournaments[tournamentURL] = Tournament{
					Name:              jsonTournament["name"].(string),
					ChallongeURL:      tournamentURL,
					ChallongeID:       jsonTournament["id"].(float64),
					Ruleset:           tournamentRulesets[i],
					DiscordCategoryID: tournamentDiscordCategoryIDs[i],
					BestOf:            tournamentBestOf[i],
				}
				break
			}
		}
		if !found {
			log.Fatal("Failed to find the \"" + tournamentURL + "\" tournament in this Challonge user's tournament list.")
		}
	}
}

func challongeGetJSON(method string, apiURL string, data io.Reader) ([]byte, error) {
	var req *http.Request
	if v, err := http.NewRequest(method, apiURL, data); err != nil {
		return nil, err
	} else {
		req = v
	}

	var resp *http.Response
	if v, err := myHTTPClient.Do(req); err != nil {
		return nil, err
	} else {
		resp = v
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Bad return status: " + strconv.Itoa(resp.StatusCode))
	}

	var raw []byte
	if v, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, err
	} else {
		raw = v
	}

	return raw, nil
}

func challongeGetParticipantName(tournament map[string]interface{}, participantID float64) string {
	// Go through all of the participants in this tournament.
	for _, v := range tournament["participants"].([]interface{}) {
		vMap := v.(map[string]interface{})
		participant := vMap["participant"].(map[string]interface{})

		// First, check the normal ID.
		if participant["id"].(float64) == participantID {
			return participant["name"].(string)
		}

		// Second, all check the group player IDs.
		// (This is needed if the tournament happens to have group stages.)
		groupIDs := participant["group_player_ids"].([]interface{})
		for _, groupID := range groupIDs {
			if groupID.(float64) == participantID {
				return participant["name"].(string)
			}
		}
	}

	return "Unknown-" + floatToString(participantID)
}
