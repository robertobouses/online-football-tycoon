package match

import (
	"errors"
	"fmt"
	"log"
	"math/rand"

	"github.com/robertobouses/online-football-tycoon/internal/domain"
)

type EventType string

const (
	EventTypeKeyPass            EventType = "KEY_PASS"
	EventTypeShot               EventType = "SHOT"
	EventTypePenaltyKick        EventType = "PENALTY_KICK"
	EventTypeLongShot           EventType = "LONG_SHOT"
	EventTypeIndirectFreeKick   EventType = "INDIRECT_FREE_KICK"
	EventTypeDribble            EventType = "DRIBBLE"
	EventTypeFoul               EventType = "FOUL"
	EventTypeYellowOrRedCard    EventType = "YELLOW_OR_RED_CARD"
	EventTypeDirectFreeKick     EventType = "DIRECT_FREE_KICK"
	EventTypeGreatScoringChance EventType = "GREAT_SCORING_CHANCE"
	EventTypeCornerKick         EventType = "CORNER_KICK"
	EventTypeInjuryDuringMatch  EventType = "INJURY_DURING_MATCH"
	EventTypeOffside            EventType = "OFFSIDE"
	EventTypeHeaded             EventType = "HEADED"
	EventTypeCounterAttack      EventType = "COUNTER_ATTACK"
	EventTypeEndOfTheMatch      EventType = "END_OF_THE_MATCH"
	EventTypeMatchBreak         EventType = "MATCH_BREAK"
)

func CalculateSuccessIndividualEvent(skill int) int {
	log.Printf("Evaluating success of individual event for skill level: %d", skill)

	switch {
	case skill < 8:
		return 0
	case skill >= 8 && skill < 14:
		return ProbabilisticIncrement14()
	case skill >= 14 && skill < 21:
		return ProbabilisticIncrement20()
	case skill >= 21 && skill < 30:
		return ProbabilisticIncrement25()
	case skill >= 30 && skill < 36:
		return ProbabilisticIncrement33()
	case skill >= 36 && skill < 44:
		return ProbabilisticIncrement40()
	case skill >= 44 && skill < 59:
		return ProbabilisticIncrement50()
	case skill >= 59 && skill < 68:
		return ProbabilisticIncrement66()
	case skill >= 68 && skill < 74:
		return ProbabilisticIncrement75()
	case skill >= 74 && skill < 92:
		return ProbabilisticIncrement90()

	default:
		return 1
	}
}

func CalculateSuccessConfrontation(atackerSkill, defenderSkill int) int {
	log.Printf("Calculating confrontation success: Attacker skill = %d, Defender skill = %d", atackerSkill, defenderSkill)

	switch {
	case atackerSkill < defenderSkill-91:
		return 0
	case atackerSkill >= defenderSkill-91 && atackerSkill < defenderSkill-74:
		return ProbabilisticIncrement14()
	case atackerSkill >= defenderSkill-74 && atackerSkill < defenderSkill-69:
		return ProbabilisticIncrement20()
	case atackerSkill >= defenderSkill-69 && atackerSkill < defenderSkill-61:
		return ProbabilisticIncrement25()
	case atackerSkill >= defenderSkill-61 && atackerSkill < defenderSkill-52:
		return ProbabilisticIncrement33()
	case atackerSkill >= defenderSkill-52 && atackerSkill < defenderSkill-43:
		return ProbabilisticIncrement40()
	case atackerSkill >= defenderSkill-43 && atackerSkill < defenderSkill-30:
		return ProbabilisticIncrement44()
	case atackerSkill >= defenderSkill-30 && atackerSkill < defenderSkill-12:
		return ProbabilisticIncrement50()
	case atackerSkill >= defenderSkill-12 && atackerSkill < defenderSkill:
		return ProbabilisticIncrement57()
	case atackerSkill >= defenderSkill && atackerSkill < defenderSkill+20:
		return ProbabilisticIncrement62()
	case atackerSkill >= defenderSkill+20 && atackerSkill < defenderSkill+33:
		return ProbabilisticIncrement66()
	case atackerSkill >= defenderSkill+33 && atackerSkill < defenderSkill-37:
		return ProbabilisticIncrement71()
	case atackerSkill >= defenderSkill+37 && atackerSkill < defenderSkill+49:
		return ProbabilisticIncrement75()
	case atackerSkill >= defenderSkill+49 && atackerSkill < defenderSkill+64:
		return ProbabilisticIncrement80()
	case atackerSkill >= defenderSkill+64 && atackerSkill < defenderSkill+77:
		return ProbabilisticIncrement90()
	case atackerSkill >= defenderSkill+77 && atackerSkill < defenderSkill+96:
		return ProbabilisticIncrement94()
	case atackerSkill >= defenderSkill+96:
		return 1

	default:
		return 0
	}
}

