package multi

import (
	"github.com/bcurren/go-hue"
	"time"
)

type rGetLights struct {
	lights []hue.Light
	err error
}

type rGetNewLights struct {
	lights []hue.Light
	lastScan time.Time
	err error
}

// GetLights() is same as hue.User.GetLights() except all light ids are mapped to
// socket ids.
func (m *MultiAPI) GetLights() ([]hue.Light, error) {
	c := make(chan rGetLights)
	for _, api := range m.apis {
		go gGetLights(c, api)
	}
	return lGetLights(c, len(m.apis))
}

func gGetLights(c chan rGetLights, api hue.API) {
	lights, err := api.GetLights()
	c <- rGetLights{mapLightIds(lights), err}
}

func lGetLights(c chan rGetLights, nResponses int) ([]hue.Light, error) {
	lErrors := make([]error, 0, 1)
	lLights := make([][]hue.Light, 0, nResponses)
	
	for i := 0; i < nResponses; i++ {
		result := <- c
		if result.err != nil {
			lErrors = append(lErrors, result.err)
		} else {
			lLights = append(lLights, result.lights)
		}
	}
	
	return mergeLights(lLights), mergeErrors(lErrors)
}

// GetNewLights() is same as hue.User.GetNewLights() except all light ids are mapped to
// socket ids.
func (m *MultiAPI) GetNewLights() ([]hue.Light, time.Time, error) {
	c := make(chan rGetNewLights)
	for _, api := range m.apis {
		go gGetNewLights(c, api)
	}
	return lGetNewLights(c, len(m.apis))
}

func gGetNewLights(c chan rGetNewLights, api hue.API) {
	lights, lastScan, err := api.GetNewLights()
	c <- rGetNewLights{mapLightIds(lights), lastScan, err}
}

func lGetNewLights(c chan rGetNewLights, nResponses int) ([]hue.Light, time.Time, error) {
	lErrors := make([]error, 0, 1)
	lLights := make([][]hue.Light, 0, nResponses)
	lLastScan := make([]time.Time, 0, nResponses)
	
	for i := 0; i < nResponses; i++ {
		result := <- c
		if result.err != nil {
			lErrors = append(lErrors, result.err)
		} 
		if result.lights != nil {
			lLights = append(lLights, result.lights)
		}
		lLastScan = append(lLastScan, result.lastScan)
	}
	
	return mergeLights(lLights), mergeTime(lLastScan), mergeErrors(lErrors)
}

// SearchForNewLights() is same as hue.User.SearchForNewLights() except all light ids are mapped to
// socket ids.
func (m *MultiAPI) SearchForNewLights() error {
	return nil
}

// GetLightAttributes() is same as hue.User.GetLightAttributes() except all light ids are mapped to
// socket ids.
func (m *MultiAPI) GetLightAttributes(socketId string) (*hue.LightAttributes, error) {
	return nil, nil
}

// SetLightName() is same as hue.User.SetLightName() except all light ids are mapped to
// socket ids.
func (m *MultiAPI) SetLightName(socketId string, name string) error {
	return nil
}

// SetLightState() is same as hue.User.SetLightState() except all light ids are mapped to
// socket ids.
func (m *MultiAPI) SetLightState(socketId string, state *hue.LightState) error {
	return nil
}

func mergeLights(lLights [][]hue.Light) []hue.Light {
	countOfLights := 0
	for _, lights := range lLights {
		countOfLights += len(lights)
	}
	
	mLights := make([]hue.Light, countOfLights)
	copyTo := 0
	for _, lights := range lLights {
		copy(mLights[copyTo:], lights)
		copyTo += len(lights)
	}
	
	return mLights
}

func mergeErrors(lErrors []error) error {
	return nil
}

func mergeTime(lTime []time.Time) time.Time {
	return lTime[0]
}

func mapLightIds(lLights []hue.Light) []hue.Light {
	return lLights
}