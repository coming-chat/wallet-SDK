package customscale

import (
	"bytes"
	"fmt"
	"github.com/centrifuge/go-substrate-rpc-client/v4/scale"
	"github.com/centrifuge/go-substrate-rpc-client/v4/types"
)

type EventRaw []byte

type Event struct {
	Phase     *types.Phase
	EventData *EventData
	Topics    []*types.Hash
}

type EventData struct {
	ModuleName types.Text
	EventName  types.Text
	Args       []*CallArg
}

func (e *EventRaw) DecodeRaw(metadata *types.Metadata) ([]*Event, error) {
	var eventList []*Event
	decoder := scale.NewDecoder(bytes.NewReader(*e))

	// determine number of events
	n, err := decoder.DecodeUintCompact()
	if err != nil {
		return nil, err
	}
	for i := uint64(0); i < n.Uint64(); i++ {
		event := &Event{
			Phase:     &types.Phase{},
			EventData: &EventData{},
		}
		// decode Phase
		err = decoder.Decode(event.Phase)
		if err != nil {
			return nil, fmt.Errorf("unable to decode Phase for event #%v: %v", i, err)
		}

		// decode EventID
		id := types.EventID{}
		err = decoder.Decode(&id)
		if err != nil {
			return nil, fmt.Errorf("unable to decode EventID for event #%v: %v", i, err)
		}

		var argField []types.Si1Field
		// ask metadata for method & event name for event
		event.EventData.ModuleName, event.EventData.EventName, argField, err = FindEventNamesForEventID(metadata, id)
		// moduleName, eventName, err := "System", "ExtrinsicSuccess", nil
		if err != nil {
			return nil, fmt.Errorf("unable to find event with EventID %v in metadata for event #%v: %s", id, i, err)
		}

		argDecoder := ArgDecoder{decoder}
		event.EventData.Args, err = ArgDecode(metadata, &argDecoder, argField)
		if err != nil {
			return nil, err
		}

		err = decoder.Decode(&event.Topics)
		if err != nil {
			return nil, err
		}

		eventList = append(eventList, event)
	}
	return eventList, nil
}
