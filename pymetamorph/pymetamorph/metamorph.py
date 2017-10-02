"""

"""

import json
import asyncio
import websockets
import base64


class Metamorph:
    def __init__(self):
        self.ws = None
        self.received = []

    def connect(self, url=None):
        url = url if url else "localhost:23572"
        loop = asyncio.get_event_loop()
        self.ws = loop.run_until_complete(websockets.connect('ws://{}'.format(url)))

    def request_kafka_reset(self, topics):
        reset_event = {'type': 'reset_kafka_system', 'topics': topics}
        loop = asyncio.get_event_loop()
        loop.run_until_complete(self.ws.send(json.dumps(reset_event)))

    def await_message(self, matcher=None):
        return self._await_event("message", matcher)

    def await_reset_complete(self):
        return self._await_event("reset_complete", matcher=None, timeout=60.0)

    def _await_event(self, event_type, matcher=None, timeout=5.0):
        loop = asyncio.get_event_loop()
        return loop.run_until_complete(self._timed_await_event(event_type, matcher, timeout))

    def _timed_await_event(self, event_type, matcher, timeout):
        message = yield from asyncio.wait_for(self._do_await_event(event_type, matcher), timeout)
        return message

    async def _do_await_event(self, event_type, matcher_arg=None):
        matcher = matcher_arg if matcher_arg else Any()
        for i in range(len(self.received)):
            event = self.received[i]
            if event['type'] == event_type and matcher.matches(event):
                self.received = self.received[:i] + self.received[i+1:]
                return event
        return await self._read_events_until_type(event_type, matcher)

    async def _read_events_until_type(self, event_type, matcher):
        while True:
            raw_event = await self.ws.recv()
            print("  Raw:", raw_event)
            event = json.loads(raw_event)
            if event['type'] == 'error':
                print("Received error from Metamorph: {}".format(event['message']))
            if event['type'] == event_type and matcher.matches(event):
                print("  Yes, right type")
                return event
            else:
                self.received.append(event)


class Any:
    def matches(self, m): return True

    def __repr__(self):
        return "Any()"

    def __str__(self):
        return repr(self)


class MatchThese:
    def __init__(self, *matchers):
        self._matchers = matchers

    def matches(self, m):
        return all([matcher.matches(m) for matcher in self._matchers])

    def __repr__(self):
        return 'MatchThese({})'.format(repr(self._matchers))

    def __str__(self):
        return repr(self)

class OnTopic:
    def __init__(self, topic):
        self._topic = topic

    def matches(self, m):
        return 'topic' in m and m['topic'] == self._topic

    def __repr__(self):
        return "OnTopic('{}')".format(self._topic)

    def __str__(self):
        return repr(self)


def value_as_string(event):
    decoded = base64.decodebytes(event['message'])
    return decoded.decode('UTF-8')