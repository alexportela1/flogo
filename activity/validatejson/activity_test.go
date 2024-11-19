package validatejson

import (
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"

	"github.com/stretchr/testify/assert"
)

var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {

	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}

func TestCreate(t *testing.T) {

	act := NewActivity(getActivityMetadata())

	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
}

func TestEval(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	jsonSchema := `{"$id":"https://example.com/person.schema.json","$schema":"https://json-schema.org/draft/2020-12/schema","title":"Person","type":"object","properties":{"firstName":{"type":"string","description":"The person's first name."},"lastName":{"type":"string","description":"The person's last name."},"age":{"description":"Age in years which must be equal to or greater than zero.","type":"integer","minimum":0}}}`
	
	jsonStringSuccess := `{"firstName":"John","lastName":"Doe","age":21}`
	jsonStringErrorInvalidFieldValue := `{"firstName":"John","lastName":"Doe","age":"21"}`
	
	//success case
	tc.SetInput("jsonSchema", jsonSchema)
	tc.SetInput("jsonString", jsonStringSuccess)
	_, err := act.Eval(tc)
	assert.Nil(t, err)
	success := tc.GetOutput("success")
	assert.Equal(t, success, true)
	
	//error case - invalid field value
	tc.SetInput("jsonSchema", jsonSchema)
	tc.SetInput("jsonString", jsonStringErrorInvalidFieldValue)
	_, err = act.Eval(tc)
	assert.Nil(t, err)
	success = tc.GetOutput("success")
	assert.Equal(t, success, false)
}
