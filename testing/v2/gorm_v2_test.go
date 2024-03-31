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

// TestTransactionGorm adalah function untuk menggunakan transaction sql
func TestTransactionGorm(t *testing.T) {
	configMock := mck.NewConfigMock()
	configMock.Mock.On("GetConfig").Return(&cfg)
	db := connection.ConnectToDB(configMock)

	// function transaction
	// menggunakan db.Begin(), kemudian menggunakan tx.Commit() atau tx.Rollback()
	// harus ditulis Commit() dan Rollback() secara manual
	// test transaction success -> commit
	t.Run("test transaction menggunakan begin rollback", func(t *testing.T) {
		// create transaction
		tx := db.Begin()
		defer tx.Rollback()

		// create user
		err := tx.Create(&entity.User{
			ID:       "12",
			Password: "rahasia",
			Name: entity.Name{
				FirstName: "user ke 12",
			},
			Information: "",
		}).Error
		assert.Nil(t, err)

		if err != nil {
			return
		}

		// commit transaction
		tx.Commit()
	})

	// function transaction
	// menggunakan method db.Transaction(), menuliskan code di dalam function callbacknya
	// test transaction gagal -> rollback
	t.Run("test transaction menggunakan db.Transaction callback", func(t *testing.T) {
		err := db.Transaction(func(tx *gorm.DB) error {
			user := entity.User{
				ID:       "13",
				Password: "rahasia",
				Name: entity.Name{
					FirstName: "user ke 13",
				},
				Information: "ini akan diignore",
			}

			err := tx.Create(&user).Error
			assert.Nil(t, err)
			if err != nil {
				return err
			}

			err = tx.Create(&user).Error
			assert.NotNil(t, err)
			assert.Error(t, err)
			if err != nil {
				return err
			}

			// will automaticly commit
			return nil
		})

		assert.NotNil(t, err)
		assert.Error(t, err)
	})
}

// TestQueryGormSingleRow untuk query data yang hasilnya satu data saja
// dapat menggunakan menthod db.First(), db.Take(), db.Last()
func TestQueryGormSingleRow(t *testing.T) {
	configMock := mck.NewConfigMock()
	configMock.Mock.On("GetConfig").Return(&cfg)
	db := connection.ConnectToDB(configMock)

	// db.First() digunakan untuk query single data yang diurutkan dari id paling kecil
	t.Run("test query single object using First method", func(t *testing.T) {
		var user entity.User
		err := db.First(&user).Error
		assert.Nil(t, err)
		assert.Equal(t, user.ID, "1")

		userJson, _ := json.Marshal(&user)
		log.Println(string(userJson))
	})

	// db.First() dengan WHERE filter -> WHERE id='5'
	t.Run("test query single object using First method with WHERE filter", func(t *testing.T) {
		var user entity.User
		err := db.First(&user, "id = ?", "5").Error
		assert.Nil(t, err)
		assert.Equal(t, "5", user.ID)

		userJson, _ := json.Marshal(&user)
		log.Println(string(userJson))
	})

	// db.Take() digunakan untuk query single data yang datanya tidak diurutkan
	t.Run("test query single object using Take method", func(t *testing.T) {
		var user entity.User
		err := db.Take(&user).Error
		assert.Nil(t, err)

		userJson, _ := json.Marshal(&user)
		log.Println(string(userJson))
	})

	// db.Take() dengan WHERE filter -> WHERE id='5'
	t.Run("test query single object using Take method with WHERE filter", func(t *testing.T) {
		var user entity.User
		err := db.Take(&user, "id = ?", "5").Error
		assert.Nil(t, err)
		assert.Equal(t, "5", user.ID)

		userJson, _ := json.Marshal(&user)
		log.Println(string(userJson))
	})

	// db.Last() digunakan untuk query single data yang datanya diurutkan dari id terakhir
	t.Run("test query single object using Last method", func(t *testing.T) {
		var user entity.User
		err := db.Last(&user).Error
		assert.Nil(t, err)
		assert.Equal(t, "9", user.ID)

		userJson, _ := json.Marshal(&user)
		log.Println(string(userJson))
	})
}

