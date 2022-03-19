// Package responses provides helper functions, structs, and enums for formatting http responses
package responses

import (
	"encoding/json"
	"fmt"
	"net/http"

	"apricate/log"
)

// Prettifies input into json string for output
func JSON(input interface{}) (string, error) {
	res, err := json.MarshalIndent(input, "", "  ")
	if err != nil {
		return "", err
	}
	return string(res), nil
}

// enum for api response codes
type ResponseCode int
const (
	CRITICAL_JSON_MARSHAL_ERROR ResponseCode = -3
	JSON_Marshal_Error ResponseCode = -2
	Unimplemented ResponseCode = -1
	Generic_Failure ResponseCode = 0
	Generic_Success ResponseCode = 1
	Auth_Failure ResponseCode = 2
	Username_Validation_Failure ResponseCode = 3
	DB_Save_Failure ResponseCode = 4
	Generate_Token_Failure ResponseCode = 5
	DB_Get_Failure ResponseCode = 6
	UDB_Get_Failure ResponseCode = 7
	Internal_Server_Error ResponseCode = 8
	No_Assitant_At_Location ResponseCode = 9
	Specified_Plant_Not_Found ResponseCode = 10
	No_AuthPair_Context ResponseCode = 11
	User_Not_Found ResponseCode = 12
	Item_Does_Not_Exist ResponseCode = 13
	Not_Enough_Items_In_Warehouse ResponseCode = 14
	Plot_Already_Planted ResponseCode = 15
	Plot_Too_Small ResponseCode = 16
	Item_Is_Not_Seed ResponseCode = 17
	Plot_Already_Empty ResponseCode = 18
	Plot_Not_Planted ResponseCode = 19
	Bad_Request ResponseCode = 20
	Invalid_Plot_Action ResponseCode = 21
	Consumable_Not_In_Options ResponseCode = 22
	Missing_Consumable_Selection ResponseCode = 23
	Plants_Still_Growing ResponseCode = 24
	Tool_Not_Found ResponseCode = 25
)

// Defines Response structure for output
type Response struct {
	Code ResponseCode `json:"code" binding:"required"`
	Message string `json:"message" binding:"required"`
	Data interface{} `json:"data"`
}

