#!/usr/bin/python

from threading import Thread, Lock

i = 0
lock = Lock()

def adder():
    global i
    for x in range(0, 1000000):
        lock.acquire()
        i += 1
        lock.release()

def subber():
    global i
    for x in range(0, 1000010):
        lock.acquire()
        i -= 1
        lock.release()

def main():
    adder_thr = Thread(target = adder)
    subber_thr = Thread(target = subber)
    adder_thr.start()
    subber_thr.start()
    adder_thr.join()
    subber_thr.join()
    print("Done: " + str(i))


main()