func KeyPass(lineup, rivalLineup domain.Team) (string, int, int, int, int, error) {

	passer := GetRandomMidfielder(lineup.Players)
	receiver := GetRandomForward(lineup.Players)
	log.Printf("Selected passer: %+v, receiver: %+v", passer, receiver)

	if passer == nil || receiver == nil {
		return "There are not enough players available to make a pass", 0, 0, 0, 0, fmt.Errorf("There are not enough players available to make a pass")
	}
	successfulPass := CalculateSuccessIndividualEvent(passer.Technique)
	var sentence string
	var lineupChances, rivalChances, lineupGoals, rivalGoals int

	if successfulPass == 1 {
		log.Printf("Pass success calculated: %d", successfulPass)
		sentence := fmt.Sprintf("%s makes a key pass to %s.", passer.LastName, receiver.LastName)
		log.Println(sentence)

		lineupChances = 1
		if resultOfEvent := ProbabilisticIncrement14(); resultOfEvent == 1 {
			PenaltyKick(lineup, rivalLineup)
		} else {
			Shot(lineup, rivalLineup, passer)
		}

		return sentence, lineupChances, rivalChances, lineupGoals, rivalGoals, nil
	}

	sentence = fmt.Sprintf("%s fails to make a key pass to %s.", passer.LastName, receiver.LastName)
	log.Println(sentence)

	lineupChances = 0

	return sentence, lineupChances, rivalChances, lineupGoals, rivalGoals, nil
}

func Shot(lineup, rivalLineup domain.Team, passer *domain.Player) (string, int, int, int, int, error) {

	shooter := GetRandomForward(lineup.Players)
	if shooter == nil {
		return "no forward player found in lineup", 0, 0, 0, 0, errors.New("no forward player found in lineup")
	}
	goalkeeper := GetGoalkeeper(rivalLineup.Players)
	if goalkeeper == nil {
		return "no goalkeeper found in rival lineup", 0, 0, 0, 0, errors.New("no goalkeeper found in rival lineup")
	}
	defender := GetRandomDefender(rivalLineup.Players)
	if defender == nil {
		return "no defender player found in lineup", 0, 0, 0, 0, errors.New("no defender player found in lineup")
	}

	log.Printf("Shooter: %+v, Defender: %+v, Goalkeeper: %+v", shooter, defender, goalkeeper)

	var sentence string
	var lineupChances, rivalChances, lineupGoals, rivalGoals int

	successfulAgainstDefender := CalculateSuccessConfrontation(shooter.Technique, defender.Technique)
	log.Printf("Success against defender: %d", successfulAgainstDefender)

	if successfulAgainstDefender == 1 {
		sentence = fmt.Sprintf("%s escapes from %s's marking...", shooter.LastName, defender.LastName)
		lineupChances = 1
		log.Println("the forward beats the defender")

		log.Printf("%s supera a %s.\n", shooter.LastName, defender.LastName)

		successfulAgainstGoalkeeper := CalculateSuccessConfrontation(shooter.Technique, goalkeeper.Technique)

		if successfulAgainstGoalkeeper == 1 {
			sentence += fmt.Sprintf(" %s shoots and also beats the goalkeeper... GOOOOOAL! %s is just a spectator in the play %s scores a goal!\n", shooter.LastName, goalkeeper.LastName, shooter.LastName)
			log.Println(sentence)
			log.Println("The striker also beats the goalkeeper, it's a GOAL")
			lineupGoals = 1

			if passer != nil {
				log.Printf("the passer is: %v", passer)

			}

			return sentence, lineupChances, rivalChances, lineupGoals, rivalGoals, nil
		} else {
			sentence += fmt.Sprintf(" %s's shot is saved by %s.\n", shooter.LastName, goalkeeper.LastName)
			log.Println(sentence)

		}
		lineupChances = 1
	} else {
		log.Printf("the defender blocked the shot in the first instance")
		sentence += fmt.Sprintf(" %s's shot is blocked by %s.\n", shooter.LastName, defender.LastName)
		log.Println(sentence)

		lineupChances = 0
		return sentence, lineupChances, rivalChances, lineupGoals, rivalGoals, nil
	}
	return sentence, lineupChances, rivalChances, lineupGoals, rivalGoals, nil
}

