package auth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	"apricate/log"
	"apricate/rdb"
	"apricate/responses"
	"apricate/schema"

	"github.com/golang-jwt/jwt"
)

// GENERATE & VALIDATE TOKENS

// validate that username meets spec
func ValidateUsername (username string, slurs *[]string) string {
	// Defines acceptable chars
	isAlphaNumeric := regexp.MustCompile(`^[A-Za-z0-9\-\_]+$`).MatchString
	streamlinedName := strings.ToUpper(strings.Trim(strings.Trim(strings.Trim(strings.Trim(username, "-"), "_"), "-"), "_"))
	log.Debug.Printf("Streamlined name received: %s", streamlinedName)
	reservedSequences := []string{"NPC", "OWNER", "MOD-", "MOD_", "CONTRIBUTOR", "ADMIN-", "ADMIN_", "ADMINISTRATOR", "MODERATOR"}
	if username == "" {
		return "CANT_BE_BLANK"
	} else if len(username) <= 0 {
		return "TOO_SHORT"
	} else if len(username) > 32 {
		return "TOO_LONG"
	} else if !isAlphaNumeric(username) {
		return "INVALID_CHARS"
	} else {
		for _, word := range reservedSequences {
			if len(word) < len(streamlinedName) && streamlinedName[0:len(word)] == word {
				return fmt.Sprintf("RESERVED_SEQUENCE-NAME_CANNOT_BEGIN_WITH_-_%s_-_CONTACT_DEVELOPER_IF_YOU_BELIEVE_THIS_IS_A_MISTAKE", word)
			}
		}
		for _, slur := range *slurs {
			if strings.Contains(streamlinedName, slur) {
				return fmt.Sprintf("CONTAINS_SLUR-NAME_CANNOT_CONTAIN_-_%s_-_CONTACT_DEVELOPER_IF_YOU_BELIEVE_THIS_IS_A_MISTAKE", slur)
			}
		}
		return "OK"
	}
}

// HANDLE TOKEN VALIDATION FOR SECURE ROUTES

// Defines struct for passing around Token-Username pairs
type ValidationPair struct{
	Username string
	Token string
}

// enum for ValidationContext
type ValidationResponseKey int
const (
	ValidationContext ValidationResponseKey = iota
)

// Extract Token from request header
func ExtractToken(r *http.Request) (token string, ok bool) {
	bearerToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1], true
	}
	return "", false
}

// Extract token from header then parse and ensure confirms to signing method, if so return decoded token
func VerifyTokenFormatAndDecode(r *http.Request) (jwt.Token, error) {
	tokenString, ok := ExtractToken(r)
	if !ok {
		// Report failure to extract token
		log.Debug.Print("VerifyTokenFormatAndDecode: Error extracting token")
		return jwt.Token{}, fmt.Errorf("token extraction from header failed")
	}
	log.Debug.Printf("Token string: %s", tokenString)
	// Function for parsing token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure the token method conforms to SigningMethodHMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// return apricate_access_secret to parser for decoding
		access_secret := []byte(os.Getenv("APRICATE_ACCESS_SECRET"))
		if len(access_secret) == 0 {
			return nil, fmt.Errorf("internal error: could not get access secret from env")
		}
		return access_secret, nil
	})
	// Pass parse errors through to calling funcs
	if err != nil {
		return jwt.Token{}, err
	}
	// Return decoded token
	return *token, nil
}

// Verify token format and decode, then extract metadata (e.g. username) and return
func ExtractTokenMetadata(r *http.Request) (ValidationPair, error) {
	// Verify format and decode
	token, err := VerifyTokenFormatAndDecode(r)
	log.Debug.Printf("ExtractTokenMetadata:\nToken:\n%v\nError:\n%v\n", token, err)
	if err != nil {
		return ValidationPair{}, err
	}
	// ensure token.Claims is jwt.MapClaims
	claims, ok := token.Claims.(jwt.MapClaims)
	log.Debug.Printf("claims %v ok %v\n", claims, ok)
	log.Debug.Printf("token.Valid %v\n", token.Valid)
	if !ok || !token.Valid {
		// Fail state, token invalid and/or error
		return ValidationPair{}, fmt.Errorf("token invalid or token.Claims != jwt.MapClaims")
	}
	// Success state
	username := fmt.Sprintf("%s", claims["username"])
	log.Debug.Printf("username %v\n", username)
	// Return token and extracted username
	return ValidationPair{
		Token: token.Raw,
		Username: username,
	}, nil
}

