from util_ft import *
from util_sklearn import *
from util_db import *

import dill


def dfunc(fn, *args, **kwargs):
    with open(fn, "rb") as fp:
        func = dill.load(fp)
    return func(*args, **kwargs)