type ResponseConfig struct {
	Message string `json:"message" binding:"required"`
	HttpResponse int `json:"http_response" binding:"required"`
}
var ResponseMap = map[ResponseCode]ResponseConfig{
	CRITICAL_JSON_MARSHAL_ERROR: {
		Message: "[CRITICAL_JSON_MARSHAL_ERROR] Server error in responses.JSON, could not marshal JSON_Marshal_Error response! PLEASE contact Developer.",
		HttpResponse: http.StatusInternalServerError,
	},
	JSON_Marshal_Error: {
		Message: "[JSON_Marshal_Error] Responses module encountered an error while marshaling response JSON. Contact Developer.",
		HttpResponse: http.StatusInternalServerError,
	},
	Unimplemented: {
		Message: "[Unimplemented] Unimplemented Feature. You shouldn't be able to hit this on the live build... Contact Developer",
		HttpResponse: http.StatusNotImplemented,
	},
	Generic_Failure: {
		Message: "[Generic_Failure] Contact developer",
		HttpResponse: http.StatusBadRequest,
	},
	Generic_Success: {
		Message: "[Generic_Success] Request Successful",
		HttpResponse: http.StatusOK,
	},
	Auth_Failure: {
		Message: "[Auth_Failure] Token was invalid or missing from request. Did you confirm sending the token as an authorization header?",
		HttpResponse: http.StatusUnauthorized,
	},
	Username_Validation_Failure: {
		Message: "[Username_Validation_Failure] Please ensure username conforms to requirements and account does not already exist!",
		HttpResponse: http.StatusBadRequest,
	},
	DB_Save_Failure: {
		Message: "[DB_Save_Failure] Failed to save to DB",
		HttpResponse: http.StatusInternalServerError,
	},
	Generate_Token_Failure: {
		Message: "[Generate_Token_Failure] Username passed initial validation but could not generate token, contact Developer.",
		HttpResponse: http.StatusInternalServerError,
	},
	DB_Get_Failure: {
		Message: "[DB_Get_Failure] Could not get requested information from DB",
		HttpResponse: http.StatusInternalServerError,
	},
	UDB_Get_Failure: {
		Message: "[UDB_Get_Failure] Could not get from user DB",
		HttpResponse: http.StatusInternalServerError,
	},
	Internal_Server_Error: {
		Message: "[Internal_Server_Error] Server encountered an error that is not the user's fault, contact Developer",
		HttpResponse: http.StatusInternalServerError,
	},
	No_Assitant_At_Location: {
		Message: "[No_Assitant_At_Location] Cannot see or interact with location without an assistant there. Verify that location exists, that name is spelled correctly (spaces may be replaced with underscores), and that an assistant is at the location specified",
		HttpResponse: http.StatusForbidden,
	},
	Specified_Plant_Not_Found: {
		Message: "[Specified_Plant_Not_Found] Could not get specified plant from dictionary",
		HttpResponse: http.StatusNotFound,
	},
	No_AuthPair_Context: {
		Message: "[No_AuthPair_Context] Failed to get AuthPair context from middleware",
		HttpResponse: http.StatusInternalServerError,
	},
	User_Not_Found: {
		Message: "[User_Not_Found] User not found!",
		HttpResponse: http.StatusNotFound,
	},
	Item_Does_Not_Exist: {
		Message: "[Item_Does_Not_Exist] Specified item does not exist in the master goods dictionary.",
		HttpResponse: http.StatusNotAcceptable,
	},
	Not_Enough_Items_In_Warehouse: {
		Message: "[Not_Enough_Items_In_Warehouse] The local warehouse does not have enough items for specified action.",
		HttpResponse: http.StatusNotFound,
	},
	Plot_Already_Planted: {
		Message: "[Plot_Already_Planted] Harvest or clear plot before attempting to plant it.",
		HttpResponse: http.StatusForbidden,
	},
	Plot_Too_Small: {
		Message: "[Plot_Too_Small] Plot too small for specified plant size & quantity. To pass validation, Quantity * SeedSize must be less-than or equal to PlotSize.",
		HttpResponse: http.StatusBadRequest,
	},
	Item_Is_Not_Seed: {
		Message: "[Item_Is_Not_Seed] Specified SeedName does not map to known Seed",
		HttpResponse: http.StatusBadRequest,
	},
	Plot_Already_Empty: {
		Message: "[Plot_Already_Empty] Plot is already empty, cannot be cleared",
		HttpResponse: http.StatusConflict,
	},
	Plot_Not_Planted: {
		Message: "[Plot_Not_Planted] Cannot interact with plot when not planted",
		HttpResponse: http.StatusConflict,
	},
	Bad_Request: {
		Message: "[Bad_Request] Invalid request payload, please validate the request body conforms to expectations",
		HttpResponse: http.StatusBadRequest,
	},
	Invalid_Plot_Action: {
		Message: "[Invalid_Plot_Action] Action specified in request body is either missing or fails to validate as either the current growth stage's Action or SkipAction (if applicable)",
		HttpResponse: http.StatusBadRequest,
	},
	Consumable_Not_In_Options: {
		Message: "[Consumable_Not_In_Options] The specified consumable does not match a valid option from the current Growth Stage",
		HttpResponse: http.StatusBadRequest,
	},
	Missing_Consumable_Selection: {
		Message: "[Missing_Consumable_Selection] The request body did not include a consumable selection, and specified action was not the SkipAction, consumables not optional",
		HttpResponse: http.StatusBadRequest,
	},
	Plants_Still_Growing: {
		Message: "[Plants_Still_Growing] The specified plot's plants are still growing, cannot interact yet",
		HttpResponse: http.StatusConflict,
	},
	Tool_Not_Found: {
		Message: "[Tool_Not_Found] The local warehouse does not have the requisite tool for the specified action",
		HttpResponse: http.StatusNotFound,
	},
}

// Returns the prettified json string of a properly structure api response given the inputs
func FormatResponse(code ResponseCode, data interface{}, messageDetail string) (string, int, error) {
	var message string
	var httpResponse int
	// Based on code choose base message text
	responseConfig, ok := ResponseMap[code]
	if !ok {
		message = "[Unexpected_Error] ResponseCode not in valid enum range! Contact developer"
		httpResponse = http.StatusInternalServerError
	} else {
		message = responseConfig.Message
		httpResponse = responseConfig.HttpResponse
	}
	
	// Define response
	var res Response = Response {
		Code: code,
		Message: message,
		Data: data,
	}

	// If messageDetail provided, append it
	if messageDetail != "" {
		res.Message = message + " | " + messageDetail
	}

	responseText, jsonErr := JSON(res)
	if jsonErr != nil {
		return "", httpResponse, jsonErr
	}
	return responseText, httpResponse, nil
}

func SendRes(w http.ResponseWriter, code ResponseCode, data interface{}, messageDetail string) {
	responseObject, httpResponse, jsonErr := FormatResponse(code, data, messageDetail)
	if jsonErr != nil {
		jsonErrMsg := fmt.Sprintf("Could not MarshallIndent json for data %v", data)
		errResponseObject, errHttpResponse, criticalJsonError := FormatResponse(JSON_Marshal_Error, nil, jsonErrMsg)
		if criticalJsonError != nil {
			log.Error.Printf("Could not format MarshallIndent response, error: %v", criticalJsonError)
			w.WriteHeader(http.StatusInternalServerError)
			msg := fmt.Sprintf("{\"code\":-3, \"message\": \"CRITICAL SERVER ERROR in responses.JSON, could not marshal JSON_Marshal_Error response! PLEASE contact developer. Error: %v\", \"data\":{}", criticalJsonError)
			w.Write([]byte(msg))
		}
		w.WriteHeader(errHttpResponse)
		w.Write([]byte(errResponseObject))
		// fmt.Fprint(w, errResponseObject)
	}

	w.WriteHeader(httpResponse)
	w.Write([]byte(responseObject))
	// fmt.Fprint(w, responseObject)
}