package gamebuildrContainers

import "errors"
import "strings"

var godotVersions = map[string]string{
	"2.1": "mr.robot-godot-2.1.2",
}

var defoldVersions = map[string]string{
	"1": "mr.robot-defold-1",
}

var unrealVersions = map[string]string{
	"4.17": "unreal-4.17.2",
}

var containers = map[string]map[string]string{
	"godot engine": godotVersions,
	"defold":       defoldVersions,
	"unreal":       unrealVersions,
}

// GamebuildrContainers is responsible for docker container image management
type GamebuildrContainers struct {
}

// GetContainerImageName returns the image name for a given engine and version
func (gamebuildrContainers GamebuildrContainers) GetContainerImageName(
	engineName string, engineVersion string) (string, error) {
	lowerName := strings.ToLower(engineName)
	lowerVersion := strings.ToLower(engineVersion)
	image, exists := containers[lowerName][lowerVersion]
	if !exists {
		return "", errors.New("Container image not found")
	}

	return image, nil
}
