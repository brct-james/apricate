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
	// No_WDB_Context ResponseCode = 9
	// No_UDB_Context ResponseCode = 10
	No_AuthPair_Context ResponseCode = 11
	User_Not_Found ResponseCode = 12
	// Not_Enough_Mana ResponseCode = 13
	// No_Such_Ritual ResponseCode = 14
	// No_Golem_Found ResponseCode = 15
	// Ritual_Not_Known ResponseCode = 16
	// No_Such_Status ResponseCode = 17
	// Golem_In_Blocking_Status ResponseCode = 18
	// New_Status_Not_Allowed ResponseCode = 19
	// Bad_Request ResponseCode = 20
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
	// No_WDB_Context: {
	// 	Message: "[No_WDB_Context] Could not get WDB context from middleware",
	// 	HttpResponse: http.StatusInternalServerError,
	// },
	// No_UDB_Context: {
	// 	Message: "[No_UDB_Context] Could not get UDB context from middleware",
	// 	HttpResponse: http.StatusInternalServerError,
	// },
	No_AuthPair_Context: {
		Message: "[No_AuthPair_Context] Failed to get AuthPair context from middleware",
		HttpResponse: http.StatusInternalServerError,
	},
	User_Not_Found: {
		Message: "[User_Not_Found] User not found!",
		HttpResponse: http.StatusNotFound,
	},
	// Not_Enough_Mana: {
	// 	Message: "[Not_Enough_Mana] Could not complete requested action due to insufficient mana",
	// 	HttpResponse: http.StatusNotAcceptable,
	// },
	// No_Such_Ritual: {
	// 	Message: "[No_Such_Ritual] The specified ritual is not recognized",
	// 	HttpResponse: http.StatusNotFound,
	// },
	// No_Golem_Found: {
	// 	Message: "[No_Golem_Found] Golem with the specified symbol could not be found in user data",
	// 	HttpResponse: http.StatusNotFound,
	// },
	// Ritual_Not_Known: {
	// 	Message: "[Ritual_Not_Known] User does not know the specified ritual, so it cannot be executed",
	// 	HttpResponse: http.StatusForbidden,
	// },
	// No_Such_Status: {
	// 	Message: "[No_Such_Status] Specified golem status does not exist",
	// 	HttpResponse: http.StatusBadRequest,
	// },
	// Golem_In_Blocking_Status: {
	// 	Message: "[Golem_In_Blocking_Status] Golem's current status does not allow changes to be made",
	// 	HttpResponse: http.StatusConflict,
	// },
	// New_Status_Not_Allowed: {
	// 	Message: "[New_Status_Not_Allowed] Specified status is not valid for the specified golem's archetype",
	// 	HttpResponse: http.StatusNotAcceptable,
	// },
	// Bad_Request: {
	// 	Message: "[Bad_Request] Invalid request payload, please validate the request body conforms to expectations",
	// 	HttpResponse: http.StatusBadRequest,
	// },
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