package lessondb

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/owezzy/schoolCRM/business/domain/lessonbus"
)

func (s *Store) Creation(ctx context.Context, actorID, requestID uuid.UUID) (lessonbus.Creation, error) {
	// Same actor/request commands serialize across service instances. Hash
	// collisions merely serialize unrelated commands; the primary key still identifies them.
	if _, err := s.tx.ExecContext(ctx, `SELECT pg_advisory_xact_lock(hashtextextended($1,0))`, "go-lesson-creation:"+actorID.String()+":"+requestID.String()); err != nil {
		return lessonbus.Creation{}, err
	}
	var row struct {
		InputHash string `db:"input_hash"`
		Result    []byte `db:"result"`
	}
	if err := s.tx.GetContext(ctx, &row, `SELECT input_hash,result FROM lesson_creation_requests
        WHERE actor_id=$1 AND request_id=$2`, actorID, requestID); err != nil {
		return lessonbus.Creation{}, translate(err)
	}
	creation := lessonbus.Creation{ActorID: actorID, RequestID: requestID, InputHash: row.InputHash}
	if err := json.Unmarshal(row.Result, &creation.Result); err != nil {
		return lessonbus.Creation{}, fmt.Errorf("decode creation receipt: %w", err)
	}
	return creation, nil
}

func (s *Store) RecordCreation(ctx context.Context, creation lessonbus.Creation) error {
	result, err := json.Marshal(creation.Result)
	if err != nil {
		return err
	}
	_, err = s.tx.ExecContext(ctx, `INSERT INTO lesson_creation_requests(actor_id,request_id,input_hash,result)
        VALUES ($1,$2,$3,$4)`, creation.ActorID, creation.RequestID, creation.InputHash, result)
	return err
}
