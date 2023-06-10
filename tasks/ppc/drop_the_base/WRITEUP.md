# drop_the_base
Ну тут в принципе ничего сверхъестественного, нам приходит строчка такого формата:  
`[1] Base(12 + 46): 2cKCWUZorNMxwcfFLUmUyn2PkbCevRYV6pGi8oD3CpRne4N5Gw3DuniDKJ8wkkhQKQP8uMnbtQzTnF96TdiHovmv`  
Формат достаточно понятный после `Base` в скобках идет выражение `(12 + 46)`, посчитали, получили `58`. Долго думаем, что же делать с этой длинной строчкой и самое необычное что, наверное, может прийти в голову это попробовать задекодить её как `base58` строку. Декодим, получаем:  
`PxbFAqZWzMGgtFohGPmFsJXPaNcNVJsdPlltWoNgsoFHQedXzdkhpcfUHUyhkjDq`  
Отправляем, получаем в ответ новую строку.  
Методом ручного решения(или не ручного) понимаем, что тут всего 4 варианта base кодировки: **base64**, **base58**, **base32**, **base16**.  
Пишим скрипт, который всё это за нас будет делать:  
```python3 from pwn import *
import base64, base58

r = remote("185.104.115.19", 26065)


def parse_task(line):
    splitted = line.split(" ")
    enc = splitted.pop()
    splitted.pop(0)

    expr = "".join(splitted)[5:-2]
    return eval(expr), enc


try:
    while True:
        r.clean(0)
        resp = r.recvuntil(b"Decoded: ").decode("utf-8")
        print(resp)

        base, enc = parse_task(resp.split("\n")[0])

        dec = "blin"
        if base == 64:
            dec = base64.b64decode(enc)
        elif base == 58:
            dec = base58.b58decode(enc)
        elif base == 32:
            dec = base64.b32decode(enc)
        elif base == 16:
            dec = bytes.fromhex(enc)

        r.sendline(dec)
except Exception as e:
    print(e)
    r.interactive()
finally:
    r.close()
```
Скрипт проходит 1000 таких строк и получает ответ:  
```
NICE CO...khm.khm...BASE
surctf_drop_the_base_or_drop_the_bass
```  

`flag: surctf_drop_the_base_or_drop_the_bass`
