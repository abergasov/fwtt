package validator

import (
	"fmt"
	"fwtt/internal/entites"
	"fwtt/internal/storage/database"

	"github.com/Masterminds/squirrel"
)

const (
	table = "challenges"
)

type Repository struct {
	db database.DBConnector
}

func NewRepository(db database.DBConnector) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) LoadChallenges() ([]*entites.Challenge, error) {
	sql := fmt.Sprintf("SELECT valid_till, challenge, difficulty, max_allowed, used, hash_algo, hash FROM %s", table)
	rows, err := r.db.Client().Queryx(sql)
	if err != nil {
		return nil, fmt.Errorf("failed to load challenges: %w", err)
	}
	defer rows.Close()
	result := make([]*entites.Challenge, 0, 10000)
	for rows.Next() {
		var ch entites.Challenge
		if err = rows.StructScan(&ch); err != nil {
			return nil, fmt.Errorf("failed to scan challenge: %w", err)
		}
		result = append(result, &ch)
	}
	return result, nil
}

func (r *Repository) SaveChallenges(results []*entites.Challenge) error {
	q := squirrel.Insert(table).Columns("valid_till", "challenge", "difficulty", "max_allowed", "used", "hash_algo", "hash")
	for _, rs := range results {
		q = q.Values(rs.ValidTill, rs.Challenge, rs.Difficulty, rs.MaxAllowed, rs.Used, rs.HashAlgo, rs.Hash)
	}
	sql, args, err := q.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	_, err = r.db.Client().Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}

func (r *Repository) DropChallenges(challenges []string) error {
	if len(challenges) == 0 {
		return nil
	}

	q := squirrel.Delete(table).Where(squirrel.Eq{"challenge": challenges})
	sql, args, err := q.PlaceholderFormat(squirrel.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}
	_, err = r.db.Client().Exec(sql, args...)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}
	return nil
}
