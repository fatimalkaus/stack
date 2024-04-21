package stack

import (
	"context"
	"errors"

	"gorm.io/gorm"
)

type item struct {
	ID    uint `gorm:"primaryKey"`
	Value string
}

func (item) TableName() string {
	return "stack"
}

// PostgresStack implements Stack using Postgres.
type PostgresStack struct {
	db *gorm.DB
}

// NewPostgresStack creates new PostgresStack.
func NewPostgresStack(g *gorm.DB) (Stack, error) {
	return &PostgresStack{db: g}, nil
}

// Push pushes data to stack.
func (s *PostgresStack) Push(ctx context.Context, val any) error {
	if err := s.db.WithContext(ctx).Create(&item{Value: val.(string)}).Error; err != nil {
		return err
	}
	return nil
}

// Pop pops data from stack.
func (s *PostgresStack) Pop(ctx context.Context) (any, error) {
	var value interface{}
	if err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var item item
		if err := tx.Order("id desc").First(&item).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrIsEmpty
			}
			return err
		}

		if err := tx.Delete(&item).Error; err != nil {
			return err
		}

		value = item.Value
		return nil
	}); err != nil {
		return nil, err
	}

	return value, nil
}

// Peek peeks data from stack.
func (s *PostgresStack) Peek(ctx context.Context) (any, error) {
	var item item
	if err := s.db.WithContext(ctx).Order("id desc").First(&item).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrIsEmpty
		}
		return nil, err
	}

	return item.Value, nil
}
