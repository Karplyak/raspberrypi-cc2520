import requests, time
from uuid import getnode as get_mac
import socket

POST_URL = 'https://rpi-disc-gbns1w57rtvk.runscope.net'

for i in range(0,10):

	macraw = get_mac()

	mactem = '{:0>12x}'.format(macraw)
	macfmt = ':'.join([mactem[i:i+2] for i in range(0, 12, 2)])

	s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
	s.connect(('www.google.com', 0))
	ipaddr = s.getsockname()[0]

	r = requests.post(POST_URL,
		data={"time":time.time(), "macaddr":macfmt, "ipaddr":ipaddr})

	if r.status_code == 200:
		break

	time.sleep(10)

