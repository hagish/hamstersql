hamstersql: A tool to turn sql into nicely named functions. No magic.
===================================================================

A tool that generate nice functions from sql queries based on go templates. It is inspired by
[PugSQL](https://pugsql.org/).

It turns this:

```sql
-- :name search_users :many
select * from users where username like :pattern

-- :name update_username :affected
-- :doc Update the username of a user
update users set username = :username
where user_id = :user_id
```
into this:

```python
... snip ... core libraries and imports ... snip ...

# many: 
def search_users(session, pattern):
    query = text("""
select * from users where username like :pattern
    """)
    params = {
        "pattern": pattern,
    }
    result = session.execute(query, params)
    return core.convert_result_to_many_rows(result)


# affected: Update the username of a user
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
```

The goal was to have the sql queries in a readable way at one place and 
use nice readable function names in the code. 

* Everything is explicit.
* It integrates well into any build pipeline and only touches things when changed.
* It is easy to adjust to your needs. It uses the go template engine to generate the code.
* Each sql file is a group of queries that turns into a generated code file.
* At the moment there is only a template for python using sqlalchemy.
  * No sql injection possible. The sql is not executed directly but via the sqlalchemy engine.

## Usage

```bash
hamstersql -i example-sql -t templates/python -o output -v
```
This example is runnable and generate the python code from the example-sql folder. The python code is fully runnable.

## SQL file format

```
-- :name update_username :affected
-- :doc Update the username of all users
update users set username = :username
```

* You can add multiple queries to one file.
* Queries can be multiline.
* A query must start with `-- :name <name> :<type>` where `<name>` is the name of the function and `<type>` is the type of the function return.
  * `<type>` can be
    * `one` for a single row
    * `many` for multiple rows
    * `affected` for the number of affected rows
    * `inserted` for the inserted id
    * `scalar` for a single value
    * `none` for no return value (you can omit none if you want)
* You can add a comment to a query with `-- :doc <one line of comment text>`
* Each sql parameter starts with `:` like `:username` and will turn into a function parameter with the same name without the `:` (e.g. `username`).

## Template config (toml format)

* `staticFiles` - A list of files that are copied to the output directory.
* `groupFile` - The template for each sql input file (group).
* `groupFileName` - The parametrized name of the generated files for each group. `{{group}}` gets replaced by the group name. The group name gets derived from the basename of the sql files. So `samples.sql` will generate a group named `samples`.

## Generate a new template

* Copy the `templates/python` directory to a new directory.
* Adjust the `config.toml`.
* Adjust the group and the static files.
* Static files are optional.
* The group file is mandatory.

## Existing templates
* `templates/python` - Python with sqlalchemy