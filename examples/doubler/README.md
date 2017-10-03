# Doubler
Doubler is an example of how to use Metamorph. It is a simple Kafka-based microservice,
which consumes a number in one topic, and produces the doubled number in another topic.

It demonstrates the use of the Metamorph service, which creates a Kafka environment,
and the Pymetamorph Python interface to Metamorph.

## Testing
The testing is using the BDD framework Behave. The requirements for Doubler are
written in Gherkin, and are available in the `features/specification.feature`.

The actual code for the tests is in `features/steps/specification.py`.

## Messages
It listens to the topic `doublerin`, to messages looking like

    { "number": 12 }
    
and it produces a message to the topic `doublerout`, looking like

    { "doubled": 24 }

If there is an error, such as a malformed message, then an error message is
produced in topic `doublerout`

    { "error": "Malformed message" }

## Demonstration
The Kafka system started by Metamorph is reset on every test scenario. This means
that a completely new system is created. This is for demonstration purposes. You can
control when you need the system reset. In the same way, we start a completely new
instance of the Doubler microservice on each scenario. This is done in the function
`before_scenario` in the `features/environment.py` file.

## Prerequisites
You will need a recent version of `librdkafka` installed on your system. At the time of this
writing, the version in Ubuntu is not recent enough to work with the `confluent-kafka`.

To install the prerequisites for executing the BDD specification, in its own virtualenv, run

    $ virtualenv venv/doublerspec
    $ . venv/doublerspec/bin/activate
    $ pip install -r features/requirements.txt

## Running
In the `doubler` directory, run

    $ behave

It should run a single feature with two scenarios.