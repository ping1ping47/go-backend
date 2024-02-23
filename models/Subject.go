package models

import "gorm.io/gorm"

// Subject represents the structure of your "subjects" table in the database.
type Subject struct {
    gorm.Model // Gorm's default model that includes fields like ID, CreatedAt, UpdatedAt, and DeletedAt

    // Your custom fields
    Name string
	Description string
    // Add more fields as needed
}

// You can add more fields and methods based on your specific requirements
