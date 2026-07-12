package config

// adding this for credit score. I dont think the scoring scheme is final yet,
// so keeping it here for now. Once the scoring scheme is finalized,
// we can move this to a database table and fetch the values from there.
const defaultParticipationCredits = 10

var participationCredits = map[string]int{
	"Winner":      20,
	"1st Place":   20,
	"Runner Up":   15,
	"2nd Place":   15,
	"3rd Place":   15,
	"Coordinator": 12,
	"Organizer":   12,
	"Participant": 10,
	"Volunteer":   8,
}

var eventLevelBonus = map[string]int{
	"National":      5,
	"International": 10,
}

func CreditsForCertificate(participationType, eventLevel string) int {
	credits, ok := participationCredits[participationType]
	if !ok {
		credits = defaultParticipationCredits
	}

	return credits + eventLevelBonus[eventLevel]
}
