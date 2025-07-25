package player

import (
	"log"
	"math/rand"
	"time"

	"github.com/robertobouses/online-football-tycoon/internal/domain"
)

const (
	highTechniqueBonus   = 90
	highMentalBonus      = 90
	highPhysiqueBonus    = 90
	mediumTechniqueBonus = 78
	mediumMentalBonus    = 76
	mediumPhysiqueBonus  = 77
	ageSalaryFactor      = 24
	youngValueBonusStart = 22
	youngValueBonusEnd   = 28
	oldValue             = 30
	lowPhsysique         = 25
)

func CalculatePlayerFeeAndSalary(technique, mental, physique, age int, country, position string) (int, int) {
	rand.Seed(time.Now().UnixNano())

	totalPlayerQuality := technique + mental + physique

	baseFee := totalPlayerQuality * totalPlayerQuality * totalPlayerQuality * (rand.Intn(4) + 1)
	baseSalary := totalPlayerQuality * totalPlayerQuality * (rand.Intn(70) + 25)

	fee := baseFee
	if age < youngValueBonusStart {
		fee += (5_000_000 * (youngValueBonusStart - age)) * (rand.Intn(4) + 1)
	} else if age >= youngValueBonusStart && age <= youngValueBonusEnd {
		fee += 10_000_000 * (rand.Intn(2) + 1)
	} else if age >= oldValue {

		fee -= int(6_000_000 * float64(rand.Intn(3)+1) * ((float64(age) * 1.9) - float64(oldValue)))
	}

	salary := int(float64(baseSalary) * (float64(age) / float64(ageSalaryFactor)))

	if technique >= highTechniqueBonus {
		salary += 3_000_000 * (rand.Intn(2) + 1)
		fee += 10_000_000 * (rand.Intn(2) + 1)
	}
	if technique >= mediumTechniqueBonus {
		fee += 2_000_000 * (rand.Intn(2) + 1)
	}
	if mental >= highMentalBonus {
		salary += int(1_000_000 * (rand.Float64()*2.3 + 1.3))
		fee += 8_000_000 * (rand.Intn(2) + 1)
	}
	if mental >= mediumMentalBonus {
		salary += 1_000_000 * (rand.Intn(2) + 1)
	}
	if physique >= highPhysiqueBonus {
		salary += int(1_000_000 * (rand.Float64()*2.4 + 1))
		fee += 8_000_000 * (rand.Intn(2) + 1)
	}
	if physique >= mediumPhysiqueBonus {
		fee += 1_000_000 * (rand.Intn(2) + 1)
	}

	if physique < lowPhsysique {
		fee -= 700_000 * (rand.Intn(3) + 1)
	}

	switch position {
	case domain.PositionGoalkeeper:
		fee -= 2_000_000 * (rand.Intn(2) + 1)
		if totalPlayerQuality < 120 {
			salary -= 200_000 * (rand.Intn(5) + 1)
		}

	case domain.PositionDefender:
		fee -= 1_200_000 * (rand.Intn(2) + 1)

	case domain.PositionMidfielder:
		fee += 500_000 * (rand.Intn(2) + 1)
	case domain.PositionForward:
		fee += 1_550_000 * (rand.Intn(3) + 1)
		salary += 50_000 * (rand.Intn(8) + 1)
	}

	log.Println("valor de salary en CalculatePlayerFeeAndSalary", salary)
	log.Println("valor de fee en CalculatePlayerFeeAndSalary", fee)

	if totalPlayerQuality < 80 {
		fee -= 800_000 * (rand.Intn(2) + 1)
		salary -= 70_000 * (rand.Intn(4) + 1)
	}
	if totalPlayerQuality < 114 {
		fee -= 50_000 * (rand.Intn(3) + 1)
		salary -= 20_000 * (rand.Intn(5) + 1)
	}

	if salary < 0 {
		salary = 160_000
	}

	if fee < 0 {
		fee = 0
	}

	return fee, salary
}
