package repository

import (
	"fmt"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/kerti/evm/04-tennis-player/database"
	"github.com/kerti/evm/04-tennis-player/model"
	"github.com/kerti/evm/04-tennis-player/util/failure"
	"github.com/kerti/evm/04-tennis-player/util/logger"
)

const (
	queryInsertContainer = `
		INSERT INTO containers (
			containers.entity_id,
			containers.player_entity_id,
			containers.capacity,
			containers.ball_count
		) VALUES (
			:entity_id,
			:player_entity_id,
			:capacity,
			:ball_count)`

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
	ExistsByID(id uuid.UUID) (exists bool, err error)
	Create(container model.Container) (err error)
	ResolveByIDs(ids []uuid.UUID) (containers []model.Container, err error)
	ResolveByPlayerID(playerID uuid.UUID) (containers []model.Container, err error)
	ResolvePage(pageNum int, pageSize int) (page *model.Page, err error)
	TxBulkUpdate(tx *sqlx.Tx, containers []model.Container) (err error)
}

// ContainerMySQLRepo is the repository for Containers implemented with MySQL backend
type ContainerMySQLRepo struct {
	DB *database.MySQL `inject:"mysql"`
}

// Startup perform startup functions
func (r *ContainerMySQLRepo) Startup() {
	logger.Trace("Container repository starting up...")
}

// Shutdown cleans up everything and shuts down
func (r *ContainerMySQLRepo) Shutdown() {
	logger.Trace("Container repository shutting down...")
}

// ExistsByID checks whether a Container exists by its ID
func (r *ContainerMySQLRepo) ExistsByID(id uuid.UUID) (exists bool, err error) {
	err = r.DB.Get(
		&exists,
		"SELECT COUNT(entity_id) > 0 FROM containers WHERE containers.entity_id = ?",
		id.String())
	if err != nil {
		logger.ErrNoStack("%v", err)
	}
	return
}

// Create creates a new Container
func (r *ContainerMySQLRepo) Create(container model.Container) (err error) {
	exists, err := r.ExistsByID(container.ID)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	if exists {
		err = failure.OperationNotPermitted("create", "Container", "already exists")
		logger.ErrNoStack("%v", err)
		return err
	}

	stmt, err := r.DB.Prepare(queryInsertContainer)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	_, err = stmt.Exec(container)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return err
	}

	return nil
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

// ResolveByPlayerID resolves Containers by their Player IDs
func (r *ContainerMySQLRepo) ResolveByPlayerID(playerID uuid.UUID) (containers []model.Container, err error) {
	query, args, err := r.DB.In(
		querySelectContainer+" WHERE containers.player_entity_id = ?",
		playerID)
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

// ResolvePage resolves a Page of Containers based on page and page size parameters
func (r *ContainerMySQLRepo) ResolvePage(pageNum int, pageSize int) (page *model.Page, err error) {
	offset := (pageNum - 1) * pageSize
	query, args, err := r.DB.In(
		querySelectContainer+" LIMIT ? OFFSET ?",
		pageSize,
		offset,
	)
	if err != nil {
		logger.ErrNoStack("%v", err)
		return
	}

	containers := make([]model.Container, 0)
	err = r.DB.Select(&containers, query, args...)
	if err != nil {
		logger.ErrNoStack("%v", err)
	}

	var count int
	err = r.DB.Get(&count, "SELECT COUNT(entity_id) FROM containers")
	if err != nil {
		return nil, err
	}

	page = &model.Page{
		Items:      containers,
		Page:       pageNum,
		PageSize:   pageSize,
		TotalCount: count,
	}
	page.CalculateTotalPages()
	return page, nil
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

	query, args, err := r.DB.NamedIn(q, param)
	if err != nil {
		return query, params, err
	}

	params = append(params, args...)
	return
}

func (r *ContainerMySQLRepo) buildCaseWhenClause(fieldName string, elements []string) string {
	return fmt.Sprintf("%s = (CASE containers.entity_id %s END)", fieldName, strings.Join(elements, " "))
}
