package testing

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go_gorm/model/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"testing"
)

func SetupDb() *gorm.DB {
	dsn := "root:root@tcp(localhost:3306)/belajar_golang_gorm?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal()
	}

	// success create instance db
	return db
}

// test untuk insert ke database
func TestExecuteSQL(t *testing.T) {
	db := SetupDb()
	assert.NotNil(t, db)

	t.Run("insert one data", func(t *testing.T) {
		err := db.Exec("insert into sample (id, name) values (?, ?)", "1", "Eko").Error
		assert.Nil(t, err)

		err = db.Exec("insert into sample(id, name) values (?, ?)", "2", "Budi").Error
		assert.Nil(t, err)

		err = db.Exec("insert into sample(id, name) values(?, ?)", "3", "Joko").Error
		assert.Nil(t, err)

		err = db.Exec("insert into sample(id, name) values(?, ?)", "4", "Rully").Error
		assert.Nil(t, err)
	})
}

// test untuk query ke database
func TestQuerySelect(t *testing.T) {
	db := SetupDb()

	// test select by id
	t.Run("test select sample", func(t *testing.T) {
		var sample = new(entity.Sample)
		err := db.Raw("select id, name from sample where id = ?", "1").Scan(&sample).Error
		assert.Nil(t, err)
		assert.Equal(t, "1", sample.Id)
	})

	// test select all data
	t.Run("select all data", func(t *testing.T) {
		var sample []entity.Sample
		err := db.Raw("select id, name from sample").Scan(&sample).Error
		assert.Nil(t, err)
		assert.Equal(t, 4, len(sample))

		sampleJson, _ := json.Marshal(&sample)
		fmt.Println(string(sampleJson))
	})

	// test select all data using sql.Rows -> looping
	t.Run("select all data using Rows", func(t *testing.T) {
		rows, err := db.Raw("select id, name from sample").Rows()
		if err != nil {
			t.Fatalf(err.Error())
		}
		defer rows.Close()

		var sample []entity.Sample
		for rows.Next() {
			var s entity.Sample
			if err := rows.Scan(&s.Id, &s.Name); err != nil {
				t.Fatalf(err.Error())
			}

			sample = append(sample, s)
		}

		assert.Equal(t, 4, len(sample))
	})

	// test select all using gorm.ScanRows
	t.Run("select all data using scanRows", func(t *testing.T) {
		rows, err := db.Raw("select id, name from sample").Rows()
		assert.Nil(t, err)
		assert.NotNil(t, rows)
		if err != nil {
			t.Fatalf(err.Error())
		}

		defer rows.Close()

		var sample []entity.Sample
		for rows.Next() {
			err := db.ScanRows(rows, &sample)
			assert.Nil(t, err)
		}

		assert.Equal(t, 4, len(sample))
	})
}
