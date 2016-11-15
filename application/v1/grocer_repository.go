package v1

import "github.com/jinzhu/gorm"

type GrocerRepositoryInterface interface {
	GetAll(pageNumber string, pageLimit string, relations string) ([]*Grocer, int)
	GetByID(id int, relations string) *Grocer
}

// GrocerRepository contain all function to retrieve list of grocer in database
type GrocerRepository struct {
	DB *gorm.DB
}

// GetAll function used to retrieve all grocers in the database
func (gr *GrocerRepository) GetAll(pageNumber string, pageLimit string, relations string) ([]*Grocer, int) {
	grocers := []*Grocer{}

	offset := SetOffsetValue(pageNumber, pageLimit)

	DB := gr.DB.Model(&Grocer{})

	if relations != "" {
		DB = LoadRelations(DB, relations)
	}

	DB.Offset(offset).Limit(pageLimit).Find(&grocers)

	var totalGrocers *int

	gr.DB.Model(&Grocer{}).Count(&totalGrocers)

	return grocers, *totalGrocers
}

// GetByID function used to retrieve grocer by ID in the database
func (gr *GrocerRepository) GetByID(id int, relations string) *Grocer {
	grocer := &Grocer{}

	gr.DB.Model(&Grocer{}).Where(&Grocer{ID: id}).First(&grocer)

	return grocer
}
