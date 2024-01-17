package testing

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go_gorm/model/dto"
	"go_gorm/model/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"strconv"
	"testing"
)

func SetupDb() *gorm.DB {
	dsn := "root:root@tcp(localhost:3306)/belajar_golang_gorm?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
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

// test insert to database
func TestInsertToDB(t *testing.T) {
	db := SetupDb()

	t.Run("insert user", func(t *testing.T) {
		user := entity.User{
			ID:       "1",
			Password: "123",
			Name: entity.Name{
				FirstName: "Reo",
			},
			Information: "first user",
		}

		newData := db.Create(&user)
		assert.Nil(t, newData.Error)
		assert.Equal(t, 1, int(newData.RowsAffected))
	})

	t.Run("batch insert", func(t *testing.T) {
		// menggunakan method db.Create([]object)
		// atau db.CreateInBatcheds(slices, size)

		users := []entity.User{}
		for i := 2; i <= 10; i++ {
			users = append(users, entity.User{
				ID:       strconv.Itoa(i),
				Password: "123",
				Name: entity.Name{
					FirstName: fmt.Sprintf("User %v", i),
				},
			})
		}

		// insert batch
		tx := db.Create(&users)
		assert.Nil(t, tx.Error)
		assert.Equal(t, 9, int(tx.RowsAffected))
	})
}

// test transaction
func TestGormTransaction(t *testing.T) {
	db := SetupDb()

	// test using transaction
	t.Run("test transction", func(t *testing.T) {
		err := db.Transaction(func(tx *gorm.DB) error {
			// insert
			if err := tx.Create(&entity.User{
				ID:       "11",
				Password: "123",
				Name: entity.Name{
					FirstName: "user 11",
				},
			}).Error; err != nil {
				return err
			}

			if err := tx.Create(&entity.User{
				ID:       "12",
				Password: "123",
				Name: entity.Name{
					FirstName: "user 12",
				},
			}).Error; err != nil {
				return err
			}

			return nil
		})

		assert.Nil(t, err)
	})
}

// test gorm manual transaction
func TestManualTransaction(t *testing.T) {
	db := SetupDb()

	// test manual transaction success
	t.Run("test manual transaction success", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		// insert data ke 13
		err := tx.Create(&entity.User{
			ID:       "13",
			Password: "123",
			Name: entity.Name{
				FirstName: "user 13",
			},
		}).Error
		assert.Nil(t, err)

		// insert data ke 14
		err = tx.Create(&entity.User{
			ID:       "14",
			Password: "123",
			Name: entity.Name{
				FirstName: "user 14",
			},
		}).Error
		assert.Nil(t, err)

		if err == nil {
			tx.Commit()
		}
	})

	// test manual transaction failed
	t.Run("test manual transaction failed", func(t *testing.T) {
		tx := db.Begin()
		defer tx.Rollback()

		// insert data ke 15
		err := tx.Create(&entity.User{
			ID:       "15",
			Password: "123",
			Name: entity.Name{
				FirstName: "user 15",
			},
		}).Error
		assert.Nil(t, err)

		// insert data ke 14
		err = tx.Create(&entity.User{
			ID:       "14",
			Password: "123",
			Name: entity.Name{
				FirstName: "user 14",
			},
		}).Error
		assert.NotNil(t, err)

		if err == nil {
			tx.Commit()
		}
	})
}

// test query
func TestQuerSingleObject(t *testing.T) {
	db := SetupDb()

	// test get first data
	t.Run("test first single object", func(t *testing.T) {
		var user entity.User
		err := db.First(&entity.User{}).Scan(&user).Error
		if err != nil {
			t.Fatalf(err.Error())
		}

		assert.Equal(t, "1", user.ID)
	})

	// test get last data
	t.Run("test get last single object", func(t *testing.T) {
		var user entity.User
		err := db.Last(&entity.User{}).Scan(&user).Error
		assert.Nil(t, err)
		assert.Equal(t, "9", user.ID)
	})

	// test get first data with condition
	t.Run("test first with condition", func(t *testing.T) {
		var user entity.User
		err := db.First(&user, "id=? AND first_name=?", "2", "User 2").Error
		assert.Nil(t, err)
		assert.Equal(t, "2", user.ID)
	})
}

// test query all object
func TestQueryAllObject(t *testing.T) {
	db := SetupDb()

	// test get all data
	t.Run("find all data", func(t *testing.T) {
		var users []entity.User
		err := db.Find(&users, "id IN ?", []string{"1", "2", "3"}).Error
		assert.Nil(t, err)
		assert.Equal(t, 3, len(users))
	})
}

// test query menggunakan where
func TestQueryWhere(t *testing.T) {
	db := SetupDb()

	// test query menggunakan where -> single data
	t.Run("test query where", func(t *testing.T) {
		var user entity.User
		err := db.Where("id=?", "1").Where("first_name = ?", "Reo").Take(&user).Error
		assert.Nil(t, err)
		assert.Equal(t, "1", user.ID)
	})

	// test multi data menggunakan query LIKE
	t.Run("test query multi data", func(t *testing.T) {
		var users []entity.User
		err := db.Where("first_name LIKE ?", "user"+"%").Find(&users).Error
		assert.Nil(t, err)
		usersJson, _ := json.Marshal(&users)
		fmt.Println(string(usersJson))
	})

	// test menggunakan OR condition
	t.Run("test query OR condition", func(t *testing.T) {
		var users []entity.User
		err := db.Where("first_name like ?", "user"+"%").Or("password = ?", "rahasia").Find(&users).Error
		assert.Nil(t, err)
		assert.Equal(t, 13, len(users))
	})
}

