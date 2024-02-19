package user

import (
	"fmt"
	"os"
	"testing"

	"github.com/gaetanDubuc/beeckend/internal/entity"
	"github.com/magiconair/properties/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestRepository_Get(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("user.db"), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&entity.User{})

	// Create
	repo := NewRepository(db)

	user := entity.User{
		Name:  "John",
		Email: "r@f.fr",
	}
	use, err := repo.Create(user)
	fmt.Println(use, err)
	repouser, err := repo.Get(1)

	assert.Equal(t, repouser.Name, "John")

	if err := os.Remove("user.db"); err != nil {
		t.Errorf("Erreur lors de la suppression du fichier gorm.db : %v", err)
	}
}
