package v1_1

import (
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
)

type BaseRepository struct{}

func (bs *BaseRepository) LoadRelations(DB *gorm.DB, relations string) *gorm.DB {
	// Split on comma.
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

func (bs *BaseRepository) processNestedRelations(DB *gorm.DB, relations []string) *gorm.DB {
	nestedRelations := make([]string, len(relations))

	for key, relation := range relations {
		nestedRelations[key] = strings.Title(relation)
	}

	return DB.Preload(strings.Join(nestedRelations, "."))
}

func (bs *BaseRepository) SetOffsetValue(pageNumber string, pageLimit string) int {
	pageNumberInt, _ := strconv.Atoi(pageNumber)
	pageLimitInt, _ := strconv.Atoi(pageLimit)

	offset := (pageNumberInt * pageLimitInt) - pageLimitInt

	if pageNumberInt == 1 || pageNumberInt == 0 {
		offset = 0
	}

	return offset
}
