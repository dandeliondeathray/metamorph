Feature: Kafka and Zookeeper is maintained by Metamorph
    Metamorph will start Kafka, along with Zookeeper, behind the scenes so that the consumer or producer under test can
    connect.

    Scenario: Reset Kafka environment
         When Pymetamorph sends a reset environment command
         Then we receive a reset complete event