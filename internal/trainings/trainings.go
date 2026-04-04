package trainings

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/personaldata"
	"github.com/Yandex-Practicum/tracker/internal/spentenergy"
)

type Training struct {
	Steps        int
	TrainingType string
	Duration     time.Duration
	personaldata.Personal
}

func (t *Training) Parse(datastring string) error {
	parts := strings.Split(datastring, ",")
	if len(parts) != 3 {
		return fmt.Errorf("неверный формат данных тренировки: ожидается 'шаги,тип,длительность'")
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
	t.Steps = steps

	t.TrainingType = strings.TrimSpace(parts[1])

	durationStr := strings.TrimSpace(parts[2])
	if durationStr != parts[2] {
		return fmt.Errorf("ошибка преобразования длительности: недопустимые пробелы")
	}
	dur, err := time.ParseDuration(durationStr)
	if err != nil {
		return fmt.Errorf("ошибка преобразования длительности: %w", err)
	}
	if dur <= 0 {
		return fmt.Errorf("ошибка преобразования длительности: значение должно быть положительным")
	}
	t.Duration = dur

	return nil
}

func (t Training) ActionInfo() (string, error) {
	dist := spentenergy.Distance(t.Steps, t.Height)
	speed := spentenergy.MeanSpeed(t.Steps, t.Height, t.Duration)

	var calories float64
	var err error

	switch t.TrainingType {
	case "Бег":
		calories, err = spentenergy.RunningSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	case "Ходьба":
		calories, err = spentenergy.WalkingSpentCalories(t.Steps, t.Weight, t.Height, t.Duration)
	default:
		return "", errors.New("неизвестный тип тренировки")
	}

	if err != nil {
		return "", fmt.Errorf("ошибка расчёта калорий: %w", err)
	}

	return fmt.Sprintf(
		"Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n",
		t.TrainingType, t.Duration.Hours(), dist, speed, calories,
	), nil
}
