package helpers

import (
	"github.com/codemicro/nota/internal/database"
	"github.com/codemicro/nota/internal/models"
	"strconv"
)

// Helpers specifically for repeated API functions

func CheckAndConvertId(numberString string, idName string, modelType interface{}) (int64, bool, int, models.GenericResponse) {

	// Converts string ID (from c.Params or c.Body) to int64, then checks for its existence
	// Args - numberString - string that contains ID
	//        idName       - Name of the field that is being checked (eg. session, file, etc) - used in response formation
	//        modelType    - Type of model that is being checked for
	// Return format - id, hasFailed, httpCode, formedResponseModel

	id, err := strconv.ParseInt(numberString, 10, 64)

	// Check the ID is a valid number
	if err != nil {
		return 0, true, 400, models.GenericResponse{
			Status:  "error",
			Message: idName + " ID is not a valid number",
		}
	}

	conn := database.Conn

	// check there's a session with the specified ID in the db
	var count int
	conn.Model(modelType).Where("id = ?", int(id)).Count(&count)

	if count == 0 {
		return 0, true, 404, models.GenericResponse{
			Status:  "error",
			Message: idName + " ID not found",
		}
	}

	return id, false, 0, models.GenericResponse{}
}
