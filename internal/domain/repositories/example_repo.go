package repositories

import(
	"github.com/0xdiaz/gin-boilerplate/internal/domain/models"
	"github.com/0xdiaz/gin-boilerplate/internal/adapters/database"
	"github.com/0xdiaz/Datatables-Gin/datatables"

	"github.com/gin-gonic/gin"
)

func GetDataDatatables(c *gin.Context) (interface{}, error) {
	query := database.GetDB().Model(&models.Example{})
	var example []*models.Example
	result, err := datatables.OfReturn(
		c,
		query,
		&example,
		[]string{"id","data"},
		map[string]string{ // orderable
			"id":           "id",
			"data":          "data",
		},
		datatables.NewOptions().
			WithIndex("DT_RowIndex", false), // global index

	)

	if err != nil {
		return nil, err
	}

	return result,nil
}