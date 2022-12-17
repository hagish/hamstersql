-- :name find_user :one
select * from users where user_id = :user_id

-- :name search_users :many
select * from users where username like :pattern

-- :name update_username :affected
update users set username = :username
where user_id = :user_id

-- :name get_username :scalar
select username from users where user_id = :user_id

-- :name update_username :insert
insert into users (username) values (:username)