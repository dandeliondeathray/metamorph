Feature: Doubler
  An example of how to use BDD and Metamorph for acceptance testing of a microservice.
  The Doubler service takes a message with a number, and doubles that number.
  It consumes messages from the `doublerin` topic and produces messages on the `doublerout` topic.

  Scenario: Double a number
     When the number 12 is sent to Doubler
     Then it responds with the value 24

  Scenario: Malformed messages
     When a malformed message is sent to Doubler
     Then it responds with an error message
