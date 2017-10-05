from confluent_kafka import Producer, Consumer, KafkaError
import json

producer = Producer({'bootstrap.servers': 'localhost:9092'})
consumer = Consumer({'bootstrap.servers': 'localhost:9092', 'group.id': 'mygroup',
                     'default.topic.config': {'auto.offset.reset': 'smallest'}})


def send_response(value):
    producer.produce("triplerout", value=json.dumps({'tripled': value}))
    producer.flush()


def send_error(message):
    producer.produce("triplerout", value=json.dumps({'error': message}))
    producer.flush()

consumer.subscribe(['triplerin'])
running = True
while running:
    msg = consumer.poll()
    if not msg.error():
        json_string = msg.value().decode('utf-8')
        print('[Tripler Microservice] Received message: %s' % json_string)
        try:
            request = json.loads(json_string)
            number = int(request['number'])
            send_response(number * 3)
        except:
            send_error("Malformed message")

    elif msg.error().code() != KafkaError._PARTITION_EOF:
        print(msg.error())
        running = False
consumer.close()