from util_ft import *
from util_sklearn import *
from util_db import *

import dill


def dload(fn):
    with open(fn, "rb") as fp:
        return dill.load(fp)


def dsave(obj, fn):
    with open(fn, "wb") as fp:
        dill.dump(obj, fp)


def dfunc(fn, *args, **kwargs):
    func = dload(fn)
    return func(*args, **kwargs)