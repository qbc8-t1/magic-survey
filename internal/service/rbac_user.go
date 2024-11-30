package service

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/QBC8-Team1/magic-survey/domain/model"
	repository "github.com/QBC8-Team1/magic-survey/persistance"
)

func (o *RbacService) MakeFakeUser() (model.User, error) {
	index := MakeRandomNumber(400)
	return o.repo.MakeUser(model.User{
		FirstName:    "ufn" + strconv.Itoa(index),
		LastName:     "u2ln" + strconv.Itoa(index),
		Email:        "u2@email" + strconv.Itoa(index),
		NationalCode: strconv.Itoa(MakeRandomNumber(9999999999)),
		Password:     "111111111" + strconv.Itoa(index),
	})
}

func (o *RbacService) MakeFakeQuestionnaire(userID uint) (model.Questionnaire, error) {
	return o.repo.MakeQuestionnaire(model.Questionnaire{
		OwnerID:                    uint(userID),
		Status:                     model.QuestionnaireStatusOpen,
		MaxMinutesToResponse:       MakeRandomNumber(100),
		MaxMinutesToChangeAnswer:   MakeRandomNumber(100),
		MaxMinutesToGivebackAnswer: MakeRandomNumber(100),
		RandomOrSequential:         model.QuestionnaireTypeSequential,
		CanBackToPreviousQuestion:  false,
		Title:                      "test q" + strconv.Itoa(MakeRandomNumber(200000000)),
		MaxAllowedSubmissionsCount: 0,
		AnswersVisibleFor:          model.QuestionnaireVisibilityEverybody,
		CreatedAt:                  time.Now(),
	})
}

func MakeRandomNumber(max int) int {
	return rand.Intn(max)
}

func (o *RbacService) GetUser(userID uint) (model.User, error) {
	return o.repo.GetUser(userID)
}

func (o *RbacService) GetUserRolesWithPermissions(userID uint) ([]repository.RoleWithPermissions, error) {
	return o.repo.GetUserRolesWithPermissions(userID)
}
