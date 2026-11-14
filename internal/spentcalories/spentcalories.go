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
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {

	parts := strings.Split(data, ",")

	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf("Неверный формат данных: %q", data)
	}

	steps, err := strconv.Atoi(parts[0])

	if err != nil {
		return 0, "", 0, fmt.Errorf("Неверное значение шагов: %q: %w", parts[0], err)
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("Неверное значение шагов (отрицательное количество или ноль): %q", parts[0])

	}

	duration, err := time.ParseDuration(parts[2])

	if err != nil {
		return 0, "", 0, fmt.Errorf("Неверное значение продолжительности активности: %q: %w", parts[2], err)
	}

	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("Неверное значение продолжительности активности (отрицательное или 0): %.2f", duration.Hours())
	}

	typeActivity := parts[1]

	return steps, typeActivity, duration, nil

}

func distance(steps int, height float64) float64 {

	stepLength := stepLengthCoefficient * height

	return (float64(steps) * stepLength) / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {

	if duration <= 0 {
		return 0
	}

	return distance(steps, height) / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {

	steps, typeActivity, duration, err := parseTraining(data)

	if err != nil {
		log.Println(err)
		return "", err
	}

	var spentCalories float64

	switch typeActivity {
	case "Ходьба":
		spentCalories, err = WalkingSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	case "Бег":
		spentCalories, err = RunningSpentCalories(steps, weight, height, duration)
		if err != nil {
			log.Println(err)
			return "", err
		}
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	dist := distance(steps, height)
	meanSpeed := meanSpeed(steps, height, duration)
	result := fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", typeActivity, duration.Hours(), dist, meanSpeed, spentCalories)

	return result, nil

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 {
		return 0, fmt.Errorf("Неверное значение количества шагов (отрицательное или 0): %d", steps)
	}
	if height <= 0 {
		return 0, fmt.Errorf("Неверное значение роста (отрицательное или 0): %.2f", height)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("Неверное значение веса (отрицательное или 0): %.2f", weight)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("Неверное значение продолжительности бега (отрицательное или 0): %.2f", duration.Hours())
	}

	durationInMinutes := duration.Minutes()
	meanSpeed := meanSpeed(steps, height, duration)

	return (weight * meanSpeed * durationInMinutes) / float64(minInH), nil

}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {

	if steps <= 0 {
		return 0, fmt.Errorf("Неверное значение количества шагов (отрицательное или 0): %d", steps)
	}
	if height <= 0 {
		return 0, fmt.Errorf("Неверное значение роста (отрицательное или 0): %.2f", height)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("Неверное значение веса (отрицательное или 0): %.2f", weight)
	}

	durationInMinutes := duration.Minutes()
	meanSpeed := meanSpeed(steps, height, duration)

	spentCalories := (weight * meanSpeed * durationInMinutes) / float64(minInH)

	return walkingCaloriesCoefficient * spentCalories, nil
}
