package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/aido123/detector/pkg/containerservice/detector"
)

type FakeDetectorClient struct {
}

func (detectorClient FakeDetectorClient) Get(ctx context.Context, resourceGroupName string, resourceName string, detectId string, startTime string, endTime string) (result detector.Detector, err error) {
	jsonFile, err := os.Open("test_data/fake_detector.json")
	if err != nil {
		return
	}
	defer jsonFile.Close()
	jsonData, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return
	}
	var detector detector.Detector
	json.Unmarshal(jsonData, &detector)

	return detector, nil
}

func TestGetDetectionLines(t *testing.T) {
	var detectorClient FakeDetectorClient
	detectionLine := getDetectionLines("aad-issues", detectorClient)
	if detectionLine[0].Status != "Success" {
		t.Error("Expected Status to be Success but found " + detectionLine[0].Status)
	}
	if detectionLine[0].Priority != 1 {
		t.Errorf("Expected Priority to be 1 but found %d", detectionLine[0].Priority)
	}

	if detectionLine[1].Status != "Info" {
		t.Error("Expected Status to be Success but found " + detectionLine[1].Status)
	}
	if detectionLine[1].Priority != 2 {
		t.Errorf("Expected Priority to be 2 but found %d", detectionLine[1].Priority)
	}

	if detectionLine[2].Status != "Warning" {
		t.Error("Expected Status to be Success but found " + detectionLine[2].Status)
	}
	if detectionLine[2].Priority != 3 {
		t.Errorf("Expected Priority to be 3 but found %d", detectionLine[2].Priority)
	}

	if detectionLine[3].Status != "Critical" {
		t.Error("Expected Status to be Success but found " + detectionLine[3].Status)
	}
	if detectionLine[3].Priority != 4 {
		t.Errorf("Expected Priority to be 4 but found %d", detectionLine[3].Priority)
	}

	if detectionLine[4].Status != "Other" {
		t.Error("Expected Status to be Success but found " + detectionLine[4].Status)
	}
	if detectionLine[4].Priority != 0 {
		t.Errorf("Expected Priority to be 0 but found %d", detectionLine[4].Priority)
	}
}
