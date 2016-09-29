package v1

import "github.com/jinzhu/gorm"

type OccasionHandlerInterface interface {
	index()
}

type OccasionHandler struct {
	DB *gorm.DB
}

// func (oh *OccasionHandler) Index(c *gin.Context) {
// 	// Retrieve filter query string in request
// 	filterQueries := c.DefaultQuery("filter", "")
// 	fmt.Println(filterQueries)

// 	if filterQueries != "" {
// 		// Find if filterQueries contain word "or" OR "and"
// 		if strings.Contains(filterQueries, " and ") || strings.Contains(filterQueries, " or ") {
// 			// Split filterQuery by comma
// 			filterQueries := strings.Split(filterQueries, ",")
// 			fmt.Println(filterQueries)

// 			for _, filterQuery := range filterQueries {
// 				// Split filterQueries by space
// 				filters := strings.Split(filterQuery, " ")

// 				fmt.Println(len(filters))
// 				// If filters length below 3 return error message
// 				if len(filters) < 3 {

// 				}

// 				result := oh.DB.Where("name = ?", "jinzhu").First(&Occasion{})

// 			}
// 		}

// 	}

// }
