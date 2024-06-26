package testing

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"go_gorm/model/dto"
	"go_gorm/model/entity"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"log"
	"strconv"
	"testing"
	"time"
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

// test insert autoIncrement
func TestInsertAutoIncrement(t *testing.T) {
	db := SetupDb()

	// test insert batch dengan 10 data -> id auto_increment
	t.Run("insert 10 data with auto increment", func(t *testing.T) {
		var newData []entity.UserLog
		for i := 1; i <= 10; i++ {
			newData = append(newData, entity.UserLog{
				UserId: strconv.Itoa(i),
				Action: "created",
			})
		}

		err := db.Create(&newData).Error
		assert.Nil(t, err)
	})
}

// test isert created_at dan update_at : millisecond
func TestInsertTimeMilli(t *testing.T) {
	db := SetupDb()

	// insert user_logs timestamp using millisecond
	t.Run("test insert timestamp millisecond", func(t *testing.T) {
		var userLogs []entity.UserLog
		for i := 1; i <= 10; i++ {
			userLogs = append(userLogs, entity.UserLog{
				UserId: fmt.Sprintf("user %v", i),
				Action: "initial created",
			})
		}

		// insert to db
		tx := db.Create(&userLogs)
		assert.Nil(t, tx.Error)
		assert.Equal(t, 10, int(tx.RowsAffected))
	})
}

// test update or insert
func TestSaveOrUpdate(t *testing.T) {
	db := SetupDb()

	// test save or update
	t.Run("test insert or update", func(t *testing.T) {
		userLog := entity.UserLog{
			UserId: "1",
			Action: "test action",
		}

		tx := db.Save(&userLog) // akan insert -> karena data ID nya tidak ada
		assert.Nil(t, tx.Error)

		userLog.UserId = "2"
		err := db.Save(&userLog).Error // update
		assert.Nil(t, err)
	})

	// test data non auto_increment
	t.Run("test data non auto increment", func(t *testing.T) {
		user := entity.User{
			ID:       "99",
			Password: "",
			Name: entity.Name{
				FirstName: "user99",
			},
		}
		tx := db.Save(&user) // create
		assert.Nil(t, tx.Error)
		assert.Equal(t, 1, int(tx.RowsAffected))

		user.Name.FirstName = "user 99 updated"
		tx = db.Save(&user)
		assert.Nil(t, tx.Error)
		assert.Equal(t, 1, int(tx.RowsAffected))
	})
}

// test onconflict
func TestConflictCreate(t *testing.T) {
	db := SetupDb()

	// test on conflict
	t.Run("test on conflict", func(t *testing.T) {
		user := entity.User{
			ID:       "88",
			Password: "rahasia",
			Name: entity.Name{
				FirstName: "user 88 update",
			},
		}

		//
		result := db.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&user) // create
		assert.Nil(t, result.Error)
	})
}

// test delete data
func TestDeleteData(t *testing.T) {
	db := SetupDb()

	// delete tanpa select data dulu
	t.Run("test delete tanpa select dulu", func(t *testing.T) {
		err := db.Delete(&entity.User{}, "id=?", "99").Error
		assert.Nil(t, err)

		log.Println("success delete")
	})

	// delete dengan get data dulu
	t.Run("test delete dengan select data dulu", func(t *testing.T) {
		var user entity.User
		err := db.Take(&user, "id=?", "88").Error
		assert.Nil(t, err)

		userJson, _ := json.Marshal(&user)
		fmt.Println(string(userJson))

		// delete
		err = db.Delete(&user).Error
		assert.Nil(t, err)
	})

	// delete dengan condition WHERE
	t.Run("test delete dengan where", func(t *testing.T) {
		err := db.Where("password=? AND first_name=?", "rahasia123", "Reo").Delete(&entity.User{}).Error
		assert.Nil(t, err)

		// test get data yang sudah dihapus
		var user entity.User
		err = db.Where("password=? AND first_name=?", "rahasia123", "Reo").Take(&user).Error
		assert.NotNil(t, err)
	})
}

