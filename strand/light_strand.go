package strand

import (
	"errors"
	"fmt"
	"github.com/bcurren/go-hue"
	"strconv"
)

// Structure that holds the mapping from socket id to light id. This implements
// the hue.API interface so it can be used as a drop in replacement.
type LightStrand struct {
	api    hue.API
	Length int
	Lights *TwoWayMap
}

// Create a new light strand with the given length and hue.API to delegate to.
func NewLightStrand(length int, api hue.API) *LightStrand {
	var lightStrand LightStrand
	lightStrand.api = api
	lightStrand.Length = length
	lightStrand.Lights = NewTwoWayMap()

	return &lightStrand
}

// Change the hue.API delegate.
func (lg *LightStrand) SetDelegateAPI(api hue.API) {
	lg.api = api
}

// An interactive way of mapping all unmapped light bulbs on the hue bridge. This
// function does the following:
//
// 1. Turn all lights white
// 2. For each unmapped light
//   a. Turn the bulb red
//   b. Call socketToLightFunc - The implementation should return the socket id for 
//      the unmapped light
//   c. Map the bulb to the socket id
//   d. Turn the white bulb and continue
//
// This should be used to interactively prompt a person to map a light to a position
// in the strand.
func (lg *LightStrand) MapUnmappedLights(socketToLightFunc func() string) error {
	unmappedLightIds, err := lg.getUnmappedLightIds()
	if err != nil {
		return err
	}

	white := createWhiteLightState()
	red := createRedLightState()

	for _, unmappedLightId := range unmappedLightIds {
		// Turn new unmapped light red
		err = lg.api.SetLightState(unmappedLightId, red)
		if err != nil {
			return err
		}

		socketId := socketToLightFunc()
		if !lg.validSocketId(socketId) {
			return errors.New(fmt.Sprintf("Invalid socket id provided %s.", socketId))
		}
		lg.Lights.Set(socketId, unmappedLightId)

		// Turn newly mapped light white
		err = lg.api.SetLightState(unmappedLightId, white)
		if err != nil {
			return err
		}
	}

	return nil
}

func (lg *LightStrand) getUnmappedLightIds() ([]string, error) {
	allHueLights, err := lg.api.GetLights()
	if err != nil {
		return nil, err
	}

	allMappedLightIds := lg.Lights.GetValues()

	unmappedLights := make([]string, 0, 5)
	for _, hueLight := range allHueLights {
		alreadyMapped := false
		for _, mappedLightId := range allMappedLightIds {
			if hueLight.Id == mappedLightId {
				alreadyMapped = true
				break
			}
		}
		if !alreadyMapped {
			unmappedLights = append(unmappedLights, hueLight.Id)
		}
	}

	return unmappedLights, nil
}

func (lg *LightStrand) validSocketId(socketId string) bool {
	socketIdAsInt, err := strconv.Atoi(socketId)
	if err != nil {
		return false
	}

	if socketIdAsInt <= 0 || socketIdAsInt > lg.Length {
		return false
	}

	return true
}

func createWhiteLightState() *hue.LightState {
	white := &hue.LightState{}
	white.ColorTemp = new(uint16)
	*white.ColorTemp = 1800

	return white
}

func createRedLightState() *hue.LightState {
	red := &hue.LightState{}
	red.Brightness = new(uint8)
	*red.Brightness = 255
	red.Hue = new(uint16)
	*red.Hue = 65535
	red.Saturation = new(uint8)
	*red.Saturation = 255

	return red
}
