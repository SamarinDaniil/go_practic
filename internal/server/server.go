package server

import (
	"context"
	"fmt"
	"sync"

	eventpb "./pkg/api/protobuf"
)

const (
	port = ":5000"
)

var ev sync.Mutex

type server struct {
	eventMap map[int64]eventpb.Event
}

func (s *server) MakeEvent(ctx context.Context, req *eventpb.MakeEventRequest) (*eventpb.MakeEventResponse, error) {
	ev.Lock()

	if s.eventMap == nil {
		s.eventMap = make(map[int64]eventpb.Event)
	}

	eventID := len(s.eventMap) + 1
	event := eventpb.Event{
		SenderID: req.SenderID,
		EventID:  int64(eventID),
		Time:     req.Time,
		Name:     req.Name,
	}
	s.eventMap[int64(eventID)] = event

	fmt.Printf("Received event: %+v\n", event)
	ev.Unlock()
	return &eventpb.MakeEventResponse{
		EventID: int64(eventID),
	}, nil
}

func (s *server) GetEvent(ctx context.Context, req *eventpb.GetEventRequest) (*eventpb.GetEventResponse, error) {
	ev.Lock()
	if event, ok := s.eventMap[req.EventID]; ok {
		ev.Unlock()
		return &eventpb.GetEventResponse{
			Status: 0,
			Event:  &event,
		}, nil
	}
	ev.Unlock()

	return &eventpb.GetEventResponse{
		Status:       1,
		ErrorMessage: "Event not found",
	}, nil
}

func (s *server) DeleteEvent(ctx context.Context, req *eventpb.DeleteEventRequest) (*eventpb.DeleteEventResponse, error) {
	ev.Lock()
	if _, ok := s.eventMap[req.EventID]; ok {
		delete(s.eventMap, req.EventID)
		ev.Unlock()
		return &eventpb.DeleteEventResponse{
			Status: 0,
		}, nil
	}
	ev.Unlock()
	return &eventpb.DeleteEventResponse{
		Status:       1,
		ErrorMessage: "Event not found",
	}, nil
}

func (s *server) GetEvents(req *eventpb.GetEventsRequest, stream eventpb.EventManager_GetEventsServer) error {
	ev.Lock()
	for _, event := range s.eventMap {
		if event.Time >= req.ToTime && event.Time <= req.FromTime {
			err := stream.Send(&eventpb.GetEventsResponse{
				Events: []*eventpb.Event{&event},
			})
			ev.Unlock()
			if err != nil {
				return err
			}
		}
	}
	ev.Unlock()
	return nil
}