// test soft delete
func TestSoftDelete(t *testing.T) {
	db := SetupDb()

	// insert data todos
	t.Run("insert data todos", func(t *testing.T) {
		err := db.Model(&entity.Todo{}).Create(&entity.Todo{
			UserId:      "1",
			Title:       "data 1",
			Description: "ini adalah data 1",
		}).Error
		assert.Nil(t, err)
		log.Println("success insert data todos")
	})

	// test delete soft delete
	t.Run("test soft delete", func(t *testing.T) {
		err := db.Delete(&entity.Todo{ID: 1}).Error
		assert.Nil(t, err)

		log.Println("sucess delete soft delete")
	})

	// get data after soft delete
	t.Run("get data after soft delete", func(t *testing.T) {
		var todos []entity.Todo
		err := db.Find(&todos).Error
		assert.Nil(t, err)
		assert.Equal(t, 0, len(todos))
	})
}

// unscope -> digunakan untuk get data yang sudah disoft_delete atau ingin hard delete
func TestUnscope(t *testing.T) {
	db := SetupDb()

	// get data yang sudah soft_delete
	t.Run("get data yang sudah disoft_delete", func(t *testing.T) {
		var todo entity.Todo
		err := db.Unscoped().Take(&todo, "id=1").Error
		assert.Nil(t, err)
		assert.Equal(t, int64(1), todo.ID)
		todoJson, _ := json.Marshal(&todo)
		log.Println(string(todoJson))
	})

	// hard delete data yang sudah soft delete
	t.Run("hard delete data yang sudah soft delete", func(t *testing.T) {
		err := db.Unscoped().Delete(&entity.Todo{}, "id=?", 1).Error
		assert.Nil(t, err)

		// get data
		var todos []entity.Todo
		err = db.Model(&entity.Todo{}).Find(&todos).Error
		assert.Nil(t, err)
		assert.Equal(t, 0, len(todos))
	})

	/**
	Peringatan
	- ketika menggunakan soft delete, perhatikan penggunaan primary key atau unique index
	- ketika data sudah dihapus secara soft delete, sebenarnya data masih ada di tabel, oleh karena itu pastikan
	data primary key atau unique index tidak duplicate dengan data yang sudah dihapus secara soft_delete
	*/
}

// Model struct
func TestTodoLogs(t *testing.T) {
	db := SetupDb()

	// insert to table todo_logs
	t.Run("insert table todo_logs", func(t *testing.T) {
		err := db.Create(&entity.TodoLog{
			UserId:      "User 1",
			Title:       "Todo Log 1",
			Description: "created user 1",
		}).Error
		assert.Nil(t, err)

		// get data by id 1
		var todoLog entity.TodoLog
		err = db.Model(&entity.TodoLog{}).Take(&todoLog, "id=?", 1).Error
		assert.Nil(t, err)
		assert.Equal(t, uint(1), todoLog.Model.ID)

		// encode to json
		todoLogJson, _ := json.Marshal(&todoLog)
		log.Println(string(todoLogJson))
	})
}

// test locking database
func TestLockingDatabase(t *testing.T) {
	db := SetupDb()

	t.Run("test locking for update", func(t *testing.T) {
		err := db.Transaction(func(tx *gorm.DB) error {
			// get data and lock for update
			var todoLog entity.TodoLog
			err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Take(&todoLog, "id=1").Error
			if err != nil {
				return err
			}

			// update data
			todoJson, _ := json.Marshal(&todoLog)
			log.Println(string(todoJson))

			todoLog.Title = "Todo Log User 1"
			err = tx.Updates(&todoLog).Error
			return err
		})

		assert.Nil(t, err)
	})
}

