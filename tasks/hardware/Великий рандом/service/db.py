import peewee as pw

database = pw.SqliteDatabase("gpsch.db")


class Users(pw.Model):
    login = pw.CharField()
    tg_user_id = pw.IntegerField()
    is_admin = pw.BooleanField()

    class Meta:
        database = database


def get_user_by_login(login):
    try:
        user = Users.get(Users.login == login)
    except pw.DoesNotExist:
        return None
    return user


def get_user_by_tg(tg_user_id):
    try:
        user = Users.get(Users.tg_user_id == tg_user_id)
    except pw.DoesNotExist:
        return None
    return user


def create_user(login, tg_user_id):
    if not get_user_by_login(login):
        user = Users.create(login=login, tg_user_id=tg_user_id, is_admin=False)
        database.commit()
        return user
    else:
        return None


def create_tables():
    with database:
        database.create_tables([Users])


if __name__ == "__main__":
    create_tables()