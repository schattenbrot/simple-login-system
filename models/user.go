package models

import "time"

type User struct {
	ID                       string    `json:"id,omitempty" bson:"_id,omitempty"`
	Name                     string    `json:"name,omitempty" bson:"name,omitempty"`
	Email                    string    `json:"email" bson:"email,omitempty"`
	Password                 string    `json:"password,omitempty" bson:"password,omitempty"`
	ResetPasswordToken       string    `json:"resetPasswordToken,omitempty" bson:"resetPasswordToken,omitempty"`
	ResetPasswordTokenExpire time.Time `json:"resetPasswordTokenExpire,omitempty" bson:"resetPasswordTokenExpire,omitempty"`
	CreatedAt                time.Time `json:"createdAt" bson:"createdAt,omitempty"`
	UpdatedAt                time.Time `json:"updatedAt" bson:"updatedAt,omitempty"`
}
