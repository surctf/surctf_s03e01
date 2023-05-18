from telebot.apihelper import ApiTelegramException
from pywebio import start_server
from pywebio.output import popup, put_text
from pywebio.input import input, actions
import serial_bridge as sb
import db
from credentials import FLAG
from tg_bot import send_code


def signin():
    while True:
        put_text("Регистрация здесь: https://t.me/surctf_secure_system_bot")
        login = input('Введите логин:')
        user = db.get_user_by_login(login)
        if user:
            break
        else:
            popup("Такого пользователя не существует!")

    while True:
        right_code = sb.get_number()
        is_sent = send_code(right_code, login)
        if not is_sent:
            popup("Не удалось отправить код!")
            break
        else:
            code = input("Введите код безопасности: ")
            if code == right_code:
                if user.is_admin:
                    put_text(FLAG)
                    break
                else:
                    put_text("У вас ничего нет!")
                    break
            else:
                popup("Неверный код!")


def main():
    # action = actions("Добро пожаловать!", [
    #     {'label': 'Войти', 'value': "signin"},
    #     {'label': 'Зарегистрироваться', 'value': "signup"}
    # ])
    # if action == 'signin':
    signin()


if __name__ == '__main__':
    start_server(main, port=8080, debug=False)
