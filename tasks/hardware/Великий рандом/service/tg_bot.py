from telebot.apihelper import ApiTelegramException

from credentials import TG_TOKEN
import telebot
import db

bot = telebot.TeleBot(TG_TOKEN)


@bot.message_handler(commands=['start'])
def welcome(message):
    bot.reply_to(message, "Отправьте мне свой логин")
    print(message)


@bot.message_handler()
def login_handler(message):
    if not db.get_user_by_tg(message.from_user.id):
        if not db.get_user_by_login(message.text):
            db.create_user(message.text, message.from_user.id)
            bot.reply_to(message, "Вы успешно зарегистрировались!\nТеперь можете войти на сайт с кодом")
        else:
            bot.reply_to(message, "Такой логин уже существует!")
    else:
        bot.reply_to(message, "Вы уже зарегистрировались!")


def send_code(code, login):
    user = db.get_user_by_login(login)
    if user:
        try:
            bot.send_message(user.tg_user_id, f"Ваш код для авторизации: {code}\nНикому не сообщайте его!")
        except ApiTelegramException:
            return False
        return True
    else:
        return False


if __name__ == "__main__":
    bot.infinity_polling()
