package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/constants"
)

// Основные константы, необходимые для расчетов.
const (
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
	walkTraining               = "Ходьба"
	runTraining                = "Бег"
)

func parseTraining(data string) (int, string, time.Duration, error) {
	arr := strings.Split(data, ",")
	if len(arr) != 3 {
		return 0, "", 0, errInvalidData
	}

	steps, err := strconv.Atoi(arr[0])
	if err != nil {
		return 0, "", 0, err
	}

	if steps <= 0 {
		return 0, "", 0, errInvalidSteps
	}

	duration, err := time.ParseDuration(arr[2])
	if err != nil {
		return 0, "", 0, err
	}

	if duration <= 0 {
		return 0, "", 0, errInvalidDuration
	}

	trainingType := arr[1]

	return steps, trainingType, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient
	return (stepLength * float64(steps)) / constants.MetersInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	distance := distance(steps, height)

	return distance / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, trainingType, duration, err := parseTraining(data)

	if err != nil {
		log.Println(err)
		return "", err
	}

	distance := distance(steps, height)
	speed := meanSpeed(steps, height, duration)

	var calories float64
	switch trainingType {
	case walkTraining:
		calories, err = WalkingSpentCalories(steps, weight, height, duration)
	case runTraining:
		calories, err = RunningSpentCalories(steps, weight, height, duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		return "", err
	}

	return fmt.Sprintf(""+
		"Тип тренировки: %s\n"+
		"Длительность: %.2f ч.\n"+
		"Дистанция: %.2f км.\n"+
		"Скорость: %.2f км/ч\n"+
		"Сожгли калорий: %.2f\n",
		trainingType,
		duration.Hours(),
		distance,
		speed,
		calories), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	calories, err := caloriesPerMinute(steps, weight, height, duration)
	if err != nil {
		return 0, errInvalidData
	}

	return calories, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	calories, err := caloriesPerMinute(steps, weight, height, duration)
	if err != nil {
		return 0, errInvalidData
	}

	return walkingCaloriesCoefficient * calories, nil
}

func caloriesPerMinute(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errInvalidData
	}

	speed := meanSpeed(steps, height, duration)
	durationMinutes := duration.Minutes()

	return (durationMinutes * speed * weight) / minInH, nil
}
