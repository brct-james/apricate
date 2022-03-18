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
	// JSON_Unmarshal_Error ResponseCode = 8
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
	// New_Status_Not_Allowed ResponseCode = 19
	Bad_Request ResponseCode = 20
	// No_Available_Routes ResponseCode = 21
	// Target_Route_Unavailable ResponseCode = 22
	// UDB_Update_Failed ResponseCode = 23
	// Leaderboard_Not_Found ResponseCode = 24
	// No_Resource_Nodes_At_Location ResponseCode = 25
	// Target_Resource_Node_Unavailable ResponseCode = 26
	// No_Packable_Items ResponseCode = 27
	// Invalid_Manifest ResponseCode = 28
	// Manifest_Overflow ResponseCode = 29
	// No_Storable_Items ResponseCode = 30
	// Blank_Manifest_Disallowed ResponseCode = 31
	// Blank_Order_Disallowed ResponseCode = 32
	// Merchant_Inventory_Empty ResponseCode = 33
	// Invalid_Order_Type ResponseCode = 34
	// Insufficient_Resources_Held ResponseCode = 35
	// Clearinghouse_Spool_Error ResponseCode = 36
	// Could_Not_Decode_Order ResponseCode = 37
	// Golem_Locked_For_Editing ResponseCode = 38
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
		Message: "[CRITICAL_JSON_MARSHAL_ERROR] Server error in responses.JSON, could not marshal JSON_Marshal_Error response! PLEASE contact developer.",
		HttpResponse: http.StatusInternalServerError,
	},
	JSON_Marshal_Error: {
		Message: "[JSON_Marshal_Error] Responses module encountered an error while marshaling response JSON. Please contact developer.",
		HttpResponse: http.StatusInternalServerError,
	},
	Unimplemented: {
		Message: "[Unimplemented] Unimplemented Feature. You shouldn't be able to hit this on the live build... Please contact developer",
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
		Message: "[Generate_Token_Failure] Username passed initial validation but could not generate token, contact Admin.",
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
	// JSON_Unmarshal_Error: {
	// 	Message: "[JSON_Unmarshal_Error] Error while attempting to unmarshal JSON from DB",
	// 	HttpResponse: http.StatusInternalServerError,
	// },
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
	// New_Status_Not_Allowed: {
	// 	Message: "[New_Status_Not_Allowed] Specified status is not valid for the specified golem's archetype",
	// 	HttpResponse: http.StatusNotAcceptable,
	// },
	Bad_Request: {
		Message: "[Bad_Request] Invalid request payload, please validate the request body conforms to expectations",
		HttpResponse: http.StatusBadRequest,
	},
	// No_Available_Routes: {
	// 	Message: "[No_Available_Routes] No routes available for the specified golem, contact developer",
	// 	HttpResponse: http.StatusInternalServerError,
	// },
	// Target_Route_Unavailable: {
	// 	Message: "[Target_Route_Unavailable] The specified route is not available at the current location",
	// 	HttpResponse: http.StatusNotAcceptable,
	// },
	// UDB_Update_Failed: {
	// 	Message: "[UDB_Update_Failed] Could not complete request due to error while saving user data to udb",
	// 	HttpResponse: http.StatusInternalServerError,
	// },
	// Leaderboard_Not_Found: {
	// 	Message: "[Leaderboard_Not_Found] Requested leaderboard not found",
	// 	HttpResponse: http.StatusNotFound,
	// },
	// No_Resource_Nodes_At_Location: {
	// 	Message: "[No_Resource_Nodes_At_Location] No resources nodes found at the location of the specified golem",
	// 	HttpResponse: http.StatusNotFound,
	// },
	// Target_Resource_Node_Unavailable: {
	// 	Message: "[Target_Resource_Node_Unavailable] The specified resource node is not available at the current location",
	// 	HttpResponse: http.StatusNotFound,
	// },
	// No_Packable_Items: {
	// 	Message: "[No_Packable_Items] No packable items in location inventory at specified golem's locale",
	// 	HttpResponse: http.StatusNotFound,
	// },
	// Invalid_Manifest: {
	// 	Message: "[Invalid_Manifest] Invalid manifest, specified item not contained in sufficient quantity in specified inventory",
	// 	HttpResponse: http.StatusBadRequest,
	// },
	// Manifest_Overflow: {
	// 	Message: "[Manifest_Overflow] Manifest is valid but requests more items than the golem can handle",
	// 	HttpResponse: http.StatusNotAcceptable,
	// },
	// No_Storable_Items: {
	// 	Message: "[No_Storable_Items] No storable items in the specified golem's inventory",
	// 	HttpResponse: http.StatusNotFound,
	// },
	// Blank_Manifest_Disallowed: {
	// 	Message: "[Blank_Manifest_Disallowed] Manifest cannot be blank, please includes items to load",
	// 	HttpResponse: http.StatusBadRequest,
	// },
	// Blank_Order_Disallowed: {
	// 	Message: "[Blank_Order_Disallowed] Order cannot be blank, please includes order type, item symbol, quantity, target price, and force_execution",
	// 	HttpResponse: http.StatusBadRequest,
	// },
	// Merchant_Inventory_Empty: {
	// 	Message: "[Merchant_Inventory_Empty] Merchant must be holding the items you wish to sell",
	// 	HttpResponse: http.StatusNotAcceptable,
	// },
	// Invalid_Order_Type: {
	// 	Message: "[Invalid_Order_Type] Type of the specified order does not match any known type",
	// 	HttpResponse: http.StatusBadRequest,
	// },
	// Insufficient_Resources_Held: {
	// 	Message: "[Insufficient_Resources_Held] The specified action could not be completed due to insufficient resources in golem inventory",
	// 	HttpResponse: http.StatusBadRequest,
	// },
	// Clearinghouse_Spool_Error: {
	// 	Message: "[Clearinghouse_Spool_Error] Order incorrectly spooled by clearinhouse, could not execute when asked",
	// 	HttpResponse: http.StatusInternalServerError,
	// },
	// Could_Not_Decode_Order: {
	// 	Message: "[Could_Not_Decode_Order] Error occurred while decoding order, ensure formatting is correct",
	// 	HttpResponse: http.StatusBadRequest,
	// },
	// Golem_Locked_For_Editing: {
	// 	Message: "[Golem_Locked_For_Editing] Server is handling another request for this golem and has locked the data, please retry request later",
	// 	HttpResponse: http.StatusConflict,
	// },
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