CC = gcc
CFLAGS = -Wall -Wextra -std=c99 -I/usr/include/dbus-1.0 -I/usr/lib/x86_64-linux-gnu/dbus-1.0/include
LDFLAGS = $(shell pkg-config --libs dbus-1)
TARGET = dbus_example

.PHONY: all clean

all: $(TARGET)

$(TARGET): main.c
	$(CC) $(CFLAGS) -o $@ $< $(LDFLAGS)

clean:
	rm -f $(TARGET)