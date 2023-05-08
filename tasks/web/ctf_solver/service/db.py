from peewee import *
from dotenv import load_dotenv
import os, random

load_dotenv()

db = PostgresqlDatabase(os.getenv("DB"), user=os.getenv("DB_USER"), password=os.getenv("DB_PASS"),
                        host=os.getenv("DB_HOST"), port=os.getenv("DB_PORT"))


class Task(Model):
    name = CharField()
    description = CharField()
    flag = CharField()
    solved = BooleanField()

    class Meta:
        database = db

def add_task(task):
    return random.randint(10000000, 100000000000000000), "%dд" % random.randint(1, 7)

db.connect()
db.create_tables([Task])