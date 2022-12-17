-- :name find_user :one
select * from users where name = :name

-- :name search_users :many
select * from users where name like :pattern

-- :name list_users :many
select * from users

-- :name update_name :affected
-- :doc Update the name of a user
update users set name = :name
where user_id = :user_id

-- :name get_name :scalar
-- :doc get the name for a given user_id
select name from users where user_id = :user_id

-- :name insert_user :insert
-- :doc insert a new user
INSERT INTO users (name, age, email) VALUES (:name, :age, :email);

-- :name create_users_table
-- :doc creates a user table
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    age INTEGER NOT NULL,
    email TEXT NOT NULL UNIQUE
);

-- :name find_user_name_by_id :scalar
-- :doc find a user's name by their id
select name from users where id = :id

-- :name delete_users_with_age :affected
-- :doc delete all users with a given age
delete from users where age = :age