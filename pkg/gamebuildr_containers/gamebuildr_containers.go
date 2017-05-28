package gamebuildrContainers

import "errors"
import "strings"

var godotVersions = map[string]string{
	"2.1.2": "mr.robot-godot-2.1.2",
}

var defoldVersions = map[string]string{
	"1": "mr.robot-defold-1",
}

var containers = map[string]map[string]string{
	"godot":  godotVersions,
	"defold": defoldVersions,
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
