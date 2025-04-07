package services

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

const SERVIVE_SCHEMA = `{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "title": "Service Schema",
  "description": "Schema for a service configuration",
  "type": "object",
  "properties": {
    "name": {
      "type": "string",
      "description": "Name of the service"
    },
    "description": {
      "type": "string",
      "description": "Description of the service"
    },
    "tailscaleProxy": {
      "type": "string",
      "description": "Indicates tailscale tunnel proxy path"
    },
    "hybrProxy": {
      "type": "string",
      "description": "Indicates hybr tunnel service path"
    },
    "variables": {
      "type": "object",
      "description": "Variables for the service templates",
      "additionalProperties": {
        "type": "array",
        "items": {
          "type": "object",
          "properties": {
            "name": {
              "type": "string",
              "description": "Name of the variable"
            },
            "default": {
              "type": "string",
              "description": "Default value of the variable"
            },
            "description": {
              "type": "string",
              "description": "Description of the variable"
            }
          },
          "required": [
            "name",
            "default",
            "description"
          ]
        }
      }
    },
    "templates": {
      "type": "array",
      "description": "List of template files for the service",
      "items": {
        "type": "string"
      }
    }
  },
  "required": [
    "name",
    "description",
    "hybrProxy",
    "tailscaleProxy",
    "variables",
    "templates"
  ],
  "additionalProperties": false
}`

func ValidateServiceJSON(path string) error {
	schemaLoader := gojsonschema.NewStringLoader(SERVIVE_SCHEMA)
	documentLoader := gojsonschema.NewReferenceLoader("file://" + path)

	result, err := gojsonschema.Validate(schemaLoader, documentLoader)
	if err != nil {
		return fmt.Errorf("Failed to validate service.json %s: %w", path, err)
	}

	if result.Valid() {
		return nil
	}

	errors := fmt.Sprintf("Invalid service.json file at: %s\n", path)
	for _, desc := range result.Errors() {
		errors += fmt.Sprintf("\t- %s\n", desc)
	}

	return fmt.Errorf(errors)
}

func ConfirmInvalidService(err error) bool {
	fmt.Println(err)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("Do you want to:")
		fmt.Println("  - Cancel and fix the errors (N)")
		fmt.Println("  - Continue without this service (Y)")
		fmt.Print("Enter your choice [Y/n]: ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(strings.ToLower(choice))

		switch strings.ToLower(choice) {
		case "n":
			fmt.Println("Please fix the errors in the service JSON file.")
			return false
		case "y":
			return true
		default:
			fmt.Println("Invalid choice. Please enter 'c' or 'y'.")
		}
	}
}
