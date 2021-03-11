from machine import Pin
from neopixel import NeoPixel

# brightness = 0.01
# color = lambda x: tuple(map(lambda i: int(i * brightness), x))

pin = Pin(0, Pin.OUT)   # set GPIO0 to output to drive NeoPixels

import esp
grb_buf = bytearray([100, 0, 100, 0, 100 ,0])
esp.neopixel_write(pin, grb_buf, True)


import uarray

class Led:
    """
    Led
    used to control 
    """
    def __init__(self, pin, size):
        self.data = uarray.array()