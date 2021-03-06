import pymetamorph.metamorph as morph
from confluent_kafka import Producer, Consumer


class Service:
    def __init__(self):
        self._producer = Producer({'bootstrap.servers': 'localhost:9092'})
        self._consumer = Consumer({'bootstrap.servers': 'localhost:9092', 'group.id': 'mygroup',
                                   'default.topic.config': {'auto.offset.reset': 'smallest'}})
        self._received = []

    def terminate(self):
        self._consumer.close()
        self._consumer = None
        self._producer = None

    def subscribe_to(self, *topics):
        print("Topics: '{}'".format(list(topics)))
        self._consumer.subscribe(list(topics))

    def send_message(self, key, value, topic):
        print("Sending {}, {} to topic {}".format(key, value, topic))
        self._producer.produce(topic, value=value, key=key)
        self._producer.flush()

    def await_message(self, value, topic):
        for i in range(0, 10):
            msg = self._consumer.poll(timeout=1.0)
            print("Polled:", msg)
            if msg is None:
                continue
            if not msg.error():
                print("Topic: '{}', Value: '{}'".format(msg.topic(), msg.value()))
                value_as_string = msg.value().decode('UTF-8')
                if msg.topic() == topic and value_as_string == value:
                    return msg
                self._received.append(msg)
        raise RuntimeError("No message {} in topic {} was received".format(value, topic))


def before_feature(context, feature):
    context.metamorph = morph.Metamorph()
    context.metamorph.connect()


def before_scenario(context, scenario):
    context.metamorph.request_kafka_reset(["test_topic", "events"])
    context.metamorph.await_reset_complete()
    context.service = Service()


def after_scenario(context, scenario):
    context.service.terminate()