func PenaltyKick(lineup, rivalLineup domain.Team) (string, int, int, int, int, error) {
	shooter := GetRandomForward(lineup.Players)
	if shooter == nil {
		return "no forward player found in lineup", 0, 0, 0, 0, errors.New("no forward player found in lineup")
	}
	goalkeeper := GetGoalkeeper(rivalLineup.Players)
	if goalkeeper == nil {
		return "no goalkeeper found in rival lineup", 0, 0, 0, 0, errors.New("no goalkeeper found in rival lineup")
	}

	increasedShooterMental := shooter.Mental + (10 * rand.Intn(3))
	decreasedGoalkeeperMental := goalkeeper.Mental - 5

	successfulPenalty := CalculateSuccessConfrontation(increasedShooterMental, decreasedGoalkeeperMental)

	var sentence string
	var lineupChances, rivalChances, lineupGoals, rivalGoals int

	if successfulPenalty == 1 {
		sentence := fmt.Sprintf("GOOOOOAL! %s scores from the penalty spot!", shooter.LastName)
		log.Println(sentence)

		lineupGoals = 1
		lineupChances = 1

		return sentence, lineupChances, rivalChances, lineupGoals, rivalGoals, nil
	} else {
		sentence = fmt.Sprintf("%s's penalty is saved by %s.\n", shooter.LastName, goalkeeper.LastName)
		log.Println(sentence)

		lineupChances = 1

		return sentence, lineupChances, rivalChances, lineupGoals, rivalGoals, nil
	}
}

func LongShot(lineup, rivalLineup domain.Team) (string, int, int, int, int, error) {

	shooter := GetRandomForward(lineup.Players)
	if shooter == nil {
		return "no forward player found in lineup", 0, 0, 0, 0, errors.New("no forward player found in lineup")
	}
	goalkeeper := GetGoalkeeper(rivalLineup.Players)
	if goalkeeper == nil {
		return "no goalkeeper found in rival lineup", 0, 0, 0, 0, errors.New("no goalkeeper found in rival lineup")
	}

	decreasedShooterTechnique := shooter.Technique - (6 * rand.Intn(4))

	successfulLongShot := CalculateSuccessConfrontation(decreasedShooterTechnique, goalkeeper.Mental)

	lineupChances := 1
	rivalChances := 0
	lineupGoals := 0
	rivalGoals := 0

	if successfulLongShot == 1 {
		sentence := fmt.Sprintf("GOOOOOAL! %s scores from long distance!\n", shooter.LastName)
		log.Printf("GOAL! %s scores a long shot!\n", shooter.LastName)

		lineupGoals = 1
		lineupChances = 1

		return sentence, lineupChances, rivalChances, lineupGoals, rivalGoals, nil
	} else {
		sentence := fmt.Sprintf("%s's long shot is saved by %s.\n", shooter.LastName, goalkeeper.LastName)
		log.Printf("SAVE! %s's long shot is saved by %s.\n", shooter.LastName, goalkeeper.LastName)

		lineupChances = 1

		return sentence, lineupChances, rivalChances, lineupGoals, rivalGoals, nil
	}
}

