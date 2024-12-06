package rbac

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/QBC8-Team1/magic-survey/config"
	"github.com/QBC8-Team1/magic-survey/domain/model"
	"github.com/QBC8-Team1/magic-survey/internal/service"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
	"github.com/QBC8-Team1/magic-survey/pkg/db"
	"github.com/QBC8-Team1/magic-survey/pkg/logger"
	"gorm.io/gorm"
)

var configPath = flag.String("c", "../../config.yml", "Path to the configuration file")

func TestNotExistGiverUser(t *testing.T) {
	db, err := getDB()
	if err != nil {
		t.Fatalf(`failed to connect to DB`)
		return
	}

	rbacInstance := service.NewRbacService(repository.NewRbacRepository(db))
	err = rbacInstance.GivePermissions(0, 1, 2, []service.PermissionType{{Name: "something"}})
	if !errors.Is(err, model.ErrorNotFoundUser) {
		t.Fatalf(`not correct error for not found giver user: %s`, err)
	}
}

func TestNotExistQuestionnaire(t *testing.T) {
	db, err := getDB()
	if err != nil {
		t.Fatalf(`failed to connect to DB`)
		return
	}

	var user model.User
	db.Where("id = ?", 1).FirstOrCreate(&user, model.User{
		ID:           1,
		FirstName:    "u1fn",
		LastName:     "u1ln",
		Email:        "u1000000000000000@email.com",
		NationalCode: "0000000000",
		Password:     "111111111",
	})

	rbacInstance := service.NewRbacService(repository.NewRbacRepository(db))
	err = rbacInstance.GivePermissions(1, 1, 0, []service.PermissionType{{Name: "something"}})
	if !errors.Is(err, model.ErrorNotFoundQuestionnaire) {
		t.Fatalf(`not correct error for not found questionnaire: %s`, err)
	}
}

func TestNotOwnerOfQuestionnaire(t *testing.T) {
	db, err := getDB()
	if err != nil {
		t.Fatalf(`failed to connect to DB`)
		return
	}

	var user1 model.User
	db.Where("id = ?", 1).FirstOrCreate(&user1, model.User{
		ID:           1,
		FirstName:    "u1fn",
		LastName:     "u1ln",
		Email:        "u1000000000000000@email.com",
		NationalCode: "0000000000",
		Password:     "111111111",
	})

	var user2 model.User
	db.Where("id = ?", 1).FirstOrCreate(&user2, model.User{
		ID:           2,
		FirstName:    "u2fn",
		LastName:     "u2ln",
		Email:        "u2@email.com",
		NationalCode: "2222222222",
		Password:     "111111111",
	})

	questionnaire := model.Questionnaire{
		OwnerID:                    1,
		Status:                     model.QuestionnaireStatusOpen,
		MaxMinutesToResponse:       100,
		MaxMinutesToChangeAnswer:   100,
		MaxMinutesToGivebackAnswer: 100,
		RandomOrSequential:         model.QuestionnaireTypeSequential,
		CanBackToPreviousQuestion:  false,
		Title:                      "test q",
		MaxAllowedSubmissionsCount: 0,
		AnswersVisibleFor:          model.QuestionnaireVisibilityEverybody,
		CreatedAt:                  time.Now(),
	}

	result := db.Create(&questionnaire)
	if result.Error != nil {
		fmt.Println("Error creating record:", result.Error)
		return
	}

	rbacInstance := service.NewRbacService(repository.NewRbacRepository(db))
	err = rbacInstance.GivePermissions(2, 1, questionnaire.ID, []service.PermissionType{{Name: "something"}})
	if !errors.Is(err, model.ErrorNotHavePermission) {
		t.Fatalf(`not correct error for not owner of questionnaire: %s`, err)
	}
}

func TestIsReceiverUserExist(t *testing.T) {
	db, err := getDB()
	if err != nil {
		t.Fatalf(`failed to connect to DB`)
		return
	}

	var user1 model.User
	db.Where("id = ?", 1).FirstOrCreate(&user1, model.User{
		ID:           1,
		FirstName:    "u1fn",
		LastName:     "u1ln",
		Email:        "u1000000000000000@email.com",
		NationalCode: "0000000000",
		Password:     "111111111",
	})

	var user2 model.User
	db.Where("id = ?", 1).FirstOrCreate(&user2, model.User{
		ID:           2,
		FirstName:    "u2fn",
		LastName:     "u2ln",
		Email:        "u2@email.com",
		NationalCode: "2222222222",
		Password:     "111111111",
	})

	questionnaire := model.Questionnaire{
		OwnerID:                    user1.ID,
		Status:                     model.QuestionnaireStatusOpen,
		MaxMinutesToResponse:       100,
		MaxMinutesToChangeAnswer:   100,
		MaxMinutesToGivebackAnswer: 100,
		RandomOrSequential:         model.QuestionnaireTypeSequential,
		CanBackToPreviousQuestion:  false,
		Title:                      "test q",
		MaxAllowedSubmissionsCount: 0,
		AnswersVisibleFor:          model.QuestionnaireVisibilityEverybody,
		CreatedAt:                  time.Now(),
	}

	result := db.Create(&questionnaire)
	if result.Error != nil {
		fmt.Println("Error creating record:", result.Error)
		return
	}

	rbacInstance := service.NewRbacService(repository.NewRbacRepository(db))
	err = rbacInstance.GivePermissions(user1.ID, 0, questionnaire.ID, []service.PermissionType{{Name: "something"}})
	if !errors.Is(err, model.ErrorNotFoundUser) {
		t.Fatalf(`not correct error for not found receiver user: %s`, err)
	}
}

