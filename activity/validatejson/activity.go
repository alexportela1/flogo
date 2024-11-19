package validatejson

import (
	"strings"
	"fmt"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"

	"github.com/xeipuuv/gojsonschema"
)

const (
	ivJsonSchema  = "jsonSchema"
	ivJsonString  = "jsonString"
	ovSuccess = "success"
)

// log is the default package logger
var log = logger.GetLogger("activity-validatejson")

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {

	log.Info("Executing Validate JSON activity")

	// Get the action
	jsonSchema := context.GetInput(ivJsonSchema).(string)
	jsonString := context.GetInput(ivJsonString).(string)

	err = Validate(jsonSchema, jsonString)
	result := true

	if err != nil {
		log.Error(err)
		result = false
	}

	// Set the output value in the context
	context.SetOutput(ovSuccess, result)

	log.Info("Validate JSON activity completed")
	return true, nil
}

func Validate(schema string, data string) error {
	ok, errs := validate(schema, data)
	if !ok {
		for _, v := range errs {
			log.Errorf("Json schema validation error [%s]", v)
		}
		return fmt.Errorf("Output validate failed")
	}

	return nil
}

func validate(schema, data string) (bool, []string) {
	schemaLoader := gojsonschema.NewStringLoader(schema)
	jsonDataLoader := gojsonschema.NewStringLoader(data)
	errors := []string{}
	result, err := gojsonschema.Validate(schemaLoader, jsonDataLoader)

	if err != nil {
		errors = append(errors, err.Error())
		return false, errors
	}

	if result.Valid() {
		log.Info("The document is valid")
		return true, nil
	} else {
		log.Info("The document is not valid")
		for _, desc := range result.Errors() {
			errString := desc.String()
			log.Error(errString)
			errString = strings.Replace(errString, ":", "->", -1)
			errors = append(errors, errString)
		}
		return false, errors
	}
}
