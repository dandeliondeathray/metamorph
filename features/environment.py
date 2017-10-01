import pymetamorph.pymetamorph.metamorph as morph


class Service:
    def send_message(self, key, value, topic):
        print("Sending {}, {} to topic {}".format(key, value, topic))


def before_feature(context, feature):
    context.metamorph = morph.Metamorph()
    context.metamorph.connect()
    context.service = Service()
