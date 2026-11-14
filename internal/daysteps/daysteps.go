package daysteps

import (
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

	parts := strings.Split(data, ",")

	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("Неверный формат данных: %q", data)
	}

	steps, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("Неверное значение шагов: %q: %w", parts[0], err)
	}

	if steps <= 0 {
		return 0, 0, fmt.Errorf("Неверное значение шагов (отрицательнок количество): %q", parts[0])
	}

	duration, err := time.ParseDuration(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("Неверное значение продолжительности: %q: %w", parts[1], err)
	}

	if duration <= 0 {
		return 0, 0, fmt.Errorf("Неверное значение продолжительности прогулки (отрицательное или 0): %.2f", duration.Hours())
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {

	steps, duration, err := parsePackage(data)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	if steps <= 0 {
		return ""
	}

	distanceMeters := float64(steps) * stepLength

	distanceKilometers := distanceMeters / float64(mInKm)

	spentCalories, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	if err != nil {
		fmt.Println(err)

		return ""
	}

	result := fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", steps, distanceKilometers, spentCalories)
	return result
}