// test insert table wallet
func TestInsertWallet(t *testing.T) {
	db := SetupDb()

	t.Run("insert wallet", func(t *testing.T) {
		err := db.Model(&entity.Wallet{}).Create(&entity.Wallet{
			Id:      "2",
			UserId:  "2",
			Balance: 1000000,
		}).Error
		assert.Nil(t, err)

		// get data after update
		var wallet entity.Wallet
		err = db.Take(&wallet, "user_id=?", "2").Error
		assert.Nil(t, err)

		walletJson, _ := json.Marshal(&wallet)
		log.Println(string(walletJson))
	})
}

// func test query join relation table
func TestQueryRelationTable(t *testing.T) {
	db := SetupDb()

	// get data from table users and wallets
	t.Run("get data menggunakan preload", func(t *testing.T) {
		var user entity.User
		err := db.Model(&entity.User{}).Preload("Wallet").Take(&user, "id=?", "2").Error
		assert.Nil(t, err)

		userJson, _ := json.Marshal(&user)
		log.Println(string(userJson))

		assert.Equal(t, "2", user.ID)
		//assert.Equal(t, "2", user.Wallet.Id)
	})

	// get data menggunakan joins
	t.Run("get data menggunakan join", func(t *testing.T) {
		var user entity.User
		err := db.Model(&entity.User{}).Joins("Wallet").Take(&user, "users.id=2").Error
		assert.Nil(t, err)

		userJson, _ := json.Marshal(&user)
		log.Println(string(userJson))
	})
}

// auto upsert relation
func TestAutoCreateUpdateRelation(t *testing.T) {
	db := SetupDb()

	// test auto insert 2 tables (yang mempunyai relasi)
	t.Run("auto create update relation", func(t *testing.T) {
		// create data that will be inserted to database
		user := entity.User{
			ID:       "20",
			Password: "rahasia",
			Name: entity.Name{
				FirstName:  "Muhammad",
				MiddleName: "Reo",
				LastName:   "Sahobby",
			},
			/*
				Wallet: &entity.Wallet{
					Id:      "3",
					UserId:  "20",
					Balance: 18000000,
				}
				,
			*/
		}

		// insert to database tabel users and wallets
		err := db.Create(&user).Error
		assert.Nil(t, err)
	})

	// test skip auto create/update -> tidak menggunakan auto create/auto update
	t.Run("skip auto create update relation", func(t *testing.T) {
		// create data that will be inserted/updated
		user := entity.User{
			ID:       "21",
			Password: "rahasia",
			Name: entity.Name{
				FirstName: "Reo",
				LastName:  "Sahobby",
			},
			/*
				Wallet: &entity.Wallet{

					Id:      "4",
					UserId:  "21",
					Balance: 10000000,
				}
			*/

		}

		// insert to database only table users
		err := db.Omit(clause.Associations).Create(&user).Error
		assert.Nil(t, err)

		log.Println("hanya insert data ke tabel users")
	})

	// insert data users sekaligus data addressesnya
	t.Run("insert data users and address directly", func(t *testing.T) {
		// create data
		user := entity.User{
			ID:       "22",
			Password: "P@ssw0rd",
			Name: entity.Name{
				FirstName:  "Cotton",
				MiddleName: "Buds",
				LastName:   "Adult",
			},
			/*
				Wallet: &entity.Wallet{
					Id:      "4",
					UserId:  "22",
					Balance: 98000000,
				},
				Addresses: []entity.Address{
					{
						UserId:  "22",
						Address: "Tegal Baru, Gumulan, Klaten Tengah",
					},
					{
						UserId:  "22",
						Address: "Ragunan, Pasar Minggu, Jakarta Selatan",
					},
				},

			*/
		}

		// insert all to database
		err := db.Save(&user).Error
		assert.Nil(t, err)

		log.Println("success insert all to database")
	})

	// get data users and address (relation 2 table -> one to many)
	t.Run("test get users and address", func(t *testing.T) {
		var user []entity.User
		err := db.Model(&entity.User{}).Preload("Addresses").Where("id=?", "22").Find(&user).Error
		assert.Nil(t, err)

		userJson, _ := json.Marshal(&user)
		fmt.Println(string(userJson))
	})

	// test get data yang memiliki relasi 3 tabel
	t.Run("get data with 3 relation tables", func(t *testing.T) {
		var users []entity.User
		err := db.Model(&entity.User{}).Preload("Addresses").Joins("Wallet").Find(&users).Error
		assert.Nil(t, err)

		usersJson, _ := json.Marshal(&users)
		log.Println(string(usersJson))
	})

	// test get one data yang memiliki 3 relasi
	t.Run("get one data with 3 relation tables", func(t *testing.T) {
		var user entity.User
		err := db.Model(&entity.User{}).Preload("Addresses").Joins("Wallet").Take(&user, "users.id=?", "22").Error
		assert.Nil(t, err)

		userJson, _ := json.Marshal(&user)
		fmt.Println(string(userJson))
	})
}

