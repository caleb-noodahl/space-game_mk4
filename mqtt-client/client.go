package mqttclient

import (
	"fmt"

	mqtc "github.com/eclipse/paho.mqtt.golang"
)

type Client struct {
	client         mqtc.Client
	opts           *mqtc.ClientOptions
	onconnect      mqtc.OnConnectHandler
	lostconnection mqtc.ConnectionLostHandler
}

func NewClient(opts ...func(*Client)) *Client {
	c := &Client{
		opts: mqtc.NewClientOptions(),
		onconnect: func(c mqtc.Client) {
			fmt.Println("connected to mqtt server")
		},
		lostconnection: func(c mqtc.Client, err error) {
			fmt.Println("lostconn. err: %s", err)
		},
	}
	for _, o := range opts {
		o(c)
	}

	c.client = mqtc.NewClient(c.opts)
	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return c
}

func WithSubscriptionHandler(topic string, eos int, handler mqtc.MessageHandler) func(*Client) {
	return func(c *Client) {
		token := c.client.Subscribe(topic, 2, handler)
		if token.Wait() && token.Error() != nil {
			panic(token)
		}
	}
}

func WithServer(url string) func(*Client) {
	return func(c *Client) {
		c.opts.AddBroker(url)
	}
}

func WithBasicAuth(username, password string) func(*Client) {
	return func(c *Client) {
		c.opts.SetUsername(username)
		c.opts.SetPassword(password)
	}
}

func (c *Client) Publish(topic string, msg []byte) {
	token := c.client.Publish(topic, 2, false, msg)
	token.Wait()
}

func (c *Client) Subscribe(topic string, handler func(mqtc.Client, mqtc.Message)) {
	token := c.client.Subscribe(topic, 2, handler)
	if token.Wait() && token.Error() != nil {
		panic(token)
	}
}

func (c *Client) Unsubscribe(topics ...string) {
	c.client.Unsubscribe(topics...)
}

func (c *Client) Disconnect() {
	c.client.Disconnect(300)
}
