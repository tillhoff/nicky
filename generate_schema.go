package main

import (
	"log"
	"reflect"

	"gopkg.in/yaml.v3"
)

func generateFoldedSchema(source string) string {

	/// the following block starts recursion over the yaml file and generates an example schema, for more details on this process see method generate_recursive

	mappedSource := make(map[string]interface{}, 0) // create empty map
	yaml.Unmarshal([]byte(source), &mappedSource)   // store yaml into map

	mappedSchema := make(map[string]interface{}, 0)               // create empty map
	generatedProperties := generateRecursiveMap("", mappedSource) // start recursion
	mappedSchema["properties"] = generatedProperties              // set the map generated by recursion as properties of root-element

	/// the recursion is now finished
	/// the next block adds some basic values to the root element

	mappedSchema["name"] = "schema_generated_by_nicky"           // set name of schema, which is used in the unfolding-process
	mappedSchema["required"] = getKeysOfMap(generatedProperties) // set all first level properties as required

	/// the whole map is now finished
	/// the next block converts the map to yaml and finally returns it
	generatedSchema, err := yaml.Marshal(&mappedSchema) // store map into yaml
	if err != nil {                                     // if error is thrown
		log.Fatal(err) // throw error
	}

	return string(generatedSchema) // return yaml
}

func generateRecursiveMap(parentkey string, m map[string]interface{}) map[string]interface{} {

	newMap := make(map[string]interface{}) // initialize empty map

	for key, value := range m { // search all keys within m
		switch value.(type) { // compare type of value
		case bool:
			newMap[key] = "boolean" // set variable type to boolean
		case int:
			newMap[key] = "integer" // set variable type to integer
		case string:
			stringMap := make(map[string]interface{})
			stringMap["pattern"] = value
			newMap[key] = stringMap
		case map[string]interface{}: // object
			if mappedValue, ok := value.(map[string]interface{}); ok { // if value is a map/object
				objectMap := make(map[string]interface{})
				objectMap["properties"] = generateRecursiveMap(key, mappedValue) // start recursion and set value to result
				newMap[key] = objectMap
			}
		case []interface{}: //list
			if interfacedValue, ok := value.([]interface{}); ok { // if value is a slice/list/array
				newMap[key] = generateRecursiveInterface(key, interfacedValue) // start recursion and add new map to local map
			}
		default:
			log.Fatal("The type '" + reflect.TypeOf(value).String() + "' is not implemented. This problem occured on a key with the name '" + key + "'.") //unkown type error
		}
	}

	return newMap
}

func generateRecursiveInterface(parentkey string, i []interface{}) map[string]interface{} {

	//newMap := make(map[string]interface{})      // initialize empty map
	newInterface := make([]interface{}, 0) // initialize empty map

	for _, value := range i { // search all keys within i
		switch value.(type) { // compare type of value
		case bool:
			newMap := make(map[string]interface{})      // initialize empty map
			newMap["type"] = "boolean"                  // set variable type to boolean
			newInterface = append(newInterface, newMap) // add new map to local interface
		case int:
			newMap := make(map[string]interface{})      // initialize empty map
			newMap["type"] = "integer"                  // set variable type to integer
			newInterface = append(newInterface, newMap) // add new map to local interface
		case string:
			stringMap := make(map[string]interface{})
			stringMap["pattern"] = value                   // set pattern to value to match exactly
			newInterface = append(newInterface, stringMap) // add new map to local interface
		case map[string]interface{}: // object
			if mappedValue, ok := value.(map[string]interface{}); ok { // if value is a map/object
				objectMap := make(map[string]interface{})
				objectMap["properties"] = generateRecursiveMap("", mappedValue) // start recursion and set value to result
				newInterface = append(newInterface, objectMap)
			}
		case []interface{}: //list
			if interfacedValue, ok := value.([]interface{}); ok { // if value is a slice/list/array
				newInterface = append(newInterface, generateRecursiveInterface("", interfacedValue)) // start recursion and add new map to local interface
			}
		default:
			log.Fatal("The type '" + reflect.TypeOf(value).String() + "' is not implemented. This problem occured on a interface with the parent named '" + parentkey + "'.") //unkown type error
		}
	}

	oneofMap := make(map[string]interface{}) // initialize empty map
	oneofMap["oneOf"] = newInterface         // make new interface a member of a mapped key "oneOf"

	itemsMap := make(map[string]interface{}) // initialize empty map
	itemsMap["items"] = oneofMap             // make oneof-map a member of a mapped key "items"

	return itemsMap
}

/// this function returns the keys of all the first level elements in the provided map
func getKeysOfMap(m map[string]interface{}) []string {
	var keys []string

	for key := range m {
		keys = append(keys, key)
	}

	return keys
}
