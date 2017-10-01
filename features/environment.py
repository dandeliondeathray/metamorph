import pymetamorph.pymetamorph.metamorph as morph


class Service:
    def send_message(self, key, value, topic):
        print("Sending {}, {} to topic {}".format(key, value, topic))


def before_feature(context, feature):
    context.metamorph = morph.Metamorph()
    context.metamorph.connect()


def before_scenario(context, scenario):
    context.metamorph.request_kafka_reset()
    context.metamorph.await_reset_complete()
    context.service = Service()
