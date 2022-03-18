// Package schema defines database and JSON schema as structs, as well as functions for creating and using these structs
package schema

import (
	"apricate/filemngr"
	"apricate/log"

	"gopkg.in/yaml.v3"
)

// Define a good
type Good struct {
	Name string `yaml:"Name" json:"name" binding:"required"`
	Quantity uint64 `yaml:"Quantity" json:"quantity" binding:"required"`
}

// Define a raw good entry, that is processed by a generator to populate a good list
type RawGoodEntry struct {
	Name string `yaml:"Name" json:"name" binding:"required"`
	Seedy bool `yaml:"Seedy" json:"seedy,omitempty"`
	Enchantable bool `yaml:"Enchantable" json:"enchantable,omitempty"`
	Prefixes []string `yaml:"Prefixes" json:"prefixes,omitempty"`
	Suffixes []string `yaml:"Suffixes" json:"suffixes,omitempty"`
}

// Load good list by unmarhsalling given yaml file
func GoodListGenerator(path_to_goods_yaml string) map[string]interface{} {
	goodsBytes := filemngr.ReadFileToBytes(path_to_goods_yaml)
	var rawGoods []RawGoodEntry
	err := yaml.Unmarshal(goodsBytes, &rawGoods)
	if err != nil {
		log.Error.Fatalf("%v", err)
		// log.Error.Fatalf("%v", err.(*json.SyntaxError))
		// log.Error.Fatalf("%v", err.(*yaml.TypeError))
	}
	goodList := make(map[string]interface{}, 0)
	for _, good := range rawGoods {
		goodList[good.Name] = nil
		if good.Seedy {
			goodList[good.Name + " Seeds"] = nil
		}
		if good.Enchantable {
			goodList["Enchanted " + good.Name] = nil
		}
		for _, prefix := range good.Prefixes {
			goodList[prefix + " " + good.Name] = nil
			if good.Enchantable {
				goodList["Enchanted " + prefix + " " + good.Name] = nil
			}
		}
		for _, suffix := range good.Suffixes {
			goodList[good.Name + " " + suffix] = nil
			if good.Enchantable {
				goodList["Enchanted " + good.Name  + " " + suffix] = nil
			}
		}
	}
	return goodList
}