// test belong to
func TestBelongsTo(t *testing.T) {
	db := SetupDb()

	// test data from belongs to using joins
	t.Run("test get data from belongs to table using joins", func(t *testing.T) {
		var address entity.Address
		err := db.Model(&entity.Address{}).Joins("User").Take(&address, "addresses.user_id=?", "22").Error
		assert.Nil(t, err)

		addressJson, _ := json.Marshal(&address)
		fmt.Println(string(addressJson))
	})

	// test get data from belongs using preload
	t.Run("test get data from belongs to table using preload", func(t *testing.T) {
		var address entity.Address
		err := db.Model(&entity.Address{}).Preload("User").Take(&address, "user_id=?", "22").Error
		assert.Nil(t, err)

		addressJson, _ := json.Marshal(&address)
		fmt.Println(string(addressJson))
	})

	// test get data from wallet -> belongs to one to one
	t.Run("get data belongs to one to one", func(t *testing.T) {
		var wallet entity.Wallet
		err := db.Model(&entity.Wallet{}).Preload("User").Take(&wallet, "user_id=?", "22").Error
		assert.Nil(t, err)

		walletJson, _ := json.Marshal(&wallet)
		fmt.Println(string(walletJson))
	})

	t.Run("get data wallet from belongs to using joins", func(t *testing.T) {
		var wallet entity.Wallet
		err := db.Model(&entity.Wallet{}).Joins("User").Take(&wallet, "wallets.user_id=?", "22").Error
		assert.Nil(t, err)

		walletJson, _ := json.Marshal(&wallet)
		fmt.Println(string(walletJson))
	})
}

// test insert many to many
func TestManyTomany(t *testing.T) {
	db := SetupDb()

	t.Run("test insert manytomany", func(t *testing.T) {
		// insert to product
		err := db.Model(&entity.Product{}).Omit(clause.Associations).Where("id='P001'").Save(&entity.Product{
			ID:        "P001",
			Name:      "iPhone 15 Pro Max 512GB",
			Price:     30000000,
			CreatedAt: time.Now(),
		}).Error
		assert.Nil(t, err)

		err = db.Table("user_like_product").Create(map[string]any{
			"user_id":    "1",
			"product_id": "P001",
		}).Error
		assert.Nil(t, err)
	})

	// test get data from many to many
	t.Run("test get data many to many", func(t *testing.T) {
		var product entity.Product
		err := db.Model(&entity.Product{}).Preload("LikeByUsers").First(&product, "id=?", "P001").Error
		assert.Nil(t, err)

		productJson, _ := json.Marshal(&product)
		fmt.Println(string(productJson))
	})

	// test get data many2many from user
	t.Run("get data many to many from user", func(t *testing.T) {
		var user entity.User
		err := db.Model(&entity.User{}).Preload("LikeProducts").First(&user, "id=?", "1").Error
		assert.Nil(t, err)

		userJson, _ := json.Marshal(&user)
		fmt.Println(string(userJson))
	})
}
