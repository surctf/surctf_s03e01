# -*- coding: utf-8 -*-
from time import sleep

import serial


def get_number():
    ser = serial.Serial('/dev/ttyACM0', 9600)
    sleep(1.6)
    i = '\n'.encode()
    ser.write(i)
    value = ser.readline().decode().strip()
    print(value)
    return value


if __name__ == "__main__":
    numbers = [get_number() for i in range(20)]
    # temp = [b'4', b'\n', b'3', b'0', b'4', b'\xf0', b'1', b'\xf1', b'0', b'4', b'\xff', b'\xf1', b'\xfd', b'\xb1', b'1', b'1', b'\x8d', b'1', b'5', b'8']
    # print(b''.join(temp))
    print(numbers)