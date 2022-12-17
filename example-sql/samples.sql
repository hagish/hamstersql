-- :name find_user :one
select * from users where username = :username

-- :name search_users :many
select * from users where username like :pattern

-- :name update_username :affected
-- :doc Update the username of a user
update users set username = :username
where user_id = :user_id

-- :name get_username :scalar
-- :doc get the username for a given user_id
select name from users where user_id = :user_id

-- :name insert_user :insert
INSERT INTO users (name, age, email) VALUES (:name, :age, :email);

-- :name create_users_table
-- :doc creates a user table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    age INTEGER NOT NULL,
    email TEXT NOT NULL UNIQUE
);