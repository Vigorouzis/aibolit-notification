package postgres

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/vigorouzis/aibolit-notification/internal/core/domain"
)

type Client struct {
	db *sql.DB
}

func New(db *sql.DB) *Client {
	return &Client{
		db: db,
	}
}

func (c *Client) CreateSchedule(ctx context.Context, schedule *domain.Schedule) error {
	query, args, err := sq.Insert("schedules").
		Columns("id", "user_id", "medication_name", "frequency", "duration", "is_permanent", "start_date").
		Values(schedule.ID, schedule.UserID, schedule.MedicationName, schedule.Frequency, schedule.Duration, schedule.IsPermanent,
			schedule.StartDate).
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return err
	}

	_, err = c.db.ExecContext(ctx, query, args...)
	return err
}

func (r *Client) GetSchedulesByUserID(ctx context.Context, userID string) ([]domain.Schedule, error) {
	query, args, err := sq.Select("id", "user_id", "medication_name", "frequency", "duration", "is_permanent", "start_date").
		From("schedules").
		Where(sq.Eq{"user_id": userID}).PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []domain.Schedule
	for rows.Next() {
		var schedule domain.Schedule
		if err := rows.Scan(&schedule.ID, &schedule.UserID, &schedule.MedicationName, &schedule.Frequency, &schedule.Duration, &schedule.IsPermanent, &schedule.StartDate); err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}

	if len(schedules) == 0 {
		return nil, sql.ErrNoRows
	}

	return schedules, nil
}

func (c *Client) GetSchedule(ctx context.Context, userID, scheduleID string) (*domain.Schedule, error) {

	query, args, err := sq.Select("id", "user_id", "medication_name", "frequency", "duration", "is_permanent", "start_date").
		From("schedules").
		Where(sq.Eq{"id": scheduleID, "user_id": userID}).
		PlaceholderFormat(sq.Dollar).ToSql()

	if err != nil {
		return nil, err
	}

	var schedule domain.Schedule
	err = c.db.QueryRowContext(ctx, query, args...).Scan(
		&schedule.ID, &schedule.UserID, &schedule.MedicationName, &schedule.Frequency, &schedule.Duration, &schedule.IsPermanent,
		&schedule.StartDate,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &schedule, nil
}
