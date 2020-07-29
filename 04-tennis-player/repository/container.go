package repository

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kerti/evm/04-tennis-player/database"
	"github.com/kerti/evm/04-tennis-player/model"
	"github.com/kerti/evm/04-tennis-player/util/logger"
)

const (
	querySelectContainer = `
		SELECT
			containers.entity_id,
			containers.player_entity_id,
			containers.capacity,
			containers.ball_count
		FROM containers`
)

// Container is the Container repository interface
type Container interface {
	Startup()
	Shutdown()
	ResolveByIDs(ids []uuid.UUID) (containers []model.Container, err error)
	TxBulkUpdate(tx *sqlx.Tx, containers []model.Container) (err error)
}

// ContainerMySQLRepo is the repository for Containers implemented with MySQL backend
type ContainerMySQLRepo struct {
	DB *database.MySQL `inject:"mysql"`
}

// Startup perform startup functions
func (r *ContainerMySQLRepo) Startup() {
	logger.Trace("BankAccount repository starting up...")
}

// Shutdown cleans up everything and shuts down
func (r *ContainerMySQLRepo) Shutdown() {
	logger.Trace("BankAccount repository shutting down...")
}

// ResolveByIDs resolves Containers by their IDs
func (r *ContainerMySQLRepo) ResolveByIDs(ids []uuid.UUID) (containers []model.Container, err error) {
	if len(ids) == 0 {
		return
	}

	query, args, err := r.DB.In(querySelectContainer+" WHERE containers.entity_id IN (?)", ids)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return
	}

	err = r.DB.Select(&containers, query, args...)
	if err != nil {
		logger.ErrNoStack("%v", err)
	}

	return
}

// TxBulkUpdate transactionally updates multiple containers with the transaction object passed from elsewhere
func (r *ContainerMySQLRepo) TxBulkUpdate(tx *sqlx.Tx, containers []model.Container) (err error) {
	if len(containers) == 0 {
		return nil
	}

	query, args, err := r.composeBulkUpdateQuery(containers)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return
	}

	stmt, err := tx.Preparex(query)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	_, err = stmt.Stmt.Exec(args...)
	if err != nil {
		logger.ErrNoStack("failed running bulk update container query: %v, args: %v", query, args)
		return err
	}

	return nil
}

func (r *ContainerMySQLRepo) composeBulkUpdateQuery(containers []model.Container) (query string, params []interface{}, err error) {
	param := map[string]interface{}{}

	entityIDs := make([]string, 0)
	playerIDs := make([]string, 0)
	capacities := make([]string, 0)
	ballCounts := make([]string, 0)

	for idx, container := range containers {
		param[fmt.Sprintf("entity_id_%d", idx)] = container.ID
		param[fmt.Sprintf("player_entity_id_%d", idx)] = container.PlayerID
		param[fmt.Sprintf("capacity_%d", idx)] = container.Capacity
		param[fmt.Sprintf("ball_count_%d", idx)] = container.BallCount

		entityIDs = append(entityIDs, fmt.Sprintf(":entity_id_%d", idx))
		playerIDs = append(playerIDs, fmt.Sprintf("WHEN :entity_id_%d THEN :player_entity_id_%d", idx, idx))
		capacities = append(capacities, fmt.Sprintf("WHEN :entity_id_%d THEN :capacity_%d", idx, idx))
		ballCounts = append(ballCounts, fmt.Sprintf("WHEN :entity_id_%d THEN :ball_count_%d", idx, idx))
	}

	fieldClauses := make([]string, 0)
	fieldClauses = append(fieldClauses, r.buildCaseWhenClause("player_entity_id", playerIDs))
	fieldClauses = append(fieldClauses, r.buildCaseWhenClause("capacity", capacities))
	fieldClauses = append(fieldClauses, r.buildCaseWhenClause("ball_count", ballCounts))

	q := fmt.Sprintf(
		"UPDATE containers SET %s WHERE containers.entity_id IN (%s)",
		strings.Join(fieldClauses, ", "),
		strings.Join(entityIDs, ", "))

	query, args, err := r.DB.In(q, param)
	if err != nil {
		return query, params, err
	}
	params = append(params, args...)

	return
}

func (r *ContainerMySQLRepo) buildCaseWhenClause(fieldName string, elements []string) string {
	return fmt.Sprintf("%s = (CASE containers.entity_id %s END)", fieldName, strings.Join(elements, " "))
}
