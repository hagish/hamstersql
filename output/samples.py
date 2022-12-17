from sqlalchemy import text

import core


# OneRow: 
def find_user(session, username):
    query = text("""
select * from users where username = :username
    """)
    params = {
        "username": username,
    }
    result = session.execute(query, params)
    return core.convert_result_to_one_row(result)


# ManyRows: 
def search_users(session, pattern):
    query = text("""
select * from users where username like :pattern
    """)
    params = {
        "pattern": pattern,
    }
    result = session.execute(query, params)
    return core.convert_result_to_many_rows(result)


# AffectedRows: Update the username of a user
def update_username(session, username, user_id):
    query = text("""
update users set username = :username
where user_id = :user_id
    """)
    params = {
        "username": username,
        "user_id": user_id,
    }
    result = session.execute(query, params)
    return core.convert_result_to_affected_rows(result)


# Scalar: get the username for a given user_id
def get_username(session, user_id):
    query = text("""
select name from users where user_id = :user_id
    """)
    params = {
        "user_id": user_id,
    }
    result = session.execute(query, params)
    return core.convert_result_to_scalar(result)


# InsertID: 
def insert_user(session, name, age, email):
    query = text("""
INSERT INTO users (name, age, email) VALUES (:name, :age, :email);
    """)
    params = {
        "name": name,
        "age": age,
        "email": email,
    }
    result = session.execute(query, params)
    return core.convert_result_to_insert_id(result)


# None: creates a user table
def create_users_table(session):
    query = text("""
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    age INTEGER NOT NULL,
    email TEXT NOT NULL UNIQUE
);
    """)
    params = {
    }
    result = session.execute(query, params)
    return result