func IndirectFreeKick(lineup, rivalLineup domain.Team) (string, int, int, int, int, error) {

	shooter := GetRandomMidfielder(lineup.Players)
	if shooter == nil {
		return "no shooter player found in lineup", 0, 0, 0, 0, errors.New("no shooter player found in lineup")
	}
	defenderOnAttack := GetRandomDefender(lineup.Players)
	if defenderOnAttack == nil {
		return "no defender player found in lineup", 0, 0, 0, 0, errors.New("no defender player found in lineup")
	}
	rivalDefender := GetRandomDefender(rivalLineup.Players)
	if rivalDefender == nil {
		return "no rivalDefender player found in lineup", 0, 0, 0, 0, errors.New("no rivalDefender player found in lineup")
	}
	goalkeeper := GetGoalkeeper(rivalLineup.Players)
	if goalkeeper == nil {
		return "no goalkeeper found in rival lineup", 0, 0, 0, 0, errors.New("no goalkeeper found in rival lineup")
	}

	increasedShooterTechnique := shooter.Technique + (4 * rand.Intn(6))
	increasedRivalDefenderPhysique := rivalDefender.Physique + rand.Intn(30)

	attackAtributes := increasedShooterTechnique + defenderOnAttack.Physique
	defenseAtributes := increasedRivalDefenderPhysique + goalkeeper.Technique

	successfulLongShot := CalculateSuccessConfrontation(attackAtributes, defenseAtributes)

	lineupChances := 1
	rivalChances := 0
	lineupGoals := 0
	rivalGoals := 0

	if successfulLongShot == 1 {
		sentence := fmt.Sprintf("%s takes the free kick... It is a center to the are.. %s and %s jump to fight for the center... %s head the ball...%s can't do anything...  GOOOOOAL! %s scores!\n", shooter.LastName, defenderOnAttack.LastName, rivalDefender.LastName, defenderOnAttack.LastName, goalkeeper.LastName, defenderOnAttack.LastName)
		log.Printf("GOAL! %s scores a long shot!\n", shooter.LastName)

		lineupGoals = 1
		lineupChances = 1

		return sentence, lineupChances, rivalChances, lineupGoals, rivalGoals, nil
	} else {
		sentence := fmt.Sprintf("%s's long shot is saved by %s.\n", shooter.LastName, goalkeeper.LastName)
		log.Printf("SAVE! %s's long shot is saved by %s.\n", shooter.LastName, goalkeeper.LastName)

		lineupChances = 1

		return sentence, lineupChances, rivalChances, lineupGoals, rivalGoals, nil
	}
}

func Dribble(lineup, rivalLineup domain.Team) (string, int, int, int, int, error) {
	var dribbler, defender *domain.Player
	var sentence string
	var lineupChances, rivalChances, lineupGoals, rivalGoals int

	prob := ProbabilisticIncrement20() + ProbabilisticIncrement20()
	if prob <= 0 {
		dribbler = GetRandomMidfielder(lineup.Players)
		defender = GetRandomMidfielder(rivalLineup.Players)
	} else if prob <= 1 {
		dribbler = GetRandomForward(lineup.Players)
		defender = GetRandomDefender(rivalLineup.Players)
	} else if prob > 1 {
		dribbler = GetRandomMidfielder(lineup.Players)
		defender = GetRandomDefender(rivalLineup.Players)
	}

	log.Printf("Selected passer: %+v, receiver: %+v", dribbler, defender)

	if dribbler == nil || defender == nil {
		return "There are not enough players available to make a pass", 0, 0, 0, 0, fmt.Errorf("There are not enough players available to make a pass")
	}

	sentence = fmt.Sprintf("%s tries a dribbling", dribbler.LastName)
	successfulDribble := CalculateSuccessIndividualEvent(dribbler.Technique)

	if successfulDribble == 1 {
		log.Printf("Pass success calculated: %d", successfulDribble)
		sentence += " and succeeds..."
		log.Println(sentence)

		successfulConfrontation := CalculateSuccessConfrontation(dribbler.Technique, defender.Technique)
		if successfulConfrontation == 1 {
			sentence += fmt.Sprintf(" %s dribbled %s...", dribbler.LastName, defender.LastName)

			lineupChances = 1

			if resultOfEvent := ProbabilisticIncrement40(); resultOfEvent == 1 {
				sentence += " the occasion ends with a shot"
				log.Println("the occasion ends with a shot")
				Shot(lineup, rivalLineup, dribbler)
			} else {
				sentence += " the occasion ends with a foul"
				log.Println("the occasion ends with a shot")
				Foul(lineup, rivalLineup, defender)
			}

			return sentence, lineupChances, rivalChances, lineupGoals, rivalGoals, nil
		} else {
			sentence += fmt.Sprintf(" but %s lost the dribbled against %s", dribbler.LastName, defender.LastName)

		}

		return sentence, lineupChances, rivalChances, lineupGoals, rivalGoals, nil
	}
	sentence += fmt.Sprintf(" Oh Noo, %s trips over the ball, and fails the dribble", dribbler.LastName)

	log.Println(sentence)

	lineupChances = 0

	return sentence, lineupChances, rivalChances, lineupGoals, rivalGoals, nil
}

