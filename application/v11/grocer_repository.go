package v11

import "github.com/jinzhu/gorm"

// GrocerRepository contain all function to retrieve list of grocer in database
type GrocerRepository struct {
	BaseRepository
	DB *gorm.DB
}

// GetAll function used to retrieve all grocers in the database
func (gr *GrocerRepository) GetAll(pageNumber string, pageLimit string, relations string) ([]*Grocer, int) {
	grocers := []*Grocer{}

	offset := gr.SetOffsetValue(pageNumber, pageLimit)

	DB := gr.DB.Model(&Grocer{})

	if relations != "" {
		DB = gr.LoadRelations(DB, relations)
	}

	DB.Offset(offset).Limit(pageLimit).Find(&grocers)

	var totalGrocers *int

	gr.DB.Model(&Grocer{}).Count(&totalGrocers)

	return grocers, *totalGrocers
}

// GetAllGrocersThoseOnlyHaveDeal function used to retrieve all grocers those only have deals.
func (gr *GrocerRepository) GetAllGrocersThoseOnlyHaveDeal() []*Grocer {
	grocers := []*Grocer{}

	gr.DB.Model(&Grocer{}).Joins("INNER JOIN ads_grocer ON ads_grocer.grocer_id = grocer.id").Where(&Grocer{Status: "publish"}).Group("grocer.id").Find(&grocers)

	return grocers
}

// GetByID function used to retrieve grocer by ID in the database
func (gr *GrocerRepository) GetByID(id int, relations string) *Grocer {
	grocer := &Grocer{}

	gr.DB.Model(&Grocer{}).Where(&Grocer{ID: id}).First(&grocer)

	return grocer
}

// GetByGUID function used to retrieve grocer by GUID in the database
func (gr *GrocerRepository) GetByGUID(grocerGUID, relations string) *Grocer {
	grocer := &Grocer{}

	gr.DB.Model(&Grocer{}).Where(&Grocer{GUID: grocerGUID}).First(&grocer)

	return grocer
}

// GetByGUIDAndStatus function used to retrieve grocer by GUID and status in the database
func (gr *GrocerRepository) GetByGUIDAndStatus(grocerGUID, status, relations string) *Grocer {
	grocer := &Grocer{}

	gr.DB.Model(&Grocer{}).Where(&Grocer{GUID: grocerGUID, Status: "publish"}).First(&grocer)

	return grocer
}
