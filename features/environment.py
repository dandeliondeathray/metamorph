import pymetamorph.pymetamorph.metamorph as morph
from confluent_kafka import Producer

class Service:
    def __init__(self):
        self._producer = Producer({'bootstrap.servers': 'localhost:9092'})

    def send_message(self, key, value, topic):
        print("Sending {}, {} to topic {}".format(key, value, topic))
        self._producer.produce(topic, value=value, key=key)
        self._producer.flush()


def before_feature(context, feature):
    context.metamorph = morph.Metamorph()
    context.metamorph.connect()


def before_scenario(context, scenario):
    context.metamorph.request_kafka_reset()
    context.metamorph.await_reset_complete()
    context.service = Service()