// TestQueryGormAllRows untuk query data yang hasilnya banyak data / lebih dari satu data
// dapat menggunakan method db.Find()
// dapat juga ditambahkan query WHERE dengan condition
func TestQueryGormAllRows(t *testing.T) {
	configMock := mck.NewConfigMock()
	configMock.Mock.On("GetConfig").Return(&cfg)
	db := connection.ConnectToDB(configMock)

	// query all rows using method db.Find()
	t.Run("test query all rows using method Find", func(t *testing.T) {
		var users []entity.User
		err := db.Find(&users).Error
		assert.Nil(t, err)

		usersJson, _ := json.Marshal(&users)
		log.Println(string(usersJson))
	})

	// db.Find() dengan WHERE filter -> WHERE id IN ('5', '6')
	t.Run("test query all rows using method Find dengan WHERE filter", func(t *testing.T) {
		var users []entity.User
		result := db.Find(&users, "id IN (?)", []string{"5", "6"})
		assert.Nil(t, result.Error)
		assert.Equal(t, 2, len(users))

		usersJson, _ := json.Marshal(&users)
		log.Println(string(usersJson))
	})
}

// gorm advance query
// TestQueryGormWhere adalah function untuk query menggunakan WHERE method
// method db.Where() akan dipanggil sebelum memanggil method Take(), atau Find()
func TestQueryGormWhereCondition(t *testing.T) {
	configMock := mck.NewConfigMock()
	configMock.Mock.On("GetConfig").Return(&cfg)
	db := connection.ConnectToDB(configMock)

	// test WHERE dua kondisi menggunakan 2 kali pemanggilan method WHERE -> akan dianggap AND
	t.Run("test query WHERE with find, multi WHERE", func(t *testing.T) {
		users := []entity.User{}
		err := db.Where("first_name LIKE ?", "user%").Where("password = ?", "rahasia").Find(&users).Error
		assert.Nil(t, err)

		usersJson, _ := json.Marshal(&users)
		log.Println(string(usersJson))
	})

	// test WHERE dua kondisi menggunakan 1 kali pemanggilan method WHERE -> langsung ditulis AND manual
	t.Run("test query WHERE with Find, single WHERE", func(t *testing.T) {
		users := []entity.User{}
		err := db.Where("first_name like ? AND password = ?", "user%", "rahasia").Find(&users).Error
		assert.Nil(t, err)

		usersJson, _ := json.Marshal(&users)
		log.Println(string(usersJson))
	})

	// OR Operator
	// test WHERE dengan kondisi kedua OR menggunakan pemanggilan method db.Where() dilanjutkan pemanggilan method Or()
	// maka gabungan dari dua method Where() dan Or() adalah query -> WHERE ... OR ...
	t.Run("test query WHERE dengan kondisi OR menggunakan method db.Where().OR()", func(t *testing.T) {
		var users = []entity.User{}

		firstName := "user ke 10"
		middleName := "Kurniawan"
		err := db.Where("first_name = ?", firstName).Or("middle_name = ?", middleName).Find(&users).Error
		assert.Nil(t, err)
		assert.Equal(t, 2, len(users))

		usersJson, _ := json.Marshal(&users)
		log.Println(string(usersJson))
	})

	// test WHERE dengan kondisi OR menggunakan 1 kali pemanggilan method db.Where() saja
	// kondisi OR dituliskan langsung di dalam parameter method Where()
	t.Run("test query WHERE dengan kondisi OR menggunakan method db.Where() saja", func(t *testing.T) {
		var users []entity.User

		firstName := "user ke 10"
		middleName := "Kurniawan"
		err := db.Where("first_name = ? OR middle_name = ?", firstName, middleName).Find(&users).Error
		assert.Nil(t, err)
		assert.Equal(t, 2, len(users))

		for _, user := range users {
			userJson, _ := json.Marshal(&user)
			log.Println(string(userJson))
		}
	})

	// NOT operator
	// test query WHERE NOT menggunakan method db.Not() saja
	// bisa juga ditambahkan method WHERE apabila memerlukan filter WHERE ... AND .. NOT ...
	t.Run("test query WHERE NOT dengan method db.Not() saja", func(t *testing.T) {
		var users []entity.User
		err := db.Not("middle_name = ?", "").Find(&users).Error
		assert.Nil(t, err)
		assert.Equal(t, 1, len(users))

		usersJson, _ := json.Marshal(&users)
		log.Println(string(usersJson))
	})

	// test query WHERE NOT menggunakan kombinasi method db.Where() dan method Not()
	t.Run("test query WHERE NOT menggunakan method Were() dan Not()", func(t *testing.T) {
		var users []entity.User
		password := "rahasia"
		firstName := "user%"
		err := db.Where("password = ?", password).Not("first_name like ?", firstName).Find(&users).Error
		assert.Nil(t, err)
		assert.Equal(t, 1, len(users))

		usersJson, _ := json.Marshal(&users)
		log.Println(string(usersJson))
	})

	t.Run("test query WHERE NOT hanya menggunakan method db.Where() saja", func(t *testing.T) {
		var user entity.User
		result := db.Where("password = ? AND first_name NOT LIKE ?", "rahasia", "user%").Take(&user)
		assert.Nil(t, result.Error)

		userJson, _ := json.Marshal(&user)
		log.Println(string(userJson))
	})
}

