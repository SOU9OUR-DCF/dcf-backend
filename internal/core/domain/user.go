package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserType string

const (
	UserTypeRegular    UserType = "regular"
	UserTypeRestaurant UserType = "restaurant"
	UserTypeVolunteer  UserType = "volunteer"
)

type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Username  string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"username"`
	Email     string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"`
	Type      UserType       `gorm:"type:varchar(20);not null;default:'volunteer'" json:"user_type"`
	CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}

type Restaurant struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	UserID        uuid.UUID      `gorm:"type:uuid;uniqueIndex;not null" json:"user_id"`
	Name          string         `gorm:"type:varchar(255);not null" json:"name"`
	Address       string         `gorm:"type:varchar(255)" json:"address"`
	ContactNumber string         `gorm:"type:varchar(50)" json:"contact_number"`
	TotalEvents   int            `gorm:"default:0" json:"total_events"`
	MealsServed   int            `gorm:"default:0" json:"meals_served"`
	Rating        float64        `gorm:"default:0" json:"rating"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	User          User           `gorm:"foreignKey:UserID" json:"-"`
}

func (r *Restaurant) BeforeCreate(tx *gorm.DB) error {
	if r.ID == uuid.Nil {
		r.ID = uuid.New()
	}
	return nil
}

type Volunteer struct {
	ID               uuid.UUID `json:"id" gorm:"primaryKey;type:uuid"`
	UserID           uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	FullName         string    `json:"full_name" gorm:"not null"`
	PhoneNumber      string    `json:"phone_number" gorm:"not null"`
	Address          string    `json:"address" gorm:"not null"`
	TasksCompleted   int       `json:"tasks_completed" gorm:"default:0"`
	HoursVolunteered int       `json:"hours_volunteered" gorm:"default:0"`
	MealsServed      int       `json:"meals_served" gorm:"default:0"`
	ReputationPoints int       `json:"reputation_points" gorm:"default:0"`
	CreatedAt        time.Time `json:"created_at" gorm:"not null"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"not null"`
}

func (v *Volunteer) BeforeCreate(tx *gorm.DB) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	return nil
}
