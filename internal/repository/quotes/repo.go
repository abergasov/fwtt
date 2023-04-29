package quotes

import (
	"fmt"
	"fwtt/internal/entites"
	"fwtt/internal/storage/database"
)

const (
	table = "quotes"
)

type Repository struct {
	db database.DBConnector
}

func NewRepository(db database.DBConnector) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) LoadQuotes() ([]*entites.Quote, error) {
	rows, err := r.db.Client().Queryx(fmt.Sprintf("SELECT q_id, quote, by FROM %s", table))
	if err != nil {
		return nil, fmt.Errorf("failed to load quotes: %w", err)
	}
	defer rows.Close()
	result := make([]*entites.Quote, 0, 10)
	for rows.Next() {
		var q entites.Quote
		if err = rows.StructScan(&q); err != nil {
			return nil, fmt.Errorf("failed to scan quote: %w", err)
		}
		result = append(result, &q)
	}
	return result, nil
}
