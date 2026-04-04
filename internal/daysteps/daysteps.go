package daysteps

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

type DaySteps struct {
	Steps    int
	Duration time.Duration
	personaldata.Personal
}

func (ds *DaySteps) Parse(datastring string) error {
	parts := strings.Split(datastring, ",")
	if len(parts) != 2 {
		return fmt.Errorf("неверный формат данных прогулки: ожидается 'шаги,длительность'")
	}

	stepsStr := strings.TrimSpace(parts[0])
	if stepsStr != parts[0] {
		return fmt.Errorf("ошибка преобразования шагов: недопустимые пробелы")
	}
	steps, err := strconv.Atoi(stepsStr)
	if err != nil {
		return fmt.Errorf("ошибка преобразования шагов: %w", err)
	}
	if steps <= 0 {
		return fmt.Errorf("ошибка преобразования шагов: значение должно быть положительным")
	}
	ds.Steps = steps

	durationStr := strings.TrimSpace(parts[1])
	if durationStr != parts[1] {
		return fmt.Errorf("ошибка преобразования длительности: недопустимые пробелы")
	}
	dur, err := time.ParseDuration(durationStr)
	if err != nil {
		return fmt.Errorf("ошибка преобразования длительности: %w", err)
	}
	if dur <= 0 {
		return fmt.Errorf("ошибка преобразования длительности: значение должно быть положительным")
	}
	ds.Duration = dur

	return nil
}

func (ds DaySteps) ActionInfo() (string, error) {
	dist := spentenergy.Distance(ds.Steps, ds.Height)
	calories, err := spentenergy.WalkingSpentCalories(ds.Steps, ds.Weight, ds.Height, ds.Duration)
	if err != nil {
		return "", fmt.Errorf("ошибка расчёта калорий прогулки: %w", err)
	}

	return fmt.Sprintf(
		"Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n",
		ds.Steps, dist, calories,
	), nil
}
