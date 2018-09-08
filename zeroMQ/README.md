# zeroMQ Testing

Focusing on ZeroMQ 4.0.0 or later.

For running on windows, you may need MinGW - See http://zeromq.org/build:mingw

After building the libzmq.dll, you will need to specify the path to the header files
such as:

set CGO_CFLAGS=-g -O2 -I C:\Users\Chip\go\src\github.com\libzmq-4.2.5\include

In this case, the original source for libzmq was in C:\Users\Chip\go\src\github.com\libzmq-4.2.5
Another way is to copy the zmq header files to your  


Copy the libzmq.dll and libzmq.dll.a to  "\mingw\x86_64-w64-mingw32\lib"



Also directions at : https://github.com/pebbe/zmq4/issues/3
