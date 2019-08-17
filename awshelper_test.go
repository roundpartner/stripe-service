package main

import "testing"

func TestGetSession(t *testing.T) {
	session := GetSession()
	if session == nil {
		t.Errorf("AWS Session returned nil")
	}
}

func TestGetTopic(t *testing.T) {
	_, err := GetTopic()
	if err != nil {
		t.Errorf("Unexpected Error: %s", err.Error())
	}
}
