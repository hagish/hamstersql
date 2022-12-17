from sqlalchemy import text

import core


def create_users_table(session):
    query = text("""
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    age INTEGER NOT NULL,
    email TEXT NOT NULL UNIQUE
);
""")
    params = {}
    result = session.execute(query, params)
    return result


def insert_user(session, name, age, email) -> int:
    query = text("""INSERT INTO users (name, age, email) VALUES (:name, :age, :email);""")
    params = {"name": name, "age": age, "email": email}
    result = session.execute(query, params)
    return core.convert_result_to_insert_id(result)


def find_user(session, name) -> dict[str, any]:
    query = text("""
SELECT * FROM users WHERE name=:name
    """)
    params = {"name": name}
    result = session.execute(query, params)
    return core.convert_result_to_one_row(result)


def find_user_name_by_id(session, id) -> any:
    query = text("""
SELECT name FROM users WHERE id=:id
    """)
    params = {"id": id}
    result = session.execute(query, params)
    return core.convert_result_to_scalar(result)


def list_users(session) -> list[dict[str, any]]:
    query = text("""
    SELECT * FROM users
        """)
    params = {}
    result = session.execute(query, params)
    return core.convert_result_to_many_rows(result)


def delete_users_with_age(session, age) -> int:
    query = text("""
    DELETE FROM users WHERE age=:age
    """)
    params = {"age": age}
    result = session.execute(query, params)
    return core.convert_result_to_affected_rows(result)