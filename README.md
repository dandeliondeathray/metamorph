# Metamorph
Metamorph is a testing framework aimed at testing microservices. It creates a Kafka broker, and lets
you verify that messages are consumed and produced as expected.

It is written in Go and has a Python interface.

## Metamorph is currently under development and does not function properly yet!
All of the documentation and examples indicate how things _will_ work.

# Usage overview
You have a microservice that is based on Kafka as a message queue. Let's say you named the service
Baz. You need to test that Baz behaves properly when connected to an actual Kafka system.

To do this you might create some tests written using Python unittest module, or maybe using 
behaviour driven development tools like Cucumber or Behave. The author of this document uses Behave,
so everything you see will be most easily integrated with it.

You will want to send messages to the Baz service, via a Kafka broker, and ensure that the Baz
service responds properly with messages of its own. This is what Metamorph helps you with.

Metamorph runs as a service and starts a Kafka broker (and the necessary Zookeeper instance), on
request. You connect to Metamorph with a WebSocket, and it will send a JSON message every time a
message is sent to the Kafka broker. Metamorph will only listen to topics you tell it to listen to.
The Python interface Pymetamorph contains a interface against this WebSocket and has functions that
allows you to wait for messages to arrive within a certain time period. 

Between each test, or between some groups of tests, you will want to reset the Kafka system that
Metamorph creates, so each test runs against a known state. We don't want your tests to interfere
with each other. Note that this implies you will also restart the Baz service.

A test will generally follow these steps:

- Start Metamorph once
- For each test case
    + Request that Metamorph (re)start the Kafka system
    + Wait for Kafka to start up
    + Start the microservice
    + Using the Pymetamorph client, send a message to some topic via the Kafka broker
    + Wait for the microservice to send a message to some topic
    + Verify that the message is as expected
    + Shut down the service
- Shut down Metamorph

# Example
The Doubler service in `examples/doubler` serves as an example of how to use Metamorph and
Pymetamorph with the BDD framework Behave.

The Tripler service, in `examples/tripler`, suspiciously similar to the Doubler service, serves as
a more straightforward, but unrealistic, example if you're unused to Behave or BDD frameworks.

# License
Apache License 2.0