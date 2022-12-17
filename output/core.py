from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker


def convert_result_to_one_row(result):
    row = result.fetchone()
    if row:
        row = dict(zip(result.keys(), row))
        return row
    return None


def convert_result_to_scalar(result):
    row = result.fetchone()
    if row:
        return row[0]
    return None


def convert_result_to_many_rows(result):
    rows = result.fetchall()
    rows_mapped = []
    for row in rows:
        rows_mapped.append(dict(zip(result.keys(), row)))
    return rows_mapped


def convert_result_to_insert_id(result):
    return result.lastrowid


def convert_result_to_affected_rows(result):
    return result.rowcount


def make_session(uri):
    engine = create_engine(uri, echo=True)
    session_fun = sessionmaker(bind=engine)
    session = session_fun()
    return session