func TestNotFoundPermission(t *testing.T) {
	db, err := getDB()
	if err != nil {
		t.Fatalf(`failed to connect to DB`)
		return
	}

	var user1 model.User
	db.Where("id = ?", 1).FirstOrCreate(&user1, model.User{
		ID:           1,
		FirstName:    "u1fn",
		LastName:     "u1ln",
		Email:        "u1000000000000000@email.com",
		NationalCode: "0000000000",
		Password:     "111111111",
	})

	var user2 model.User
	db.Where("id = ?", 1).FirstOrCreate(&user2, model.User{
		ID:           2,
		FirstName:    "u2fn",
		LastName:     "u2ln",
		Email:        "u2@email.com",
		NationalCode: "2222222222",
		Password:     "111111111",
	})

	questionnaire := model.Questionnaire{
		OwnerID:                    user1.ID,
		Status:                     model.QuestionnaireStatusOpen,
		MaxMinutesToResponse:       100,
		MaxMinutesToChangeAnswer:   100,
		MaxMinutesToGivebackAnswer: 100,
		RandomOrSequential:         model.QuestionnaireTypeSequential,
		CanBackToPreviousQuestion:  false,
		Title:                      "test q",
		MaxAllowedSubmissionsCount: 0,
		AnswersVisibleFor:          model.QuestionnaireVisibilityEverybody,
		CreatedAt:                  time.Now(),
	}

	result := db.Create(&questionnaire)
	if result.Error != nil {
		fmt.Println("Error creating record:", result.Error)
		return
	}

	rbacInstance := service.NewRbacService(repository.NewRbacRepository(db))
	err = rbacInstance.GivePermissions(user1.ID, user2.ID, questionnaire.ID, []service.PermissionType{{Name: "something"}})
	if !errors.Is(err, model.ErrorNotFoundPermission) {
		t.Fatalf(`not correct error for not found receiver user: %s`, err)
	}
}

func TestMakeNewRole(t *testing.T) {
	db, err := getDB()
	if err != nil {
		t.Fatalf(`failed to connect to DB`)
		return
	}

	var user1 model.User
	db.Where("id = ?", 1).FirstOrCreate(&user1, model.User{
		ID:           1,
		FirstName:    "u1fn",
		LastName:     "u1ln",
		Email:        "u1000000000000000@email.com",
		NationalCode: "0000000000",
		Password:     "111111111",
	})

	var user2 model.User
	db.Where("id = ?", 1).FirstOrCreate(&user2, model.User{
		ID:           2,
		FirstName:    "u2fn",
		LastName:     "u2ln",
		Email:        "u2@email.com",
		NationalCode: "2222222222",
		Password:     "111111111",
	})

	questionnaire := model.Questionnaire{
		OwnerID:                    user1.ID,
		Status:                     model.QuestionnaireStatusOpen,
		MaxMinutesToResponse:       100,
		MaxMinutesToChangeAnswer:   100,
		MaxMinutesToGivebackAnswer: 100,
		RandomOrSequential:         model.QuestionnaireTypeSequential,
		CanBackToPreviousQuestion:  false,
		Title:                      "test q",
		MaxAllowedSubmissionsCount: 0,
		AnswersVisibleFor:          model.QuestionnaireVisibilityEverybody,
		CreatedAt:                  time.Now(),
	}

	result := db.Create(&questionnaire)
	if result.Error != nil {
		fmt.Println("Error creating record:", result.Error)
		return
	}

	rbacInstance := service.NewRbacService(repository.NewRbacRepository(db))
	rbacInstance.GivePermissions(user1.ID, user2.ID, questionnaire.ID, []service.PermissionType{{Name: model.PERMISSION_SEE_SELECTED_USERS_ANSWERS, SelectedUsersIds: []uint{1}}})
}

func getDB() (*gorm.DB, error) {
	flag.Parse()

	conf, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Panic(fmt.Errorf("load config error: %w", err))
		return nil, err
	}

	appLogger := logger.NewAppLogger(conf)
	appLogger.InitLogger(conf.Logger.Path)
	db, err := db.InitDB(conf, appLogger)
	if err != nil {
		appLogger.Panic("Counldnt init the db")
		return nil, err
	}

	return db, nil
}
