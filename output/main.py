import core
import sample


def main():
    session = core.make_session("sqlite:///:memory:")
    sample.create_users_table(session)
    print(sample.insert_user(session, "John", 25, "john@example.com"))
    print(sample.find_user(session, "John"))
    print(sample.list_users(session))
    print(sample.find_user_name_by_id(session, 1))
    print(sample.delete_users_with_age(session, 12))
    print(sample.delete_users_with_age(session, 25))


if __name__ == '__main__':
    main()
