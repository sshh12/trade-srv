from urllib.parse import urlparse
from collections import defaultdict
import pandas as pd
import psycopg2
import os

DATABASE_URI = os.environ.get("DATABASE_URI", "postgresql://postgres:postgres@127.0.0.1/dbname")


def connect():
    url = urlparse(DATABASE_URI)
    conn = psycopg2.connect(
        user=url.username,
        password=url.password,
        host=url.hostname,
        port=url.port,
        database=url.path[1:],
    )
    return conn


def fetch_all(query, *args):
    conn = connect()
    cur = conn.cursor()
    try:
        cur.execute(query, *args)
        return cur.fetchall()
    finally:
        conn.close()


def fetch_all_as_df(query, *args, columns=None):
    data = fetch_all(query, *args)
    if columns is None:
        columns = [chr(i + 97) for i in range(len(data[0]))]
    return pd.DataFrame(data, columns=columns)


def expand_df_dicts(df, dict_defs=None):
    if dict_defs is None:
        dict_defs = defaultdict(set)
        for _, row in df.iterrows():
            for cidx, val in enumerate(row):
                if isinstance(val, dict):
                    for key in val.keys():
                        dict_defs[cidx].add(key)
    org_cols = list(df.columns)
    new_cols = []
    new_rows = []
    for i in range(len(org_cols)):
        if i not in dict_defs:
            new_cols.append(org_cols[i])
        else:
            for k in dict_defs[i]:
                new_cols.append(org_cols[i] + "_" + k)
    for _, row in df.iterrows():
        new_row = []
        for cidx, val in enumerate(row):
            if cidx not in dict_defs:
                new_row.append(val)
            else:
                for k in dict_defs[cidx]:
                    new_row.append(val.get(k))
        new_rows.append(new_row)
    return pd.DataFrame(new_rows, columns=new_cols), dict_defs