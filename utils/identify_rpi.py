import requests, time
from uuid import getnode as get_mac
import socket

POST_URL = 'https://woc94o08ipz0.runscope.net'

for i in range(0,10):

	macraw = get_mac()

	mactem = '{:0>12x}'.format(macraw)
	macfmt = ':'.join([mactem[i:i+2] for i in range(0, 12, 2)])

	ipaddr = socket.gethostbyname(socket.gethostname())

	r = requests.post(POST_URL,
		data={"time":time.time(), "macaddr":macfmt, "ipaddr":ipaddr})

	time.sleep(10)

