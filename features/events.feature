Feature: Messages are sent as events over a WebSocket connection
  METAMORPH is a testing framework which enables a test to verify that a microservice performs according to
  specification, when connected to an actual Kafka broker.
  Metamorph acts as a service which sets up a Kafka broker and consumes and produces messages on topics.
  The microservice under test connects directly to the Kafka broker, while the test itself connects to Metamorphs
  event interface.
  The EVENT INTERFACE is a WebSocket connection from the tester to Metamorph. Consumed messages are sent as JSON events
  of the event interface. A test can produce messages to the Kafka broker by sending an event to Metamorph over the
  interface.

  Scenario: A message is sent and a corresponding event is produced
     When a message is sent by the service to Kafka
     Then a message event is received on the event interface

  Scenario: A message is sent from the test
     When a message event is sent from the test to Metamorph
     Then the message can be consumed the service

  Scenario: Topic is included in a message event
    Given the topic "events"
     When a message is sent by the service to Kafka on that topic
     Then the message event contains the topic