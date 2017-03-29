package v11

import (
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
)

// BaseRepository contain generic functions that can be used in other repository.
type BaseRepository struct{}

// LoadRelations function used to load model relation using eager loading technique provideby GORM ORM.
// Example value of relations parameter: `items.categories,items.subcategories,grocerexclusives`
//
// Refer this URL: http://jinzhu.me/gorm/crud.html#preloading-eager-loading
func (bs *BaseRepository) LoadRelations(DB *gorm.DB, relations string) *gorm.DB {
	// Split relations parameter into array by comma.
	splitRelations := strings.Split(relations, ",")

	for _, relation := range splitRelations {
		// Example: shopping_lists become []string{shopping, lists}
		splitUnderscoreInRelations := strings.Split(relation, "_")

		parseRelation := make([]string, len(splitUnderscoreInRelations))

		for key, splitUnderscoreInRelation := range splitUnderscoreInRelations {
			splitUnderscoreInRelation = strings.Title(splitUnderscoreInRelation)

			parseRelation[key] = splitUnderscoreInRelation
		}

		relation := strings.Join(parseRelation, "")

		splitNestedRelations := strings.Split(relation, ".")

		if len(splitNestedRelations) > 0 {
			DB = bs.processNestedRelations(DB, splitNestedRelations)
		} else {
			DB = DB.Preload(strings.Title(relation))
		}
	}

	return DB
}

// processNestedRelations used to load nested model relationship using eager loading method provide by GORM ORM.
// For example: API want to load deals with item relation but item relation also have
// other relation (Item Category, Item Subcategory) like below and API want to load those relations.
//
// Available relationship for Item:
// - Item has one Item Category.
// - Item has one Item Subcategory.
//
// Available relationship for Deal:
// - Deal has one Item.
// - Deal has one Category.
// - Deal has one Grocer Exclusive.
//
// API will know it needs to load nested relation if API found symbol `.` in the relations paramater.
// For example of nested relations : `items.categories,items.subcategories,`
//
// The example above will tell API to load deals with item relation inluding categories and subcategories inside item model.
// Refer this URL: http://jinzhu.me/gorm/crud.html#preloading-eager-loading
func (bs *BaseRepository) processNestedRelations(DB *gorm.DB, relations []string) *gorm.DB {
	nestedRelations := make([]string, len(relations))

	for key, relation := range relations {
		nestedRelations[key] = strings.Title(relation)
	}

	return DB.Preload(strings.Join(nestedRelations, "."))
}

// SetOffsetValue function used to set offset number based on page number and page limit parameters.
// For example: API receives request that contain page number with value 1.
//
// To retrieve it from database, you cannot set the offset value based on page number because
// in database, offset value start with 0.
// This function will set the offset value to 0 if page number equal to 0 or 1.
func (bs *BaseRepository) SetOffsetValue(pageNumber string, pageLimit string) int {
	pageNumberInt, _ := strconv.Atoi(pageNumber)
	pageLimitInt, _ := strconv.Atoi(pageLimit)

	offset := (pageNumberInt * pageLimitInt) - pageLimitInt

	if pageNumberInt == 1 || pageNumberInt == 0 {
		offset = 0
	}

	return offset
}
