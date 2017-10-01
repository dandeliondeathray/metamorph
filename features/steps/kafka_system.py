from behave import *
from hamcrest import *


@when(u'Pymetamorph sends a reset environment command')
def step_impl(context):
    context.metamorph.request_kafka_reset()


@then(u'we receive a reset complete event')
def step_impl(context):
    context.metamorph.await_reset_complete()
