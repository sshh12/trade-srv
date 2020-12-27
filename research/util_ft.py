import featuretools as ft
import numpy as np
from featuretools.primitives import make_trans_primitive
from featuretools.variable_types import Variable, Boolean, Numeric

Log10Sign = make_trans_primitive(
    function=lambda array: np.sign(array.astype(np.float64)) * np.log10(np.abs(array.astype(np.float64))),
    input_types=[Numeric],
    return_type=Numeric,
    name="log10_sign",
)
IsZero = make_trans_primitive(
    function=lambda array: (array == 0),
    input_types=[Numeric],
    return_type=Boolean,
    name="is_zero",
)
IsNaN = make_trans_primitive(
    function=lambda array: array.isna(),
    input_types=[Numeric],
    return_type=Boolean,
    name="is_nan",
)


def ft_encode_data(key, feat_defs, df):
    if isinstance(feat_defs, str):
        feat_defs = ft.load_features(feat_defs)
    ent_set = ft.EntitySet(
        key + "-data",
        {
            key + "s": (df, key),
        },
        [],
    )
    return (
        ft.calculate_feature_matrix(
            feat_defs,
            ent_set,
            verbose=True,
            n_jobs=-1,
        )
        .replace([np.inf, -np.inf], 0)
        .fillna(0)
    )