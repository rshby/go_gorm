package testing_v2

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go_gorm/config"
	"go_gorm/infrastructure/database/connection"
	"go_gorm/model/entity"
	mck "go_gorm/testing/mock/config"
	"gorm.io/gorm"
	"log"
	"strconv"
	"testing"
)

var cfg = config.Config{
	App: nil,
	Database: &config.Database{
		User:     "root",
		Password: "root",
		Host:     "localhost",
		Port:     3306,
		Name:     "belajar_golang_gorm",
	},
}

func TestConnectionDB(t *testing.T) {
	var configMock = mck.NewConfigMock()

	// mock
	configMock.Mock.On("GetConfig").Return(&cfg)
	var db *gorm.DB = connection.ConnectToDB(configMock)
	assert.NotNil(t, db)
}

// TestRawQuery adalah function untuk melakukan execute query dengan raw query
// menggunakan db.Exec() dan db.ExecContext()
func TestRawSQLExec(t *testing.T) {
	configMock := mck.NewConfigMock()
	configMock.Mock.On("GetConfig").Return(&cfg)
	db := connection.ConnectToDB(configMock)

	err := db.Exec("INSERT INTO sample(id, name) VALUES (?, ?);", "1", "eko").Error
	assert.Nil(t, err)

	err = db.Exec("INSERT INTO sample(id, name) VALUES (?, ?);", "2", "budi").Error
	assert.Nil(t, err)

	err = db.Exec("INSERT INTO sample(id, name) VALUES (?, ?);", "3", "joko").Error
	assert.Nil(t, err)

	err = db.Exec("INSERT INTO sample(id, name) VALUES (?, ?);", "4", "rully").Error
	assert.Nil(t, err)
}

// TestRAWSQLQuery adalah function untuk melakukan query select dengan raw sql
// menggunakan db.Queqry() dan db.QueryContext()
func TestRAWSQLQuery(t *testing.T) {
	configMock := mck.NewConfigMock()
	configMock.Mock.On("GetConfig").Return(&cfg)
	db := connection.ConnectToDB(configMock)

	type Sample struct {
		Id   string `json:"id,omitempty" gorm:"column=id"`
		Name string `json:"nama,omitempty" gorm:"column=name"`
	}

	scenario := []struct {
		Nama        string
		id          string
		ExpectError bool
	}{
		{
			Nama:        "test select data eko success",
			id:          "1",
			ExpectError: false,
		},
		{
			Nama:        "test select data budi success",
			id:          "2",
			ExpectError: false,
		},
	}

	for _, testScenario := range scenario {
		t.Run(testScenario.Nama, func(t *testing.T) {
			// select data
			var sampe Sample
			err := db.Raw("SELECT id, name FROM sample WHERE id=?;", testScenario.id).Scan(&sampe).Error
			assert.Equal(t, err != nil, testScenario.ExpectError)

			sampleJson, _ := json.Marshal(&sampe)
			log.Println(string(sampleJson))
		})
	}

	// test get all data
	t.Run("test get all data samples", func(t *testing.T) {
		var samples []Sample
		err := db.Raw("SELECT id, name FROM sample;").Scan(&samples).Error
		if err == nil {
			sampleJson, _ := json.Marshal(&samples)
			log.Println(string(sampleJson))
		}

		assert.Nil(t, err)
		assert.Equal(t, 4, len(samples))
	})

}

