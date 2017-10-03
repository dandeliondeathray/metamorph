"""
A straightforward example of how to use Metamorph to test a microservice.

Note that this test expects Metamorph to already be up and running.

Usage: python tripler_test.py

"""
import pymetamorph.metamorph as metamorph
import json
import subprocess
import time

# This is the client we'll use to tell Metamorph to start up a Kafka system, and for it to tell us whenever a message
# is received.
metamorph_client = metamorph.Metamorph()
metamorph_client.connect()

#
# Start the Kafka system
#

print("Starting the Kafka system...", end='', flush=True)
# Ask Metamorph to (re)start the Kafka system and listen to the specified topics.
metamorph_client.request_kafka_reset(["triplerout"])
# Wait for Metamorph to restart the Kafka system.
metamorph_client.await_reset_complete()
print(" DONE.")
print("")

#
# Start the Tripler service
#
print("Starting Tripler.")
tripler_process = subprocess.Popen(["python", "tripler.py"])

print("Give Tripler a few seconds to start up...", end='', flush=True)
time.sleep(3)
print(" DONE.")
print("")

#
# Send the number 12 to Tripler, using the JSON format that Tripler expects.
#
print("Sending 12 to Tripler...")
value = json.dumps({'number': 12})
metamorph_client.send_message(topic="triplerin", value=value)

#
# Expect to get 36 back
#

print("Expecting a response with value 36...")
print("")
# The `matcher` argument ensures that we only match messages in the `triplerout` topic.
# One could create a matcher that inspects the value and finds a message with the exact value we're looking for.
# That would actually be better, as this test would fail if any message came before the message we're waiting for.
message_event = metamorph_client.await_message(matcher=metamorph.OnTopic('triplerout'))

# We have now received a message event on topic 'triplerout'. The event contains the message itself in the `message`
# field, and the topic in the `topic` field. We can be sure the topic is `triplerout`, because of the matcher we
# used.
# Get the message value as a string, and decode it as JSON. The decoded JSON should be {'tripled': 36}.
json_string = metamorph.value_as_string(message_event)
response = json.loads(json_string)
if response['tripled'] == 36:
    print("Test SUCCESSFUL.")
else:
    print("Test FAILED: expected 36, got", json.dumps(response))

#
# Terminate Tripler
#
tripler_process.terminate()
