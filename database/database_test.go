package database_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"ninja-swords/database"
	"ninja-swords/logging"
	"ninja-swords/models"
)

type Animal struct {
	ID   uint `gorm:"primarykey;autoIncrement"`
	Name string
	Age  int64
}

func testDatabase(t *testing.T, opt database.ConnectOption) {
	assert := assert.New(t)
	logging.TestingInitialize()
	defer logging.TestingFinalize()

	database.TestingInitialize(database.Default, opt)
	database.AutoMigrate(database.Default, models.Models())
	defer database.TestingFinalize()

	conn := database.GetDB(database.Default)

	err := conn.Migrator().CreateTable(Animal{})
	assert.NoError(err)

	testCreateAndFirst(t, conn)
	testNestedTransaction(t, conn)
	testWhereOr(t, conn)
}

func testCreateAndFirst(t *testing.T, conn *gorm.DB) {
	assert := assert.New(t)
	var err error

	animal := Animal{ID: 999, Name: "Bear", Age: 33}
	err = conn.Create(&animal).Error
	assert.NoError(err)

	var bear Animal
	err = conn.First(&bear).Error
	assert.NoError(err)

	assert.Equal(animal, bear)
}

func testNestedTransaction(t *testing.T, conn *gorm.DB) {
	assert := assert.New(t)
	var err error

	var deer99 Animal
	err = conn.Transaction(func(tx *gorm.DB) error {
		// create one.
		var deer Animal
		ierr := tx.Transaction(func(tx1 *gorm.DB) error {
			deer.Name = "Deer"
			deer.Age = 100
			return tx1.Create(&deer).Error
		})

		if ierr != nil {
			return ierr
		}

		assert.NotEqual(uint(0), deer.ID)

		// update one.
		ierr = tx.Transaction(func(tx2 *gorm.DB) error {
			return tx2.Model(&Animal{}).Where("id = ?", deer.ID).Update("age", 99).Error
		})

		if ierr != nil {
			return ierr
		}

		deer99.ID = deer.ID
		return nil
	})
	assert.NoError(err)

	err = conn.Where(&deer99).First(&deer99).Error
	assert.NoError(err)

	assert.Equal(deer99.Age, int64(99))
}

func testWhereOr(t *testing.T, conn *gorm.DB) {
	assert := assert.New(t)
	var err error

	animal := Animal{ID: 888, Name: "Monkey", Age: 77}
	err = conn.Create(&animal).Error
	assert.NoError(err)

	var monkey Animal
	err = conn.Or("name = ?", "Monkey").First(&monkey).Error
	assert.NoError(err)

	assert.Equal(animal, monkey)

	var monkey2 Animal
	err = conn.Where("age = ?", 44).Or("name = ?", "Monkey").First(&monkey2).Error
	assert.NoError(err)

	assert.Equal(animal, monkey2)

	var monkey3 Animal
	err = conn.Where("age = ?", 44).Or("name = ?", "Monkey").Where("id = ?", 888).First(&monkey3).Error
	assert.NoError(err)

	assert.Equal(animal, monkey3)
}

func TestDatabase(t *testing.T) {
	testDatabase(t, database.SQLiteOpt)
}
