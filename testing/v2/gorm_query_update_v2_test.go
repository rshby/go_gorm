package testing_v2

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"go_gorm/infrastructure/database/connection"
	"go_gorm/model/entity"
	mck "go_gorm/testing/mock/config"
	"log"
	"testing"
)

// TestGormQueryUpdate adalah function UT untuk UPDATE di database
// untuk mengupdate, di gorm menggunakan method db.Save() dan db.Updates()
func TestGormQueryUpdate(t *testing.T) {
	configMock := mck.NewConfigMock()
	configMock.Mock.On("GetConfig").Return(&cfg)
	db := connection.ConnectToDB(configMock)

	// update semua kolom
	// jika kita menggunakan method db.Save(), maka itu akan mengupdate semua kolom yang bisa diupdate
	t.Run("update semua kolom yang ada di tabel", func(t *testing.T) {
		var user entity.User
		err := db.Where("id = ?", "1").Take(&user).Error
		assert.Nil(t, err)

		user.Name.FirstName = "Reo"
		user.Name.LastName = "Sahobby"
		user.Password = "rahasia123"

		// update to database
		err = db.Save(&user).Error
		assert.Nil(t, err)
	})

	// update selected column
	// jika kita hanya ingin spesifik update beberapa column, dapat menggunakan method db.Updates(map)
	t.Run("update selected column yang ada di tabel", func(t *testing.T) {
		var user entity.User

		// nanti querynya -> UPDATE users SET password = 'rahasia123', middle_name = '' WHERE id = '1';
		updateCondition := map[string]any{
			"password":    "rahasia123",
			"middle_name": "",
		}

		err := db.Model(&entity.User{}).Where("id = ?", "1").Updates(&updateCondition).Error
		assert.Nil(t, err)

		// get data after update
		err = db.Where("id = ?", "1").Take(&user).Error
		assert.Nil(t, err)

		userJson, _ := json.Marshal(&user)
		log.Println(string(userJson))
	})

	// update menggunakan input parameter struct
	t.Run("updates selected column with db.Updates(struct)", func(t *testing.T) {
		// query menjadi -> UPDATE users SET password = 'diubah lagi' WHERE id = '3'
		err := db.Model(&entity.User{}).Where("id = ?", "3").Updates(entity.User{
			Password: "diubah lagi",
			Name: entity.Name{
				MiddleName: "", // default value tidak akan diubah
			},
			Information: "ini akan diignore",
		}).Error
		assert.Nil(t, err)

		// get data after update
		var user entity.User
		err = db.Take(&user, "id = ?", "3").Error
		assert.Nil(t, err)

		userJson, _ := json.Marshal(&user)
		log.Println(string(userJson))
	})

	// update satu kolom saja
	// jika ingin melakukan update data dan hanya satu kolom yang diupdate, maka menggunakan db.Update(kolom, value)
	t.Run("query update satu kolom saja", func(t *testing.T) {
		err := db.Model(&entity.User{}).Where("id = ?", "2").Update("first_name", "Eko").Error
		assert.Nil(t, err)

		// get data after update
		var user entity.User
		err = db.Take(&user, "id = ?", "2").Error
		assert.Nil(t, err)

		userJson, _ := json.Marshal(&user)
		log.Println(string(userJson))
	})
}
