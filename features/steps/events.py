from behave import *
from hamcrest import *
from pymetamorph.metamorph import OnTopic


@when(u'a message is sent by the service to Kafka')
def step_impl(context):
    context.topic = "test_topic"
    context.service.send_message("key", "value", context.topic)


@then(u'a message event is received on the event interface')
def step_impl(context):
    context.metamorph.await_message()


@given(u'the topic "{topic_name}"')
def step_impl(context, topic_name):
    context.topic = topic_name


@when(u'a message is sent by the service to Kafka on that topic')
def step_impl(context):
    context.service.send_message("key", "value", context.topic)


@then(u'the message event contains the topic')
def step_impl(context):
    m = context.metamorph.await_message(OnTopic(context.topic))
    assert_that(m['topic'], equal_to(context.topic))


@when(u'a message event is sent from the test to Metamorph on topic "servicetopic"')
def step_impl(context):
    context.topic = "servicetopic"
    context.value = "Some value"
    context.metamorph.send_string_message(topic=context.topic, value=context.value)


@then(u'the message can be consumed by the service')
def step_impl(context):
    context.service.await_message(context.value, context.topic)

@given(u'the service subscribes to the topic "servicetopic"')
def step_impl(context):
    context.service.subscribe_to("servicetopic")