package audit

import (
	"context"
	"database/sql"
	"encoding/json"
)

type Event struct {
	OrganizationID string
	ActorUserID    string
	Action         string
	EntityType     string
	EntityID       string
	RequestID      string
	Metadata       map[string]string
}

type Writer interface {
	Write(context.Context, Event) error
}

type SQLWriter struct {
	DB *sql.DB
}

func (w SQLWriter) Write(ctx context.Context, e Event) error {
	data, err := json.Marshal(e.Metadata)
	if err != nil {
		return err
	}
	_, err = w.DB.ExecContext(ctx, "INSERT INTO audit_events(organization_id, actor_user_id, action, entity_type, entity_id, metadata, request_id) VALUES(NULLIF($1,'')::uuid, NULLIF($2,'')::uuid, $3, $4, NULLIF($5,'')::uuid, $6, $7)", e.OrganizationID, e.ActorUserID, e.Action, e.EntityType, e.EntityID, data, e.RequestID)
	return err
}

type NopWriter struct{}

func (NopWriter) Write(context.Context, Event) error { return nil }