// test menggunakan query not
func TestQueryNotOperator(t *testing.T) {
	db := SetupDb()

	t.Run("test query not operator", func(t *testing.T) {
		var users []entity.User
		err := db.Not("first_name LIKE ?", "user"+"%").Find(&users).Error
		assert.Nil(t, err)
		userJson, _ := json.Marshal(&users)
		fmt.Println(string(userJson))
	})
}

// test select specific column
func TestSelectColumn(t *testing.T) {
	db := SetupDb()

	// select specific column menggunakan .Select
	t.Run("test select specific column", func(t *testing.T) {
		var users []entity.User
		err := db.Select("id", "first_name").Find(&users).Error
		assert.Nil(t, err)

		usersJson, _ := json.Marshal(&users)
		fmt.Println(string(usersJson))
	})

	// select specific column menggunakan Struct
	t.Run("test using struct and map", func(t *testing.T) {
		var users []entity.User

		userCondition := entity.User{
			Name: entity.Name{
				FirstName: "User 5",
			},
		}

		// get from database
		err := db.Where(userCondition).Find(&users).Error
		assert.Nil(t, err)
		assert.Equal(t, 1, len(users))

		usersJson, _ := json.Marshal(&users)
		fmt.Println(string(usersJson))
	})

	// select specific column menggunakan map condition
	t.Run("test using map conditoni", func(t *testing.T) {
		db := SetupDb()

		condition := map[string]any{
			"middle_name": "",
		}

		var users []entity.User

		// get from database
		err := db.Where(condition).Find(&users).Error
		assert.Nil(t, err)
		assert.Equal(t, 14, len(users))

		usersJson, _ := json.Marshal(&users)
		fmt.Println(string(usersJson))
	})
}

// test Order, Limit, dan Offset
func TestOrderLimitOffset(t *testing.T) {
	/**
	untuk melakukan sorting, kita bisa menggunakan menthod Order()
	untuk melakukan paging, kita bisa menggunakan method Limit() dan Offset()
	**/

	db := SetupDb()

	// test limit order and offset
	t.Run("test order limit offset", func(t *testing.T) {
		var users []entity.User
		err := db.Order("id asc, first_name asc").Limit(5).Offset(5).Find(&users).Error
		assert.Nil(t, err)
		assert.Equal(t, 5, len(users))
		usersJson, _ := json.Marshal(&users)
		fmt.Println(string(usersJson))
	})
}

// test Model()
func TestModel(t *testing.T) {
	db := SetupDb()

	// test with dto
	t.Run("test with dto", func(t *testing.T) {
		var response []dto.UserResponse
		err := db.Model(&entity.User{}).Select("id", "first_name", "last_name").Find(&response).Error
		assert.Nil(t, err)
		assert.Equal(t, 14, len(response))

		responseJson, _ := json.Marshal(&response)
		fmt.Println(string(responseJson))
	})
}

// test update data
func TestUpdateData(t *testing.T) {
	db := SetupDb()

	// test update data with gorm menggunakan .Save()
	t.Run("test update data", func(t *testing.T) {
		// get data with id=1
		var user entity.User
		err := db.Where("id=?", "1").Take(&user).Error
		assert.Nil(t, err)

		// update data
		user.Name.FirstName = "Reo"
		user.Name.LastName = "Sahobby"
		err = db.Save(&user).Error
		assert.Nil(t, err)
	})

	// test update selected column menggunakan .Updates()
	t.Run("update selected multi column", func(t *testing.T) {
		/**
		secara default .Save() akan mengupdate semua kolom
		jika ingin hanya mengupdate kolom tertentu dapat menggunakan .Update(kolom, value)
		atau .Updates()
		**/
		err := db.Model(&entity.User{}).Where("id=?", 2).Updates(map[string]any{"first_name": "Budi", "last_name": "Haryanto"}).Error
		assert.Nil(t, err)
	})

	t.Run("updated selected multi columns with struct", func(t *testing.T) {
		// update
		err := db.Model(&entity.User{}).Where("id=?", "11").Updates(&entity.User{
			ID:   "11",
			Name: entity.Name{FirstName: "To", LastName: "Moro"},
		}).Error
		assert.Nil(t, err)

		// get data after update
		var user entity.User
		err = db.Take(&user, "id=?", "11").Error
		assert.Nil(t, err)

		// encode to json
		userJson, _ := json.Marshal(&user)
		fmt.Println(string(userJson))
	})

	// test update satu kolom menggunakan .Update(kolom, value_baru)
	t.Run("update single column", func(t *testing.T) {
		err := db.Model(&entity.User{}).Where("id=?", "10").Update("password", "pass_user10").Error
		assert.Nil(t, err)

		// get data
		var user entity.User
		err = db.Where("id=?", "10").Take(&user).Error
		assert.Nil(t, err)

		userJson, _ := json.Marshal(&user)
		fmt.Println(string(userJson))

	})
}
