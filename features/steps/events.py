from behave import *
from hamcrest import *
from pymetamorph.pymetamorph.metamorph import OnTopic


@when(u'a message is sent to Kafka')
def step_impl(context):
    context.topic = "test_topic"
    context.service.send_message("key", "value", context.topic)


@then(u'a message event is received on the event interface')
def step_impl(context):
    context.metamorph.await_message()


@given(u'the topic "{topic_name}"')
def step_impl(context, topic_name):
    context.topic = topic_name


@when(u'a message is sent to Kafka on that topic')
def step_impl(context):
    context.service.send_message("key", "value", context.topic)


@then(u'the message event contains the topic')
def step_impl(context):
    m = context.metamorph.await_message(OnTopic(context.topic))
    assert_that(m['topic'], equal_to(context.topic))