func Foul(lineup, rivalLineup domain.Team, defender *domain.Player) (string, int, int, int, int, error) {
	var sentence string
	var lineupChances, rivalChances, lineupGoals, rivalGoals int

	probabilyYellowOrRedCard := ProbabilisticIncrement40() + ProbabilisticIncrement20()

	resultOfEvent := ProbabilisticIncrement50() + ProbabilisticIncrement33() + ProbabilisticIncrement33()

	if probabilyYellowOrRedCard >= 1 {
		sentence = "the referee puts his hand in his pocket"
		YellowOrRedCard(lineup, defender)
	}
	if resultOfEvent >= 2 {
		sentence = "the foul is in the middle of the field"
		return sentence, lineupChances, rivalChances, lineupGoals, rivalGoals, nil
	}
	if resultOfEvent >= 1 {
		sentence = "the foul is in the middle of the field"
		IndirectFreeKick(lineup, rivalLineup)

	} else {
		sentence = "the foul is in a dangerous area of the field"
		DirectFreeKick(lineup, rivalLineup)
	}
	return sentence, lineupChances, rivalChances, lineupGoals, rivalGoals, nil

}

func YellowOrRedCard(lineup domain.Team, defender *domain.Player) (string, int, int, int, int, error) {
	var sentence string
	var lineupChances, rivalChances, lineupGoals, rivalGoals, probabilyYellowCard int
	sentence = "The referee puts his hand in his pocket"

	if defender == nil {
		defender = GetRandomDefender(lineup.Players)
		if defender == nil {
			return "no defender player found in lineup", 0, 0, 0, 0, errors.New("no defender player found in lineup")
		}
	}

	probabilyIncrementByAgressive := CalculateSuccessIndividualEvent(defender.Mental)

	if probabilyIncrementByAgressive >= 1 {
		probabilyYellowCard = ProbabilisticIncrement62()
	} else {
		probabilyYellowCard = ProbabilisticIncrement75()
	}

	if probabilyYellowCard >= 1 {
		sentence += fmt.Sprintf("The referee gives %v a yellow card", defender.LastName)

	} else {

		sentence += fmt.Sprintf("The referee gives %v a red card", defender.LastName)
		log.Println("the player was expelled from the lineup")

	}

	return sentence, lineupChances, rivalChances, lineupGoals, rivalGoals, nil

}

