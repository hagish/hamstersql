import core
import samples


def main():
    session = core.make_session("sqlite:///:memory:")
    samples.create_users_table(session)
    print("insert_user", samples.insert_user(session, "John", 25, "john@example.com"))
    print("find_user", samples.find_user(session, "John"))
    print("list_users", samples.list_users(session))
    print("find_user_name_by_id", samples.find_user_name_by_id(session, 1))
    print("delete_users_with_age", samples.delete_users_with_age(session, 12))
    print("delete_users_with_age", samples.delete_users_with_age(session, 25))


if __name__ == '__main__':
    main()