// Verify that claimed authentication details are stored in database, if so return stored username, token, and ok=true
func AuthenticateWithDatabase(authD ValidationPair, userDB rdb.Database) (username string, token string, err error) {
	// Get user with claimed token
	dbuser, userFound, getUserErr := schema.GetUserFromDB(authD.Token, userDB)
	if getUserErr != nil {
		return "", "", getUserErr
	}
	if getUserErr != nil {
		// fail state
		getErrorMsg := fmt.Sprintf("in AuthenticateWithDatabase, could not get from DB for username: %s, token: %s, error: %v", authD.Username, authD.Token, getUserErr)
		log.Important.Println(getErrorMsg)
		return "", "", errors.New(getErrorMsg)
	}
	if !userFound {
		// fail state - user not found
		userNotFoundMsg := fmt.Sprintf("in AuthenticateWithDatabase, no user found in DB with username: %s, token: %s", authD.Username, authD.Token)
		log.Debug.Println(userNotFoundMsg)
		return "", "", errors.New("user not found")
	}
	log.Debug.Printf("AuthenticateWithDatabase, successfully got Username: %v, Token: %v\n", dbuser.Username, dbuser.Token)
	return dbuser.Username, dbuser.Token, nil
}

// Extract token metadata and check claimed token against database
func ValidateUserToken(r *http.Request, userDB rdb.Database) (username string, token string, err error) {
	// Extract metadata & validate
	tokenAuth, err := ExtractTokenMetadata(r)
	tokenAuthJsonString, tokenAuthJsonStringErr := responses.JSON(tokenAuth)
	if tokenAuthJsonStringErr != nil {
		log.Error.Printf("Error in ValidateUserToken, could not format tokenAuth as JSON. tokenAuth: %v, error: %v", tokenAuth, tokenAuthJsonStringErr)
	}
	log.Debug.Printf("ValidateUserToken:\nTokenAuth:\n%v\nError:\n%v\n", tokenAuthJsonString, err)
	if err != nil {
		return "", "", err
	}
	// Check against database for existing user
	dbusername, dbtoken, dbAuthErr := AuthenticateWithDatabase(tokenAuth, userDB)
	if dbAuthErr != nil {
		// Fail state, did not find user or could not get
		return "", "", dbAuthErr
	}
	// Success state, found user and matches
	return dbusername, dbtoken, nil
}

// Generates a middleware function for handling token validation on secure routes
func GenerateTokenValidationMiddlewareFunc(userDB rdb.Database) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Debug.Println(log.Yellow("-- GenerateTokenValidationMiddlewareFunc --"))
			// Validate bearer token
			username, token, validateTokenErr := ValidateUserToken(r, userDB)
			if validateTokenErr != nil {
				// Failed to validate, return failure message
				msg := fmt.Sprintf("%v", validateTokenErr)
				responses.SendRes(w, responses.Auth_Failure, nil, msg)
				return
			}
			// Create validation pair
			validationPair := ValidationPair{
				Username: username,
				Token: token,
			}
			validationPairJsonString, validationPairJsonStringErr := responses.JSON(validationPair)
			if validationPairJsonStringErr != nil {
				log.Error.Printf("Error in GenerateTokenValidationMiddlewareFunc, could not format validationPair as JSON. validationPair: %v, error: %v", validationPair, validationPairJsonStringErr)
			}
			log.Debug.Printf("validationPair:\n%v", validationPairJsonString)
			// Utilize context package to pass validation pair to secure routes from the middleware
			ctx := r.Context()
			ctx = context.WithValue(ctx, ValidationContext, validationPair)
			r = r.WithContext(ctx)
			// Continue serving route
			next.ServeHTTP(w,r)
			log.Debug.Println(log.Cyan("-- End GenerateTokenValidationMiddlewareFunc --"))
		})
	}
}