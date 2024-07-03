package user

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type FullUser struct {
	gorm.Model     `json:"-"`
	PassportSerie  int       `json:"passportSerie" `
	PassportNumber int       `json:"passportNumber"`
	Name           string    `json:"name"`
	Surname        string    `json:"surname"`
	Patronymic     string    `json:"patronymic"`
	Address        string    `json:"address"`
	UserId         uuid.UUID `json:"userId"`
}

type NewUser struct {
	PassportNumber string `json:"passportNumber" binding:"required"`
}

type UpdateUser struct {
	PassportSerie  int       `json:"passportSerie" `
	PassportNumber int       `json:"passportNumber"`
	Name           string    `json:"name"`
	Surname        string    `json:"surname"`
	Patronymic     string    `json:"patronymic"`
	Address        string    `json:"address"`
	UserId         uuid.UUID `json:"-"`
}

func (f *FullUser) TableName() string {
	return "users"
}

func (f *FullUser) Create() error {
	var err error

	err = DB.Create(f).Error
	if err != nil {
		return err
	}

	return nil
}

func (f *FullUser) ReadOne() error {
	var err error
	err = DB.Where("user_id = ?", f.UserId).First(f).Error
	if err != nil {
		return err
	}
	return nil
}

func (f *FullUser) ReadMany(filters map[string]interface{}, params map[string]int) ([]FullUser, error) {
	var users []FullUser

	query := DB.Model(&FullUser{})

	for k, v := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", k), v)
	}

	query = query.
		Offset((params["page"] - 1) * params["per_page"]).
		Limit(params["per_page"])

	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UpdateUser) Update() error {
	result := DB.Model(&FullUser{}).Where("user_id = ?", u.UserId).Updates(u)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (f *FullUser) Delete() error {
	result := DB.Where("user_id = ?", f.UserId).Delete(f)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}

func exists(serie, number int) uuid.UUID {
	var usr FullUser
	result := DB.Where("passport_serie = ? AND passport_number = ?", serie, number).First(&usr)
	if result.Error != nil {
		log.Error(usr.UserId)
		return usr.UserId
	} else {
		log.Error(uuid.Nil)
		return uuid.Nil
	}
}
