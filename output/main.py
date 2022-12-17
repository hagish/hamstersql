import core
import samples


def main():
    session = core.make_session("sqlite:///:memory:")
    samples.create_users_table(session)
    print(samples.insert_user(session, "John", 25, "john@example.com"))
    print(samples.find_user(session, "John"))
    print(samples.list_users(session))
    print(samples.find_user_name_by_id(session, 1))
    print(samples.delete_users_with_age(session, 12))
    print(samples.delete_users_with_age(session, 25))


if __name__ == '__main__':
    main()
