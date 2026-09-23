package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	arr := strings.Split(data, ",")
	if len(arr) != 2 {
		return 0, 0, errors.New("ошибка парсинга данных, некорректное число элементов после парсинга")
	}
	steps, err := strconv.Atoi(arr[0])
	if steps <= 0 {
		return 0, 0, errors.New("ошибка парсинга шагов")
	}
	if err != nil {
		return 0, 0, err
	}
	duration, err := time.ParseDuration(arr[1])
	if duration <= 0 {
		return 0, 0, errors.New("ошибка парсинга времени")
	}
	if err != nil {
		return 0, 0, err
	}
	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		return ""
	}
	distance := float64(steps) * stepLength / mInKm
	energy, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)
	if err != nil {
		return "ошибка парсинга калорий"
	}
	result := fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		steps, distance, energy,
	)

	return result
}
