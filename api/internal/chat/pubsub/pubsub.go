package pubsub

import "sync"

// Client is an in-memory Pub/Sub client.
type Client struct {
	topics map[Topic]map[Subscription]struct{}
	mu     sync.RWMutex
}

type Subscription struct {
	topic    Topic
	messages chan<- Message
}

type Topic = string

type Message = []byte

func New() *Client {
	return &Client{
		topics: make(map[Topic]map[Subscription]struct{}),
	}
}

// Publish publishes the message to all subscribers of the topic.
func (c *Client) Publish(topic Topic, message Message) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for subscription := range c.topics[topic] {
		select {
		case subscription.messages <- message:
			// ok
		default:
			// not ok, drop it
		}
	}
}

// Subscribe adds a new subscription for the topic.
// New messages will be published to the provided channel.
func (c *Client) Subscribe(topic Topic, messages chan<- Message) Subscription {
	c.mu.Lock()
	defer c.mu.Unlock()

	subscription := Subscription{
		topic:    topic,
		messages: messages,
	}
	if c.topics[topic] == nil {
		c.topics[topic] = make(map[Subscription]struct{})
	}
	c.topics[topic][subscription] = struct{}{}
	return subscription
}

// Unsubscribe removes the subscription.
// New messages will not be published to it anymore.
func (c *Client) Unsubscribe(subscription Subscription) {
	c.mu.Lock()
	defer c.mu.Unlock()

	subscriptions := c.topics[subscription.topic]
	delete(subscriptions, subscription)
	if len(subscriptions) == 0 && subscriptions != nil {
		delete(c.topics, subscription.topic)
	}
}