func DirectFreeKick(lineup, rivalLineup domain.Team) (string, int, int, int, int, error) {

	shooter := GetRandomForward(lineup.Players)
	if shooter == nil {
		return "no forward player found in lineup", 0, 0, 0, 0, errors.New("no forward player found in lineup")
	}
	goalkeeper := GetGoalkeeper(rivalLineup.Players)
	if goalkeeper == nil {
		return "no goalkeeper found in rival lineup", 0, 0, 0, 0, errors.New("no goalkeeper found in rival lineup")
	}

	decreasedShooterTechnique := shooter.Technique - (6 * rand.Intn(7))

	successfulLongShot := CalculateSuccessConfrontation(decreasedShooterTechnique, goalkeeper.Technique)

	lineupChances := 1
	rivalChances := 0
	lineupGoals := 0
	rivalGoals := 0

	if successfulLongShot == 1 {
		sentence := fmt.Sprintf("GOOOOOAL! %s scores from direct free kick distance!\n", shooter.LastName)
		log.Printf("GOAL! %s scores a free kick long shot!\n", shooter.LastName)

		lineupGoals = 1
		lineupChances = 1

		return sentence, lineupChances, rivalChances, lineupGoals, rivalGoals, nil
	} else {
		sentence := fmt.Sprintf("%s's free kick shot is saved by %s.\n", shooter.LastName, goalkeeper.LastName)
		log.Printf("SAVE! %s's free kick is saved by %s.\n", shooter.LastName, goalkeeper.LastName)

		lineupChances = 1

		return sentence, lineupChances, rivalChances, lineupGoals, rivalGoals, nil
	}
}

func GreatScoringChance(lineup domain.Team) (string, int, int, int, int, error) {
	var shooter *domain.Player
	var sentence string
	var lineupChances, rivalChances, lineupGoals, rivalGoals int
	prob := ProbabilisticIncrement66()
	if prob == 1 {
		shooter = GetRandomForward(lineup.Players)
		if shooter == nil {
			return "No player available for scoring", 0, 0, 0, 0, fmt.Errorf("no player available for scoring")
		}
	} else {
		shooter = GetRandomMidfielder(lineup.Players)
		if shooter == nil {
			return "No player available for scoring", 0, 0, 0, 0, fmt.Errorf("no player available for scoring")
		}
	}
	prob = ProbabilisticIncrement71()
	lineupChances = 1
	if prob == 1 {
		lineupGoals = 1
		sentence = fmt.Sprintf("%s score a great easy chance", shooter.LastName)

		return sentence, lineupChances, rivalChances, lineupGoals, rivalGoals, nil
	} else {
		sentence = fmt.Sprintf("%s fails miserably with a very clear scoring chance", shooter.LastName)

		return sentence, lineupChances, rivalChances, lineupGoals, rivalGoals, nil

	}
}

func CornerKick(lineup, rivalLineup domain.Team) (string, int, int, int, int, error) {
	var centerer *domain.Player
	var attacker, defender *domain.Player
	var sentence string
	var lineupChances, rivalChances, lineupGoals, rivalGoals int

	centerer = GetRandomMidfielder(lineup.Players)
	if centerer == nil {
		return "", 0, 0, 0, 0, fmt.Errorf("no midfielder found for centerer")
	}

	incrementedTechnique := centerer.Technique + rand.Intn(20)
	prob := CalculateSuccessIndividualEvent(incrementedTechnique)
	lineupChances = 1

	if prob == 1 {
		defender = GetRandomDefender(rivalLineup.Players)
		prob = ProbabilisticIncrement62()
		if prob == 1 {
			attacker = GetRandomDefender(lineup.Players)
		} else {
			attacker = GetRandomMidfielder(lineup.Players)
		}
		prob = CalculateSuccessConfrontation(attacker.Physique, defender.Physique)
		if prob == 1 {
			sentence = fmt.Sprintf("GOOOOOAL, %s took the corner very well, and %s beats %s with a incredible jump and heads at goal", centerer.LastName, attacker.LastName, defender.LastName)
			lineupGoals = 1

			return sentence, lineupChances, rivalChances, lineupGoals, rivalGoals, nil
		} else {
			sentence = fmt.Sprintf("%s takes the corner... but %s beats %s to the jump and clears the ball", centerer.LastName, defender.LastName, attacker.LastName)

			return sentence, lineupChances, rivalChances, lineupGoals, rivalGoals, nil
		}
	} else {
		sentence = fmt.Sprintf("the corner was wasted by %s", centerer.LastName)

		return sentence, lineupChances, rivalChances, lineupGoals, rivalGoals, nil
	}
}

