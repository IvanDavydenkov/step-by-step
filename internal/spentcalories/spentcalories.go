package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	//	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

// "3456,Ходьба,3h00m"
func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")

	if len(parts) != 3 {
		return 0, "", 0, errors.New("ошибка парсинга")
	}
	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, errors.New("шаги не могут быть отрицательными")
	}
	duration, err := time.ParseDuration(parts[2])

	if err != nil {
		return 0, "", 0, err
	}
	if duration <= 0 {
		return 0, "", 0, errors.New("длительность не должна быть равно нулю или меньше его")
	}
	return steps, parts[1], duration, nil
}

func distance(steps int, height float64) float64 {
	return height * stepLengthCoefficient * float64(steps) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return float64(0)
	}
	d := distance(steps, height)
	return d / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, typeTraining, duration, err := parseTraining(data)

	if err != nil {
		log.Println(err)
		return "", err
	}

	var energy float64
	switch typeTraining {
	case "Ходьба":
		value, err := WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		energy = value
	case "Бег":
		value, err := RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			return "", err
		}
		energy = value
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	result := fmt.Sprintf(
		"Тип тренировки: %s\n"+
			"Длительность: %.2f ч.\n"+
			"Дистанция: %.2f км.\n"+
			"Скорость: %.2f км/ч\n"+
			"Сожгли калорий: %.2f\n",
		typeTraining,
		duration.Hours(),
		distance(steps, height),
		meanSpeed(steps, height, duration),
		energy,
	)
	return result, nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || height <= 0 || weight <= 0 || duration <= 0 {
		return float64(0), errors.New("некорректный формат")
	}
	speed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()

	return (weight * speed * durationInMinutes) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || height <= 0 || weight <= 0 || duration <= 0 {
		return float64(0), errors.New("некорректный формат")
	}

	speed := meanSpeed(steps, height, duration)

	durationInMinutes := duration.Minutes()

	return (weight * speed * durationInMinutes) / minInH * walkingCaloriesCoefficient, nil
}
