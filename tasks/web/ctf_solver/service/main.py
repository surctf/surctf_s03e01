from flask import Flask, request, send_file
from dotenv import load_dotenv
import os, hashlib
from db import Task, add_task
from functools import wraps

load_dotenv()

app = Flask(__name__)
API_KEY_HASH = hashlib.md5(os.getenv("API_KEY").encode("utf-8")).hexdigest()


def api_key_required(f):
    err_resp = {"error": "неправильный 'api_key'", "result": None}

    @wraps(f)
    def wrapper(*args, **kwargs):
        api_key = request.args.get("api_key")
        if api_key is None:
            return err_resp

        try:
            h = hashlib.md5(api_key.encode("utf-8")).hexdigest()
            if h != API_KEY_HASH:
                return err_resp
        except Exception:
            return err_resp

        return f(*args, **kwargs)
    return wrapper


@app.route("/solve", methods=["POST"])
@api_key_required
def solve():
    resp = {"result": None, "error": None}

    body = request.json
    if body is None or any([field not in body for field in ["name", "description"]]):
        resp["error"] = "Кажется ты не так метод используешь, вызови /help чтобы понять что-как"
    else:
        task_id, wait_time = add_task(body)
        resp["result"] = "Таск добавлен в очередь на решение(id=%s), примерное время решения - %s" % (task_id, wait_time)

    return resp


@app.route("/solved_list", methods=["GET"])
@api_key_required
def solved_list():
    resp = {"result": [], "error": None}

    tasks = Task.select().where(Task.solved == True)
    for task in tasks:
        resp["result"].append({"id": task.get_id(), "name": task.name, "description": task.description, "flag": task.flag})

    return resp


@app.route("/help", methods=["GET"])
def _help():
    resp = {"result": None, "error": None}

    method = request.args.get("method")

    path_to_help = "./help/%s"
    if method is None:
        path_to_help = path_to_help % "help"
    else:
        path_to_help = path_to_help % method

    try:
        with open(path_to_help, "r") as f:
            result = f.read()
        resp["result"] = result
    except Exception as e:
        print(e)
        resp["error"] = "Кажется такого метода нет, вызови /help, чтобы получить список методов"

    return resp

@app.errorhandler(404)
def not_found(e):
    return _help()
