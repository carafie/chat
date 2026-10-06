package pubsub

import (
	"slices"
	"testing"
)

func TestPubSub_Isolation(t *testing.T) {
	ps := New()

	topicA := Topic("topicA")
	chA := make(chan Message, 1)
	subA := ps.Subscribe(topicA, chA)
	defer ps.Unsubscribe(subA)

	topicB := Topic("topicB")
	chB := make(chan Message, 1)
	subB := ps.Subscribe(topicB, chB)
	defer ps.Unsubscribe(subB)

	messageA := Message("messageA")
	ps.Publish(topicA, messageA)

	select {
	case <-chA:
	default:
		t.Error("did not receive expected message")
	}

	select {
	case <-chB:
		t.Error("received unexpected message")
	default:
	}
}

func TestPubSub_SingleSubscriber(t *testing.T) {
	ps := New()

	topic := Topic("topic")
	ch := make(chan Message, 1)
	sub := ps.Subscribe(topic, ch)
	defer ps.Unsubscribe(sub)

	message := Message("message")
	ps.Publish(topic, message)

	select {
	case gotMessage := <-ch:
		if !slices.Equal(gotMessage, message) {
			t.Errorf("got %q, want %q", gotMessage, message)
		}
	default:
		t.Error("did not receive expected message")
	}
}

func TestPubSub_MultipleSubscribers(t *testing.T) {
	ps := New()

	topic := Topic("topicA")

	chA := make(chan Message, 1)
	subA := ps.Subscribe(topic, chA)
	defer ps.Unsubscribe(subA)

	chB := make(chan Message, 1)
	subB := ps.Subscribe(topic, chB)
	defer ps.Unsubscribe(subB)

	message := Message("message")
	ps.Publish(topic, message)

	select {
	case <-chA:
	default:
		t.Error("did not receive expected message")
	}

	select {
	case <-chB:
	default:
		t.Error("did not receive expected message")
	}
}

func TestPubSub_Unsubscribe(t *testing.T) {
	ps := New()

	topic := Topic("topic")
	message := Message("message")

	chA := make(chan Message, 1)
	subA := ps.Subscribe(topic, chA)

	chB := make(chan Message, 1)
	subB := ps.Subscribe(topic, chB)

	ps.Unsubscribe(subA)
	if len(ps.topics) != 1 {
		t.Fatal("expected only one topic")
	}
	if len(ps.topics[topic]) != 1 {
		t.Fatal("expected only one subscriber")
	}

	ps.Publish(topic, message)
	select {
	case <-chA:
		t.Fatal("received unexpected message")
	default:
	}

	ps.Unsubscribe(subB)
	if len(ps.topics) != 0 {
		t.Error("did not expect any topics")
	}

	// Unsubscribing multiple times should not fail:
	ps.Unsubscribe(subB)
}

func TestPubSub_NonBlockingDrop(t *testing.T) {
	ps := New()

	topic := Topic("topic")
	ch := make(chan Message, 1)
	sub := ps.Subscribe(topic, ch)
	defer ps.Unsubscribe(sub)

	message := Message("message")
	ps.Publish(topic, message)
	droppedMessage := Message("dropped message")
	ps.Publish(topic, droppedMessage)

	select {
	case gotMessage := <-ch:
		if !slices.Equal(gotMessage, message) {
			t.Fatalf("got %q, want %q", gotMessage, message)
		}
	default:
		t.Fatal("did not receive expected message")
	}

	select {
	case <-ch:
		t.Fatal("received unexpected message")
	default:
	}
}
