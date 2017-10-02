from behave import *
from hamcrest import *
import json
from pymetamorph.pymetamorph.metamorph import OnTopic, value_as_string


@when(u'the number 12 is sent to Doubler')
def step_impl(context):
    value = json.dumps({'number': 12})
    context.metamorph.send_message(topic="doublerin", value=value)


@then(u'it responds with the value 24')
def step_impl(context):
    # The `matcher` argument ensures that we only match messages in the `doublerout` topic.
    # One could create a matcher that inspects the value and finds a message with the exact value we're looking for.
    # That would actually be better, as this test would fail if any message came before the message we're waiting for.
    message_event = context.metamorph.await_message(matcher=OnTopic('doublerout'))

    # We have now received a message event on topic 'doublerout'. The event contains the message itself in the `message`
    # field, and the topic in the `topic` field. We can be sure the topic is `doublerout`, because of the matcher we
    # used.
    # Get the message value as a string, and decode it as JSON. The decoded JSON should be {'double': 24}.
    json_string = value_as_string(message_event)
    response = json.loads(json_string)
    assert_that(response['doubled'], equal_to(24))


@when(u'a malformed message is sent to Doubler')
def step_impl(context):
    malformed_message = json.dumps({'number': 'foo'})
    context.metamorph.send_message(topic="doublerin", value=malformed_message)


@then(u'it responds with an error message')
def step_impl(context):
    message_event = context.metamorph.await_message(matcher=OnTopic('doublerout'))
    json_string = value_as_string(message_event)
    response = json.loads(json_string)
    # Assert that the JSON message is on the form {'error': '<any message>'}.
    assert_that('error' in response)
