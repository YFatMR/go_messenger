import random
import requests

HOST = "http://localhost:8010"


class UserData:
    def __init__(self, id: int, token: str) -> None:
        self.id = id
        self.token = token


def create_user(i: int) -> int:
    names = [
        "Максим",
        "Dima",
        "Denis",
        "Егор",
        "Макар",
    ]

    surnames = [
        "Лебедев",
        "Смирнов",
        "Федоров",
        "Зубков",
        "Emshanov",
    ]

    data = {
        "userData": {
            "nickname": f"{i}",
            "name": f"{names[random.randrange(len(names))]}",
            "surname": f"{surnames[random.randrange(len(surnames))]}",
        },
        "credential": {
            "email": f"{i}",
            "password": f"{i}",
        }
    }
    response = requests.post(f'{HOST}/v1/users', json=data)
    assert(200 == response.status_code)
    return int(response.json()['ID'])


def create_user_with_token(i: int) -> UserData:
    user_id = create_user(i)

    data = {
        "email": f"{i}",
        "password": f"{i}",
    }
    response = requests.post(f'{HOST}/v1/token', json=data)
    assert(200 == response.status_code)
    return UserData(id=user_id, token=response.json()['accessToken'])


def create_dialog(data1: UserData, data2: UserData) -> int:
    headers = {"Content-Type": "application/json; charset=utf-8", "Authorization": data1.token}

    response = requests.post(f'{HOST}/v1/dialogs?userID={data2.id}', headers=headers)
    assert(200 == response.status_code)
    return int(response.json()['dialogID']['ID'])


def create_message(user_data: UserData, dialog_id: int, text: str) -> int:
    headers = {"Content-Type": "application/json; charset=utf-8", "Authorization": user_data.token}

    response = requests.post(f'{HOST}/v1/dialogs/{dialog_id}/messages', headers=headers, json={
        'text': text,
    })
    assert(200 == response.status_code)


def create_dialog_with_messages(user_data1: UserData, user_data2: UserData, messages_cnt: int):
    d_id = create_dialog(user_data1, user_data2)
    msg = [
        u"Hello!",
        u"Очень рад тебя слышать, потому что мы давно не переписывались. Как сам?",
        u"All good!",
        u"Рад слышать",
        u"Hi!",
        u"Собираюсь завтра к бабушке в деревню, тебе что-нибудт привезти. Возможно салатик, картошку и т д?",
    ]

    for i in range(messages_cnt):
        d = user_data1
        if random.randrange(100) % 2 == 0:
            d = user_data2
        create_message(d, d_id, msg[random.randrange(len(msg))] + f" {i}")


lst = [create_user_with_token(i) for i in range(60)]

for i in range(1, len(lst)):
    cnt = 10
    if i < 5:
        cnt = 60
    create_dialog_with_messages(lst[0], lst[i], cnt)
print("OK!")
