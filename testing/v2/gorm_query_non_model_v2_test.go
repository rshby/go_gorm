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

// TestQueryGormNonModel adalah function untuk query namun datanya tidak ditampung ke object model
// data hasil query akan ditampung ke object DTO yang sesuai
// namun gorm harus tau model apa yang digunakan, dengan menggunakan db.Model(&entity.User{}).Take(&response).Error
// maka hasil query akan disimpan ke object response
func TestQueryGormNonModel(t *testing.T) {
	configMock := mck.NewConfigMock()
	configMock.Mock.On("GetConfig").Return(&cfg)
	db := connection.ConnectToDB(configMock)

	// buat struct untuk menampung data id dan first_name saja
	type UserResponse struct {
		Id        string `json:"id,omitempty" gorm:"column:id"`
		FirstName string `json:"first_name,omitempty" gorm:"column:first_name"`
	}

	// query ke model User
	var response UserResponse
	err := db.Model(&entity.User{}).Select("id, first_name").Where("id = ?", "2").Take(&response).Error
	assert.Nil(t, err)
	assert.Equal(t, "2", response.Id)

	responseJson, _ := json.Marshal(&response)
	log.Println(string(responseJson))
}
