package http

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vigorouzis/aibolit-notification/internal/core/domain"
)

func (s *Server) createSchedule(c *gin.Context) {
	ctx := c.Request.Context()

	var req struct {
		UserID         string `json:"user_id"`
		MedicationName string `json:"medication_name"`
		Frequency      int    `json:"frequency"`
		Duration       int    `json:"duration"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	schedule := domain.NewSchedule(
		req.UserID,
		req.MedicationName,
		req.Frequency,
		req.Duration,
	)

	scheduleID, err := s.srv.CreateSchedule(ctx, schedule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create schedule"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"schedule_id": scheduleID})
}

func (s *Server) getSchedulesByUserID(c *gin.Context) {
	ctx := c.Request.Context()

	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	schedules, err := s.srv.GetSchedulesByUserId(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch schedules"})
		return
	}

	c.JSON(http.StatusOK, schedules)
}

func (s *Server) getSchedule(c *gin.Context) {
	userID := c.Query("user_id")
	scheduleID := c.Query("schedule_id")

	if userID == "" || scheduleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id and schedule_id are required"})
		return
	}

	schedule, intakeTimes, err := s.srv.GetSchedule(context.Background(), userID, scheduleID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Schedule not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"schedule":     schedule,
		"intake_times": intakeTimes,
	})
}

func (s *Server) nextTakings(c *gin.Context) {
	userID := c.Query("user_id")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	takings, err := s.srv.NextTakings(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"next_takings": takings})
}
