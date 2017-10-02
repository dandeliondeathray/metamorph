from confluent_kafka import Producer, Consumer, KafkaError
import json

producer = Producer({'bootstrap.servers': 'localhost:9092'})
consumer = Consumer({'bootstrap.servers': 'localhost:9092', 'group.id': 'mygroup',
                     'default.topic.config': {'auto.offset.reset': 'smallest'}})


def send_response(value):
    producer.produce("doublerout", value=json.dumps({'number': value}))
    producer.flush()


def send_error(message):
    producer.produce("doublerout", value=json.dumps({'error': message}))
    producer.flush()

consumer.subscribe(['doublerin'])
running = True
while running:
    msg = consumer.poll()
    if not msg.error():
        json_string = msg.value().decode('utf-8')
        print('Received message: %s' % json_string)
        try:
            request = json.loads(json_string)
            number = int(request['number'])
            send_response(number * 2)
        except:
            send_error("Malformed message")

    elif msg.error().code() != KafkaError._PARTITION_EOF:
        print(msg.error())
        running = False
consumer.close()