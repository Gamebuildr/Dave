package gamebuildrContainers

import "testing"

func TestGetContainerImageNameGodot212(t *testing.T) {
	gamebuildrContainers := GamebuildrContainers{}
	image, err := gamebuildrContainers.GetContainerImageName("godot engine", "2.1")

	if err != nil {
		t.Fatalf(err.Error())
	}

	expectedImageName := "mr.robot-godot-2.1.2"
	if image != expectedImageName {
		t.Errorf("Expected image to be %v but got %v", expectedImageName, image)
	}
}

func TestGetContainerImageNameDefold1(t *testing.T) {
	gamebuildrContainers := GamebuildrContainers{}
	image, err := gamebuildrContainers.GetContainerImageName("defold", "1")

	if err != nil {
		t.Fatalf(err.Error())
	}

	expectedImageName := "mr.robot-defold-1"
	if image != expectedImageName {
		t.Errorf("Expected image to be %v but got %v", expectedImageName, image)
	}
}

func TestGetContainerImageNameNotExistantVersion(t *testing.T) {
	gamebuildrContainers := GamebuildrContainers{}
	_, err := gamebuildrContainers.GetContainerImageName("godot engine", "test")

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestGetContainerImageNameNotExistantName(t *testing.T) {
	gamebuildrContainers := GamebuildrContainers{}
	_, err := gamebuildrContainers.GetContainerImageName("test", "test")

	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
