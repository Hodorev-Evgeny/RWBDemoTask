package core_domain

import "time"

type SearchEvent struct {
	Query     string    `json:"query"`
	UserID    int64     `json:"user_id"`
	SessionID string    `json:"session_id"`
	TimeEvent time.Time `json:"time_event"`
}

func CreateNewEvent(
	query string,
	userID int64,
	sessionID string,
	timeEvent time.Time,
) SearchEvent {
	return SearchEvent{
		Query:     query,
		UserID:    userID,
		SessionID: sessionID,
		TimeEvent: timeEvent,
	}
}
