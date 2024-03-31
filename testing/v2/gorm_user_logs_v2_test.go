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

// TestGormInsertUserLogsAutoIncrement function UT untuk insert dan membuktikan
// bahwa id benar-benar auto increment
func TestGormInsertUserLogsAutoIncrement(t *testing.T) {
	configMock := mck.NewConfigMock()
	configMock.Mock.On("GetConfig").Return(&cfg)
	db := connection.ConnectToDB(configMock)

	// make data to insert
	var userLogs []entity.UserLog
	for i := 0; i < 10; i++ {
		userLogs = append(userLogs, entity.UserLog{
			UserId: "1",
			Action: "test action",
		})
	}

	err := db.Create(&userLogs).Error
	assert.Nil(t, err)

	// cek get data after insert
	var response []entity.UserLog
	err = db.Find(&response).Error
	assert.Nil(t, err)

	responseJson, _ := json.Marshal(&response)
	log.Println(string(responseJson))
}
