package detector

import (
	"encoding/json"

	"github.com/Azure/go-autorest/autorest"
)

type Detector struct {
	autorest.Response `json:"-"`
	ID                string     `json:"id,omitempty"`
	Name              string     `json:"name,omitempty"`
	Type              string     `json:"type,omitempty"`
	Location          string     `json:"location,omitempty"`
	Properties        Properties `json:"properties,omitempty"`
}

type Properties struct {
	DataProvidersMetadata interface{} `json:"dataProvidersMetadata,omitempty"`
	Dataset               []struct {
		RenderingProperties struct {
			Description interface{} `json:"description,omitempty"`
			Title       interface{} `json:"title,omitempty"`
			Type        int         `json:"type,omitempty"`
		} `json:"renderingProperties,omitempty"`
		Table struct {
			Columns []struct {
				ColumnName string      `json:"columnName,omitempty"`
				ColumnType interface{} `json:"columnType,omitempty"`
				DataType   string      `json:"dataType,omitempty"`
			} `json:"columns,omitempty"`
			Rows      [][]string `json:"rows,omitempty"`
			TableName string     `json:"tableName,omitempty"`
		} `json:"table,omitempty"`
	} `json:"dataset,omitempty"`
	Metadata struct {
		AnalysisTypes    []string      `json:"analysisTypes,omitempty"`
		Author           string        `json:"author,omitempty"`
		Category         string        `json:"category,omitempty"`
		Description      string        `json:"description,omitempty"`
		ID               string        `json:"id,omitempty"`
		Name             string        `json:"name,omitempty"`
		Score            int           `json:"score,omitempty"`
		SupportTopicList []interface{} `json:"supportTopicList,omitempty"`
		Type             string        `json:"type,omitempty"`
		TypeID           string        `json:"typeId,omitempty"`
	} `json:"metadata,omitempty"`
	Status struct {
		Message  interface{} `json:"message,omitempty"`
		StatusID int         `json:"statusId,omitempty"`
	} `json:"status,omitempty"`
	SuggestedUtterances interface{} `json:"suggestedUtterances,omitempty"`
}

func (d Detector) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	if d.ID != "" {
		objectMap["id"] = d.ID
	}
	return json.Marshal(objectMap)
}

func (d *Detector) UnmarshalJSON(body []byte) error {

	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "id":
			if v != nil {
				var id string
				err = json.Unmarshal(*v, &id)
				if err != nil {
					return err
				}
				d.ID = id
			}
		case "name":
			if v != nil {
				var name string
				err = json.Unmarshal(*v, &name)
				if err != nil {
					return err
				}
				d.Name = name
			}
		case "properties":
			if v != nil {
				var properties Properties
				err = json.Unmarshal(*v, &properties)
				if err != nil {
					return err
				}
				d.Properties = properties
			}
		}
	}

	return nil
}
