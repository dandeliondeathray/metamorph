"""
Helper functions for starting and stopping processes as part of an acceptance test.
"""

import subprocess
import os
import os.path


class Process:
    def __init__(self, popen):
        self._popen = popen

    def stop(self):
        self._popen.terminate()


def start(go=None, args=""):
    try:
        gopath = os.environ['GOPATH']
    except KeyError:
        gopath = os.path.join(os.environ['HOME'], 'go')
    process_path = os.path.join(gopath, go)
    popen_obj = subprocess.Popen([process_path] + args.split())
    return Process(popen_obj)
