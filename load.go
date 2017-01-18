package multiarg

import (
	"encoding/json"
	"os"
	"reflect"
	"regexp"
	"strings"
)

// Config provides control over the loading process. This allows the location
// for CLI arguments and the JSON filenames to be customized, among other
// things.
type Config struct {
	Args          []string
	JSONFilenames []string
}

var (
	snakeCaseRegexp1 = regexp.MustCompile(`(.)([A-Z][a-z]+)`)
	snakeCaseRegexp2 = regexp.MustCompile(`([a-z0-9])([A-Z])`)
)

// snakeCase converts a name in CamelCase to snake_case.
// (Based on http://stackoverflow.com/a/1176023/193619.)
func snakeCase(name string) string {
	name = snakeCaseRegexp1.ReplaceAllString(name, `${1}_${2}`)
	name = snakeCaseRegexp2.ReplaceAllString(name, `${1}_${2}`)
	return strings.ToLower(name)
}

// assignJSONValue assigns an interface{} value to a struct if possible. Any
// errors that occur during assignment are ignored.
func assignJSONValue(dest reflect.Value, src interface{}) {
	defer func() { recover() }()
	dest.Set(reflect.ValueOf(src).Convert(dest.Type()))
}

// assignEnvValue attempts to assign the value of the appropriate environment
// variable to to the variable if it has a non-empty value.
func assignEnvValue(v reflect.Value, components []string) {
	name := strings.ToUpper(strings.Join(components, "_"))
	if x := os.Getenv(name); len(x) != 0 {
		json.Unmarshal([]byte(x), v.Addr().Interface())
	}
}

// walk loads values for a struct from JSON, env. variables, and CLI arguments.
func walk(v reflect.Value, vJSON interface{}, cliMap map[string]reflect.Value, components ...string) {
	// Dereference any pointers
	for v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}
	switch v.Kind() {
	// For a struct, loop through each of the fields
	case reflect.Struct:
		m, _ := vJSON.(map[string]interface{})
		for i := 0; i < v.NumField(); i++ {
			var (
				fElem = v.Field(i)
				fName = snakeCase(v.Type().Field(i).Name)
				fJSON interface{}
			)
			// Check the map for the JSON value
			if x, ok := m[fName]; ok {
				fJSON = x
			}
			walk(fElem, fJSON, cliMap, append(components, fName)...)
		}
	// For everything else, attempt to assign using the appropriate methods
	default:
		assignJSONValue(v, vJSON)
		assignEnvValue(v, components)
		cliMap["--"+strings.Join(components, "-")] = v
	}
}

// Enumerate the CLI arguments and assign to variables as necessary
func assignCLIValues(cliMap map[string]reflect.Value, args []string) {
	for i := 0; i < len(args); i++ {
		if strings.HasPrefix(args[i], "--") {
			v, ok := cliMap[args[i]]
			if !ok {
				continue
			}
			// Boolean values do not require an argument, but all other types
			// do; in that case, pull the argument that follows
			if v.Kind() == reflect.Bool {
				v.SetBool(true)
			} else {
				i++
				if i < len(args) {
					json.Unmarshal([]byte(args[i]), v.Addr().Interface())
				}
			}
		}
	}
}

// Load attempts to load application configuration from multiple sources. The
// v parameter should be a pointer to a struct and the config parameter
// determines behavior.
func Load(v interface{}, config *Config) {
	// Build a map from all of the JSON files
	jsonMap := make(map[string]interface{})
	if config.JSONFilenames != nil {
		for _, filename := range config.JSONFilenames {
			r, err := os.Open(filename)
			if err != nil {
				continue
			}
			json.NewDecoder(r).Decode(&jsonMap)
			r.Close()
		}
	}
	// Walk the struct, assigning to the CLI map along the way
	cliMap := make(map[string]reflect.Value)
	walk(reflect.ValueOf(v), jsonMap, cliMap)
	// Use os.Args if nothing was specified
	args := config.Args
	if args == nil {
		args = os.Args
	}
}
