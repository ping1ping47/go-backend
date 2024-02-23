// models/Teacher.go
package models

import "gorm.io/gorm"

type Teacher struct {
	gorm.Model
	FirstName      string
	LastName       string
	Age            uint   // เพิ่มฟิลด์อายุ
	TeachingSubject string // เพิ่มฟิลด์วิชาที่สอน
}
