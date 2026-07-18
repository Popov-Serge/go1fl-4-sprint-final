package daysteps

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/constants"
	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

func parsePackage(data string) (int, time.Duration, error) {
	arr := strings.Split(data, ",")
	if len(arr) != 2 {
		return 0, 0, errInvalidData
	}

	stepsCount, err := strconv.Atoi(arr[0])
	if err != nil {
		return 0, 0, err
	}

	if stepsCount <= 0 {
		return 0, 0, errInvalidSteps
	}

	duration, err := time.ParseDuration(arr[1])
	if err != nil {
		return 0, 0, err
	}

	if duration <= 0 {
		return 0, 0, errInvalidDuration
	}

	return stepsCount, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		log.Println(err)
		return ""
	}
	distance := (float64(steps) * constants.StepLength) / constants.MetersInKm

	calories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		log.Println(err)
		return ""
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\n"+
			"Дистанция составила %.2f км.\n"+
			"Вы сожгли %.2f ккал.\n",
		steps, distance, calories)
}
