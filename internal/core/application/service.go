package application

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/vigorouzis/aibolit-notification/internal/core/domain"
	"github.com/vigorouzis/aibolit-notification/internal/infrastructure/postgres"
	"github.com/vigorouzis/aibolit-notification/internal/utils"
)

type ScheduleService struct {
	pg *postgres.Client
}

func NewService(pg *postgres.Client) *ScheduleService {
	return &ScheduleService{pg: pg}
}

func (s *ScheduleService) CreateSchedule(ctx context.Context, schedule *domain.Schedule) (string, error) {
	schedule.ID = uuid.New().String()
	schedule.StartDate = time.Now()

	if err := s.pg.CreateSchedule(ctx, schedule); err != nil {
		return "", err
	}

	return schedule.ID, nil
}

func (s *ScheduleService) GetSchedulesByUserId(ctx context.Context, userId string) ([]domain.Schedule, error) {
	if userId == "" {
		return nil, errors.New("user_id is required")
	}

	schedules, err := s.pg.GetSchedulesByUserID(ctx, userId)
	if err == sql.ErrNoRows {
		return nil, errors.New("no schedules found")
	} else if err != nil {
		return nil, err
	}

	return schedules, nil
}

func (s *ScheduleService) GetSchedule(ctx context.Context, userID string, scheduleID string) (*domain.Schedule, []string, error) {
	if userID == "" || scheduleID == "" {
		return nil, nil, errors.New("user_id and schedule_id are required")
	}

	schedule, err := s.pg.GetSchedule(ctx, userID, scheduleID)
	if err != nil {
		return nil, nil, err
	}
	if schedule == nil {
		return nil, nil, errors.New("schedule not found")
	}

	if !schedule.IsPermanent && schedule.Duration > 0 {
		endDate := time.Now().AddDate(0, 0, schedule.Duration)
		if time.Now().After(endDate) {
			return nil, nil, errors.New("medication course completed")
		}
	}

	intakeTimes := utils.CalculateIntakeTimes(schedule.Frequency)

	return schedule, intakeTimes, nil
}

func (s *ScheduleService) NextTakings(ctx context.Context, userID string) ([]map[string]string, error) {
	if userID == "" {
		return nil, errors.New("user_id is required")
	}

	schedules, err := s.pg.GetSchedulesByUserID(ctx, userID)
	if err == sql.ErrNoRows {
		return nil, errors.New("no schedules found for the user")
	} else if err != nil {
		return nil, err
	}

	periodStr := os.Getenv("NEXT_TAKINGS_PERIOD")
	period, err := strconv.Atoi(periodStr)
	if err != nil {
		period = 60
	}
	nextTakingsPeriod := time.Duration(period) * time.Minute

	currentTime := time.Now()
	nextTime := currentTime.Add(nextTakingsPeriod)

	var result []map[string]string

	for _, s := range schedules {

		endDate := s.StartDate.AddDate(0, 0, s.Duration)
		if !s.IsPermanent && s.Duration > 0 && time.Now().After(endDate) {
			continue
		}

		if s.IsPermanent {
			result = append(result, map[string]string{
				"medication_name": s.MedicationName,
				"time":            "Постоянно",
			})
			continue
		}

		intakeTimes := utils.CalculateIntakeTimes(s.Frequency)

		for _, t := range intakeTimes {
			parsedTime, err := time.Parse("15:04", t)
			if err != nil {
				continue
			}

			if parsedTime.After(currentTime) && parsedTime.Before(nextTime) {
				result = append(result, map[string]string{
					"medication_name": s.MedicationName,
					"time":            t,
				})
			}
		}
	}

	if len(result) == 0 {
		return nil, errors.New("no upcoming takings found")
	}

	return result, nil
}
