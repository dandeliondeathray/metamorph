"""
Helper functions for starting and stopping processes as part of an acceptance test.
"""

import subprocess
import os
import os.path
import tempfile


class Process:
    def __init__(self, popen, out, err):
        self._popen = popen
        self._out = out
        self._err = err

    def stop(self):
        self._popen.terminate()

    def out_name(self):
        return self._out.name

    def err_name(self):
        return self._err.name


def ensure_exists(path):
    try:
        os.makedirs(path)
    except OSError:
        pass


def start(go=None, args="", env=None):
    try:
        gopath = os.environ['GOPATH']
    except KeyError:
        gopath = os.path.join(os.environ['HOME'], 'go')

    ensure_exists('./behave_output')
    out = tempfile.NamedTemporaryFile(prefix=os.path.basename(go) + "_test_out", delete=False, dir='./behave_output')
    err = tempfile.NamedTemporaryFile(prefix=os.path.basename(go) + "_test_err", delete=False, dir='./behave_output')

    process_path = os.path.join(gopath, go)
    popen_obj = subprocess.Popen([process_path] + args.split(), stdout=out, stderr=err, env=env)
    return Process(popen_obj, out, err)
