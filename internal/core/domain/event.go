package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventStatus string

const (
	EventStatusUpcoming EventStatus = "upcoming"
	EventStatusActive   EventStatus = "active"
	EventStatusPast     EventStatus = "past"
	EventStatusCanceled EventStatus = "canceled"
)

type Event struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	RestaurantID  uuid.UUID      `gorm:"type:uuid;not null" json:"restaurant_id"`
	Title         string         `gorm:"type:varchar(255);not null" json:"title"`
	Description   string         `gorm:"type:text" json:"description"`
	Date          time.Time      `gorm:"not null" json:"date"`
	StartTime     time.Time      `gorm:"not null" json:"start_time"`
	EndTime       time.Time      `gorm:"not null" json:"end_time"`
	Location      string         `gorm:"type:varchar(255)" json:"location"`
	MaxGuests     int            `gorm:"default:0" json:"max_guests"`
	CurrentGuests int            `gorm:"default:0" json:"current_guests"`
	MaxVolunteers int            `gorm:"default:0" json:"max_volunteers"`
	Status        EventStatus    `gorm:"type:varchar(20);not null" json:"status"`
	MealsServed   int            `gorm:"default:0" json:"meals_served"`
	CreatedAt     time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Restaurant    Restaurant     `gorm:"foreignKey:RestaurantID" json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (e *Event) BeforeCreate(tx *gorm.DB) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	return nil
}

type VolunteerApplication struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	EventID     uuid.UUID      `gorm:"type:uuid;not null" json:"event_id"`
	VolunteerID uuid.UUID      `gorm:"type:uuid;not null" json:"volunteer_id"`
	Role        string         `gorm:"type:varchar(100);not null" json:"role"`
	Status      string         `gorm:"type:varchar(20);not null;default:'pending'" json:"status"` // pending, approved, declined
	AppliedAt   time.Time      `gorm:"autoCreateTime" json:"applied_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Event       Event          `gorm:"foreignKey:EventID" json:"-"`
	Volunteer   Volunteer      `gorm:"foreignKey:VolunteerID" json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (va *VolunteerApplication) BeforeCreate(tx *gorm.DB) error {
	if va.ID == uuid.Nil {
		va.ID = uuid.New()
	}
	return nil
}

type EventVolunteer struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	EventID     uuid.UUID      `gorm:"type:uuid;not null" json:"event_id"`
	VolunteerID uuid.UUID      `gorm:"type:uuid;not null" json:"volunteer_id"`
	Role        string         `gorm:"type:varchar(100);not null" json:"role"`
	CheckedIn   bool           `gorm:"default:false" json:"checked_in"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Event       Event          `gorm:"foreignKey:EventID" json:"-"`
	Volunteer   Volunteer      `gorm:"foreignKey:VolunteerID" json:"-"`
}

// BeforeCreate will set a UUID rather than numeric ID
func (ev *EventVolunteer) BeforeCreate(tx *gorm.DB) error {
	if ev.ID == uuid.Nil {
		ev.ID = uuid.New()
	}
	return nil
}
