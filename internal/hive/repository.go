package hive

import (
	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/gaetanDubuc/beeckend/pkg/repository"
)

func NewRepository() *repository.Repository[entity.Hive] {
	return repository.NewRepository[entity.Hive]()
}