// TestQueryRows adalah function untuk melakukan query dengan method Rows()
// returnnya adalah (*sql.rows, error)
// dilakukan scan satu-satu sesuai dengan kolom yang diSelect
func TestQueryRows(t *testing.T) {
	configMock := mck.NewConfigMock()
	configMock.Mock.On("GetConfig").Return(&cfg)
	db := connection.ConnectToDB(configMock)

	type Sample struct {
		Id   string `json:"id,omitempty" gorm:"column=id"`
		Name string `json:"name,omitempty" gorm:"column=name"`
	}

	scenario := []struct {
		Name        string
		Id          string
		ExpectError bool
	}{
		{
			Name:        "test get rows id 1 success",
			Id:          "1",
			ExpectError: false,
		},
		{
			Name:        "test get row id 4 success",
			Id:          "4",
			ExpectError: false,
		}}

	// looping each test case
	for _, testScenario := range scenario {
		t.Run(testScenario.Name, func(t *testing.T) {
			var sample Sample
			rows, err := db.Raw("SELECT id,sample.name FROM sample WHERE id=?;", testScenario.Id).Rows()
			assert.Equal(t, err != err, testScenario.ExpectError)
			if err != nil {
				panic(err)
			}

			defer rows.Close()

			if rows.Next() {
				if err := rows.Scan(&sample.Id, &sample.Name); err != nil {
					panic(err)
				}
			}

			sampleJson, _ := json.Marshal(&sample)
			log.Println(string(sampleJson))
		})
	}

	// test menggunakan db.ScanRows()
	// jika menggunakan Scan() maka harus satu-satu menuliskan field tujuannya
	// jika menggunakan ScanRows() maka tidak satu-satu, namun langsung hanya menuliskan variabelnya
	t.Run("test get all data sample using ScanRows", func(t *testing.T) {
		rows, err := db.Raw("SELECT id, name FROM sample;").Rows()
		assert.Nil(t, err)
		defer rows.Close()

		var samples []Sample
		for rows.Next() {
			err = db.ScanRows(rows, &samples)
			assert.Nil(t, err)
		}

		// convert to json
		sampleJson, _ := json.Marshal(&samples)
		log.Println(string(sampleJson))

		assert.Equal(t, 4, len(samples))
	})
}

// TestCreateUser adalah function untuk insert new data user
// menggunakan db.Create()
func TestCreateUser(t *testing.T) {
	configMock := mck.NewConfigMock()
	configMock.Mock.On("GetConfig").Return(&cfg)
	db := connection.ConnectToDB(configMock)

	scenario := []struct {
		Name        string
		Input       entity.User
		ExpectError bool
	}{
		{
			Name: "test create user id 1",
			Input: entity.User{
				ID:       "1",
				Password: "rahasia",
				Name: entity.Name{
					FirstName:  "Eko",
					MiddleName: "Kurniawan",
					LastName:   "Khannedy",
				},
				Information: "ini akan diignore",
			},
			ExpectError: false,
		},
	}

	for _, testCase := range scenario {
		t.Run(testCase.Name, func(t *testing.T) {
			result := db.Create(&testCase.Input)
			assert.Equal(t, 1, int(result.RowsAffected))
			assert.Equal(t, result.Error != nil, testCase.ExpectError)
		})
	}
}

// TestBatchInsert adalah function untuk insert banyak rows sekaligus
// apabila Create() hanya akan melakukan insert satu data saja
// untuk bisa insert banyak data, tetap meggunakan Create() tapi parameternya slice []entity.struct
func TestBatchInsert(t *testing.T) {
	configMock := mck.NewConfigMock()
	configMock.Mock.On("GetConfig").Return(&cfg)
	db := connection.ConnectToDB(configMock)

	// create data
	var users []entity.User
	for i := 0; i < 10; i++ {
		users = append(users, entity.User{
			ID:       strconv.Itoa(i + 2),
			Password: "rahasia",
			Name: entity.Name{
				FirstName: fmt.Sprintf("user ke %v", i+2),
			},
			Information: "ini akan diignore",
		})
	}

	// create scenario testing
	scenario := []struct {
		Name       string
		Input      []entity.User
		ExpecError bool
	}{
		{
			Name:       "test insert batch succes",
			Input:      users,
			ExpecError: false,
		},
	}

	for _, testCase := range scenario {
		t.Run(testCase.Name, func(t *testing.T) {
			result := db.Create(&testCase.Input)
			assert.Equal(t, len(users), int(result.RowsAffected))
			assert.Equal(t, result.Error != nil, testCase.ExpecError)
		})
	}
}
