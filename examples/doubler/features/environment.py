#import pymetamorph.pymetamorph.metamorph as metamorph


def before_scenario(context, scenario):
    pass
    #context.metamorph = metamorph.Metamorph()

    # Tell Metamorph to reset the Kafka system.
    # Provide it with a list of topics it should listen to.
    #context.metamorph.request_kafka_reset(["doublerout"])

    # Wait for metamorph to restart the Kafka system.
    #context.metamorph.await_reset_complete()