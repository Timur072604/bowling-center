package bowling

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

func biasedRoll(rng *rand.Rand, pinsLeft int, playerType string) int {
	if pinsLeft == 0 {
		return 0
	}

	switch playerType {
	case "beginner":
		roll1 := rng.Intn(pinsLeft + 1)
		roll2 := rng.Intn(pinsLeft + 1)
		if roll1 < roll2 {
			return roll1
		}
		return roll2
	case "pro":
		roll1 := rng.Intn(pinsLeft + 1)
		roll2 := rng.Intn(pinsLeft + 1)
		if roll1 > roll2 {
			return roll1
		}
		return roll2
	default:
		return rng.Intn(pinsLeft + 1)
	}
}

func GenerateRandomGame(rng *rand.Rand) string {
	var builder strings.Builder

	playerTypes := []string{"beginner", "average", "pro"}
	playerType := playerTypes[rng.Intn(len(playerTypes))]

	for frame := 0; frame < 9; frame++ {
		firstRoll := biasedRoll(rng, 10, playerType)
		if firstRoll == 10 {
			builder.WriteRune('X')
		} else {
			builder.WriteString(strconv.Itoa(firstRoll))
			pinsLeft := 10 - firstRoll
			secondRoll := biasedRoll(rng, pinsLeft, playerType)
			if firstRoll+secondRoll == 10 {
				builder.WriteRune('/')
			} else {
				builder.WriteString(strconv.Itoa(secondRoll))
			}
		}
	}

	firstRoll := biasedRoll(rng, 10, playerType)
	if firstRoll == 10 {
		builder.WriteRune('X')
		secondRoll := biasedRoll(rng, 10, playerType)
		if secondRoll == 10 {
			builder.WriteRune('X')
			thirdRoll := biasedRoll(rng, 10, playerType)
			if thirdRoll == 10 {
				builder.WriteRune('X')
			} else {
				builder.WriteString(strconv.Itoa(thirdRoll))
			}
		} else {
			builder.WriteString(strconv.Itoa(secondRoll))
			pinsLeft := 10 - secondRoll
			thirdRoll := biasedRoll(rng, pinsLeft, playerType)
			if secondRoll+thirdRoll == 10 {
				builder.WriteRune('/')
			} else {
				builder.WriteString(strconv.Itoa(thirdRoll))
			}
		}
	} else {
		builder.WriteString(strconv.Itoa(firstRoll))
		pinsLeft := 10 - firstRoll
		secondRoll := biasedRoll(rng, pinsLeft, playerType)
		if firstRoll+secondRoll == 10 {
			builder.WriteRune('/')
			thirdRoll := biasedRoll(rng, 10, playerType)
			if thirdRoll == 10 {
				builder.WriteRune('X')
			} else {
				builder.WriteString(strconv.Itoa(thirdRoll))
			}
		} else {
			builder.WriteString(strconv.Itoa(secondRoll))
		}
	}
	return builder.String()
}

func Score(game string) (int, error) {
	rolls := make([]int, 0, 21)
	for _, char := range game {
		switch {
		case char >= '0' && char <= '9':
			rolls = append(rolls, int(char-'0'))
		case char == '-':
			rolls = append(rolls, 0)
		case char == 'X' || char == 'x':
			rolls = append(rolls, 10)
		case char == '/':
			if len(rolls) == 0 {
				return 0, fmt.Errorf("spare не может быть первым броском")
			}
			prev := rolls[len(rolls)-1]
			rolls = append(rolls, 10-prev)
		default:
			return 0, fmt.Errorf("недопустимый символ в строке игры: %c", char)
		}
	}

	score := 0
	rollIndex := 0
	for frame := 0; frame < 10; frame++ {
		if rollIndex >= len(rolls) {
			break
		}
		if rolls[rollIndex] == 10 {
			if rollIndex+2 >= len(rolls) {
				return 0, fmt.Errorf("недостаточно бросков после strike")
			}
			score += 10 + rolls[rollIndex+1] + rolls[rollIndex+2]
			rollIndex++
		} else if rollIndex+1 < len(rolls) && rolls[rollIndex]+rolls[rollIndex+1] == 10 {
			if rollIndex+2 >= len(rolls) {
				return 0, fmt.Errorf("недостаточно бросков после spare")
			}
			score += 10 + rolls[rollIndex+2]
			rollIndex += 2
		} else {
			if rollIndex+1 >= len(rolls) {
				return 0, fmt.Errorf("незавершенный фрейм")
			}
			score += rolls[rollIndex] + rolls[rollIndex+1]
			rollIndex += 2
		}
	}
	return score, nil
}
