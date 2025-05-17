package common

import (
	"log"
	"reflect"

	"github.com/BurntSushi/toml"
)

// LoadConfig dynamically loads TOML config into any struct
func LoadConfig[lTomlStruct any](filePath string) (lTomlStruct, error) {
	var lDecodeValues lTomlStruct

	// lData, lErr := os.ReadFile(filePath)
	// if lErr != nil {
	// 	log.Printf("Error reading TOML file %s: %v", filePath, lErr)
	// 	return lDecodeValues, lErr
	// }
	// log.Printf("TOML file content:\n%s", string(lData))

	if _, lErr := toml.DecodeFile(filePath, &lDecodeValues); lErr != nil {
		log.Printf("Error loading config from %s: %v", filePath, lErr)
		return lDecodeValues, lErr
	}

	log.Printf("Loaded config (%s): %+v", reflect.TypeOf(lDecodeValues).Name(), lDecodeValues)
	return lDecodeValues, nil
}