// SELECT condition
// TestQueryGormSelect adalah function UT yang berisi query SELECT
// method db.Select() digunakan untuk memilih kolom apa saja yang ingin ditampilkan
func TestQueryGormSelect(t *testing.T) {
	configMock := mck.NewConfigMock()
	configMock.Mock.On("GetConfig").Return(&cfg)
	db := connection.ConnectToDB(configMock)

	t.Run("test query SELECT menggunakan method db.Select()", func(t *testing.T) {
		var user entity.User
		err := db.Select("id, first_name").Not("middle_name = ?", "").Take(&user).Error
		assert.Nil(t, err)
		assert.NotEqual(t, "", user.Name.FirstName)

		userJson, _ := json.Marshal(&user)
		log.Println(string(userJson))
	})
}

// TestQueryGormWhereStructConditon adalah function UT yang berisi query WHERE menggunakan struct
// jadi kondisi WHERE filter bisa beda-beda sesuai dengan struct yang diisi
// ini cocok digunakan untuk kondisi dinamis, sehingga kolom yang dicari bisa beda-beda sesuai structnya
func TestQueryGormWhereStructConditon(t *testing.T) {
	configMock := mck.NewConfigMock()
	configMock.Mock.On("GetConfig").Return(&cfg)
	db := connection.ConnectToDB(configMock)

	// test query WHERE menggunakan method db.Where()
	// memasukkan filter tipe data struct ke dalam parameter method db.Where() nya
	t.Run("test query WHERE menggunakan struct", func(t *testing.T) {
		// create filter menggunakan struct
		// nanti querynya akan menjadi -> SELECT * FROM users WHERE password = 'rahasia' AND first_name = 'user ke 5'
		filter := entity.User{
			Password: "rahasia",
			Name:     entity.Name{FirstName: "user ke 5"},
		}

		var user entity.User
		err := db.Where(filter).Take(&user).Error
		assert.Nil(t, err)
		assert.Equal(t, "5", user.ID)

		userJson, _ := json.Marshal(&user)
		log.Println(string(userJson))
	})
}

// TestQueryGormWhereMapCondition adalah function UT yang berisi query WHERE menggunakan map
// jadi kondisi WHERE filter bisa beda-beda sesuai dengan map yang ditulis
// ini cocok digunakan untuk kondisi dinamis, sehingga kolom yang dicari bisa beda-beda sesuai dengan mapnya
func TestQueryGormWhereMapCondition(t *testing.T) {
	configMock := mck.NewConfigMock()
	configMock.Mock.On("GetConfig").Return(&cfg)
	db := connection.ConnectToDB(configMock)

	// create filter WHERE
	// nanti querynya menjadi -> SELECT * FROM users WHERE password = 'rahasia' AND last_name = ''
	filter := map[string]any{
		"password":  "rahasia",
		"last_name": "",
	}

	var users []entity.User
	err := db.Where(filter).Find(&users).Error
	assert.Nil(t, err)

	// print response
	usersJson, _ := json.Marshal(&users)
	log.Println(string(usersJson))
}

// TestQueryGormLimitOffset adalah function UT yang berisi query LIMIT dan OFFSET
// method db.Order() digunakan untuk menambah query ORDER
// method db.Offset() digunakan untuk menambah query OFFSET
// method db.Limit() digunakan untuk menambah query LIMIT
func TestQueryGormOrderLimitOffset(t *testing.T) {
	configMock := mck.NewConfigMock()
	configMock.Mock.On("GetConfig").Return(&cfg)
	db := connection.ConnectToDB(configMock)

	var users []entity.User
	err := db.Where("first_name LIKE ?", "user%").Order("id asc").Offset(2).Limit(2).Find(&users).Error
	assert.Nil(t, err)
	assert.Equal(t, 2, len(users))

	usersJson, _ := json.Marshal(&users)
	log.Println(string(usersJson))
}
