package orchestration

import (
	"fmt"
	"sync"
)

type EventType string

type EventChannelData struct {
	Data      string
	EventType EventType
	Extras    map[string]string
}

type EventSource interface {
	Start(doneChan <-chan struct{}, eventChan chan<- *EventChannelData)
}

type SubscriptionManager struct {
	mu              sync.Mutex
	eventSources    map[EventType]EventSource
	runningServices map[EventType]chan struct{}
	cache           map[EventType]*EventChannelData
	listeners       map[EventType][]chan *EventChannelData
}

func NewSubscriptionManager() *SubscriptionManager {
	return &SubscriptionManager{
		eventSources:    make(map[EventType]EventSource),
		runningServices: make(map[EventType]chan struct{}),
		cache:           make(map[EventType]*EventChannelData),
		listeners:       make(map[EventType][]chan *EventChannelData),
	}
}

func (sm *SubscriptionManager) RegisterEventSource(eventType EventType, source EventSource) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, ok := sm.eventSources[eventType]; ok {
		return fmt.Errorf("Event source already registered for event type %s\n", eventType)
	}

	fmt.Printf("Registering event source for %s\n", eventType)
	sm.eventSources[eventType] = source

	return nil
}

func (sm *SubscriptionManager) Subscribe(eventChan chan *EventChannelData, eventTypes ...EventType) (cleanup func()) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for _, eventType := range eventTypes {
		sm.listeners[eventType] = append(sm.listeners[eventType], eventChan)
		fmt.Printf("%d Listeners for %s\n", len(sm.listeners[eventType]), eventType)

		if _, ok := sm.runningServices[eventType]; !ok {
			sm.runningServices[eventType] = make(chan struct{})
			if eventType == SYSTEM_NOTIFICATION_EVENT {
				go startNotificationMonitor(sm, sm.runningServices[eventType])
			} else {
				go sm.startEventService(eventType, sm.runningServices[eventType])
			}
		}

		if sm.cache[eventType] != nil {
			go func() {
				eventChan <- sm.cache[eventType]
			}()
		}
	}

	return func() {
		for _, eventType := range eventTypes {
			sm.Unsubscribe(eventType, eventChan)
		}
		close(eventChan)
	}
}

func (sm *SubscriptionManager) Unsubscribe(eventType EventType, eventChan chan *EventChannelData) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if _, ok := sm.listeners[eventType]; !ok {
		return
	}

	for i, ch := range sm.listeners[eventType] {
		if ch == eventChan {
			sm.listeners[eventType] = append(sm.listeners[eventType][:i], sm.listeners[eventType][i+1:]...)
			break
		}
	}

	if len(sm.listeners[eventType]) == 0 {
		if stopChan, ok := sm.runningServices[eventType]; ok {
			stopChan <- struct{}{}
			close(stopChan)
			delete(sm.runningServices, eventType)
		}
	}

	fmt.Printf("%d Listeners for %s\n", len(sm.listeners[eventType]), eventType)
}

func (sm *SubscriptionManager) startEventService(eventType EventType, stopChan <-chan struct{}) {
	source, ok := sm.eventSources[eventType]
	if !ok {
		fmt.Printf("No event source found for %s\n", eventType)
		return
	}

	eventChan := make(chan *EventChannelData)
	go source.Start(stopChan, eventChan)

	for {
		select {
		case <-stopChan:
			return
		case event := <-eventChan:
			sm.mu.Lock()
			sm.cache[eventType] = event
			for _, listenerChan := range sm.listeners[eventType] {
				go func(ch chan *EventChannelData) {
					ch <- event
				}(listenerChan)
			}
			sm.mu.Unlock()
		}
	}
}

var (
	subscriptionManager *SubscriptionManager
)

func init() {
	subscriptionManager = NewSubscriptionManager()
}

func GetSubscriptionManager() *SubscriptionManager {
	return subscriptionManager
}

func GetSubscriptionManagerWithEventChan() (*SubscriptionManager, chan *EventChannelData) {
	ch := make(chan *EventChannelData)
	return subscriptionManager, ch
}
