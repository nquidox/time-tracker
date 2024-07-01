package user

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
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

func (u *User) Create() error {
	var err error
	u.UserId = uuid.New()

	err = DB.Create(u).Error
	if err != nil {
		return err
	}

	return nil
}

func (u *User) ReadOne() error {
	var err error
	err = DB.Where("user_id = ?", u.UserId).First(u).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *User) ReadMany(filters map[string]interface{}) ([]User, error) {
	var users []User

	query := DB.Model(&User{})

	for k, v := range filters {
		query = query.Where(fmt.Sprintf("%s = ?", k), v)
	}

	err := query.Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *User) Update() error {
	result := DB.Where("user_id = ?", u.UserId).Updates(u)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (u *User) Delete() error {
	result := DB.Where("user_id = ?", u.UserId).Delete(u)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return nil
}