func InjuryDuringMatch(lineup domain.Team) (string, int, int, int, int, error) {
	var injuredPlayer *domain.Player
	var sentence string
	injuredPlayer = GetRandomPlayerExcludingGoalkeeper(lineup.Players)
	if injuredPlayer == nil {
		return "", 0, 0, 0, 0, fmt.Errorf("no midfielder found for injuredPlayer")
	}
	sentence = fmt.Sprintf("There is a player lying on the ground... wow he is %s, he looks like he will need assistance...", injuredPlayer.LastName)

	return sentence, 0, 0, 0, 0, nil
}

func Offside(lineup, rivalLineup domain.Team) (string, int, int, int, int, error) {
	var passer, playerOffside *domain.Player
	var lineupChances int
	var sentence string

	passer = GetRandomMidfielder(lineup.Players)
	if passer == nil {
		return "", 0, 0, 0, 0, fmt.Errorf("no midfielder found for passer")
	}

	playerOffside = GetRandomForward(lineup.Players)
	if playerOffside == nil {
		return "", 0, 0, 0, 0, fmt.Errorf("no midfielder found for playerOffside")
	}

	lineupChances = 1

	sentence = fmt.Sprintf("%s looks a pass... ", passer.LastName)
	sentence += fmt.Sprintf("%s runs behind the rival defense", playerOffside.LastName)

	prob := ProbabilisticIncrement66()
	if prob >= 1 {
		sentence += "%s its offside, the opportunity is lost"

	} else {
		sentence += "great pass bordering on offside"

		Shot(lineup, rivalLineup, passer)
	}

	return sentence, lineupChances, 0, 0, 0, nil
}

func Headed(lineup, rivalLineup domain.Team) (string, int, int, int, int, error) {
	var header, rivalHeader *domain.Player
	var sentence string
	var lineupChances, rivalChances int

	header = GetRandomPlayerExcludingGoalkeeper(lineup.Players)
	rivalHeader = GetRandomPlayerExcludingGoalkeeper(rivalLineup.Players)
	if header == nil || rivalHeader == nil {
		fmt.Println("header or rivalHeader es nil")
		return "", 0, 0, 0, 0, fmt.Errorf("no rival player available for the header duel")
	}
	sentence = "The ball comes through the air, here we have an aerial duel"

	success := CalculateSuccessConfrontation(header.Physique, rivalHeader.Physique)
	if success == 1 {
		lineupChances = 1
		sentence += fmt.Sprintf("%s wins a header in midfield against %s", header.LastName, rivalHeader.LastName)

		LongShot(lineup, rivalLineup)
	} else {
		sentence += fmt.Sprintf("%s loses a header in midfield against %s", header.LastName, rivalHeader.LastName)
		prob := ProbabilisticIncrement75()
		if prob >= 1 {
			rivalChances = 1
			sentence += fmt.Sprintf("%s makes a long pass, and his teammates run away", rivalHeader.LastName)

			CounterAttack(rivalLineup, lineup)

		} else {
			sentence += fmt.Sprintf("%s kick the ball into the air, and there are no second plays", rivalHeader.LastName)

		}

	}
	return sentence, lineupChances, rivalChances, 0, 0, nil
}

func CounterAttack(lineup, rivalLineup domain.Team) (string, int, int, int, int, error) {
	var sentence string

	sentence = "Some players run out in counterattack"
	prob := ProbabilisticIncrement66()
	if prob >= 1 {
		LongShot(lineup, rivalLineup)
	} else {
		prob := ProbabilisticIncrement57()
		if prob >= 1 {
			sentence += "The rival stopped the counterattack with a foul"
			IndirectFreeKick(lineup, rivalLineup)
		} else {
			sentence += "The opponent breaks the counterattack cleanly"
		}
	}

	return sentence, 0, 0, 0, 0, nil
}
