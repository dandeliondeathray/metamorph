Feature: Messages are sent as events over a WebSocket connection
	The EVENT INTERFACE is a WebSocket connection from the tester to Metamorph. Metamorph sends
	JSON messages as events over the connection whenever a message is received.

	Scenario: A message is sent and a corresponding event is produced
		 When a message is sent to Kafka
		 Then a message event is received on the event interface

	Scenario: Topic is included in a message event
		Given the topic "events"
		 When a message is sent to Kafka on that topic
		 Then the message event contains the topic