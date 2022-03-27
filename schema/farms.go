// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/log"
	"apricate/rdb"
	"bytes"
	"encoding/json"
	"fmt"
)

// enum for farm bonuses
type FarmBonuses uint8
const (
	FarmBonus_PristineSoil FarmBonuses = 0
	FarmBonus_Portals FarmBonuses = 1
	FarmBonus_Forested FarmBonuses = 2
	FarmBonus_NaturalFertilizer FarmBonuses = 3
	FarmBonus_ChronomicField FarmBonuses = 4
)

// Defines a farm
type Farm struct {
	UUID string `json:"uuid" binding:"required"`
	LocationSymbol string `json:"location_symbol" binding:"required"`
	Bonuses []FarmBonuses `json:"bonuses" binding:"required"`
	Buildings map[BuildingTypes]uint8 `json:"buildings" binding:"required"`
	Plots map[string]Plot `json:"plots" binding:"required"`
}

func NewFarm(pdb rdb.Database, totalplotcount uint64, username string, locationSymbol string) *Farm {
	var bonuses []FarmBonuses
	var buildings map[BuildingTypes]uint8
	var plots map[string]Plot
	switch locationSymbol {
	case "TS-PR-HF":
		bonuses = make([]FarmBonuses, 0)
		buildings = map[BuildingTypes]uint8{Building_Home: 1, Building_Field: 1, Building_SummoningCircle: 1}
		plots = NewPlots(pdb, username, totalplotcount, "TS-PR-HF", []Size{Huge, Large, Large, Average, Average, Modest, Modest})
	default:
		log.Error.Printf("Hit NewFarm with unknown LocationSymbol username: %s, locationSymbol: %s", username, locationSymbol)
		bonuses = make([]FarmBonuses, 0)
		buildings = make(map[BuildingTypes]uint8)
		plots = make(map[string]Plot)
	}
	return &Farm{
		UUID: username + "|Farm-" + locationSymbol,
		LocationSymbol: locationSymbol,
		Bonuses: bonuses,
		Buildings: buildings,
		Plots: plots,
	}
}

// Check DB for existing farm with given uuid and return bool for if exists, and error if error encountered
func CheckForExistingFarm (uuid string, tdb rdb.Database) (bool, error) {
	// Get farm
	_, getError := tdb.GetJsonData(uuid, ".")
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// error
			return false, getError
		}
		// farm not found
		return false, nil
	}
	// Got successfully
	return true, nil
}

// Get farm from DB, bool is farm found
func GetFarmFromDB (uuid string, tdb rdb.Database) (Farm, bool, error) {
	// Get farm json
	someJson, getError := tdb.GetJsonData(uuid, ".")
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// farm not found
			return Farm{}, false, nil
		}
		// error
		return Farm{}, false, getError
	}
	// Got successfully, unmarshal
	someData := Farm{}
	unmarshalErr := json.Unmarshal(someJson, &someData)
	if unmarshalErr != nil {
		log.Error.Fatalf("Could not unmarshal farm json from DB: %v", unmarshalErr)
		return Farm{}, false, unmarshalErr
	}
	return someData, true, nil
}

// Get farm from DB, bool is farm found
func GetFarmsFromDB (uuids []string, tdb rdb.Database) ([]Farm, bool, error) {
	// Get farm json
	someJson, getError := tdb.MGetJsonData(".", uuids)
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// farm not found
			return []Farm{}, false, nil
		}
		// error
		return []Farm{}, false, getError
	}
	// Got successfully, unmarshal
	someData := make([]Farm, len(someJson))
	for i, tempjson := range someJson {
		data := Farm{}
		unmarshalErr := json.Unmarshal(tempjson, &data)
		if unmarshalErr != nil {
			log.Error.Fatalf("Could not unmarshal farm json from DB: %v", unmarshalErr)
			return []Farm{}, false, unmarshalErr
		}
		someData[i] = data
	}
	
	return someData, true, nil
}

// Get farmdata at path from DB, bool is farm found
func GetFarmDataAtPathFromDB (uuid string, path string, tdb rdb.Database) (interface{}, bool, error) {
	// Get farm json
	someJson, getError := tdb.GetJsonData(uuid, path)
	if getError != nil {
		if fmt.Sprint(getError) != "redis: nil" {
			// farm not found
			return nil, false, nil
		}
		// error
		return nil, false, getError
	}
	// Got successfully, unmarshal
	var someData interface{}
	unmarshalErr := json.Unmarshal(someJson, &someData)
	if unmarshalErr != nil {
		log.Error.Fatalf("Could not unmarshal farm json from DB: %v", unmarshalErr)
		return nil, false, unmarshalErr
	}
	return someData, true, nil
}

// Attempt to save farm, returns error or nil if successful
func SaveFarmToDB(tdb rdb.Database, farmData *Farm) error {
	log.Debug.Printf("Saving farm %s to DB", farmData.UUID)
	err := tdb.SetJsonData(farmData.UUID, ".", farmData)
	// creationSuccess := rdb.CreateFarm(tdb, farmname, uuid, 0)
	return err
}

// Attempt to save farm data at path, returns error or nil if successful
func SaveFarmDataAtPathToDB(tdb rdb.Database, uuid string, path string, newValue interface{}) error {
	log.Debug.Printf("Saving farm data at path %s to DB for uuid %s", path, uuid)
	err := tdb.SetJsonData(uuid, path, newValue)
	return err
}

func (s FarmBonuses) String() string {
	return farmBonusesToString[s]
}

var farmBonusesToString = map[FarmBonuses]string {
	FarmBonus_PristineSoil: "Pristine Soil | Doubles base yield of plants.",
	FarmBonus_Portals: "Portals | Halve the travel time for any Assistant travelling to OR from the farm.",
	FarmBonus_Forested: "Forested | Two unique Buildings, the Logging Camp and Lumber Mill, are available for construction. Also comes with a free Axe tool, allowing the on-site harvest and processing of logs, lumber, and planks.",
	FarmBonus_NaturalFertilizer: "Natural Fertilizer | Thanks to the unique geography of the location, plants can naturally grow to Gigantic size",
	FarmBonus_ChronomicField: "Chronomic Field | A perpetual Chronomic Field covers the fields, allowing the accelerated growth of trees. Grow an orchard in your own lifetime!",
}

var farmBonusesToID = map[string]FarmBonuses {
	"Pristine Soil | Doubles base yield of plants.": FarmBonus_PristineSoil,
	"Portals | Halve the travel time for any Assistant travelling to OR from the farm.": FarmBonus_Portals,
	"Forested | Two unique Buildings, the Logging Camp and Lumber Mill, are available for construction. Also comes with a free Axe tool, allowing the on-site harvest and processing of logs, lumber, and planks.": FarmBonus_Forested,
	"Natural Fertilizer | Thanks to the unique geography of the location, plants can naturally grow to Gigantic size": FarmBonus_NaturalFertilizer,
	"Chronomic Field | An ancient Fae cast a perpetual Chronomic Field over the farm, allowing the accelerated growth of trees. Grow an orchard or lumberyard in days, not years!": FarmBonus_ChronomicField,
}

// MarshalJSON marshals the enum as a quoted json string
func (s FarmBonuses) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(farmBonusesToString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *FarmBonuses) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	// Note that if the string cannot be found then it will be set to the zero value, 'Created' in this case.
	*s = farmBonusesToID[j]
	return nil
}