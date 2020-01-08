import httplib
import time

host = "localhost"
port = 4481
cons = []

def req(con, path):
  con.request("GET", path)
  result = con.getresponse()
  result.read()

# connection 1 - initial request. establishes frontend/backend connections
con1 = httplib.HTTPConnection(host, port)
req(con1, "/fecon1")

# connection 1 - second request. reuses connection. after first reuse in
# `aggressive` mode, the connection can be reused by other frontend connections
req(con1, "/fecon1")

# connection 2 - initial request. establishes frontend connection. reuses
# backend connection
con2 = httplib.HTTPConnection(host, port)
req(con2, "/fecon2")

# connection 1 - close frontend connection. backend connection remains open
con1.close()

# connection 2 - second request. reuses same backend connection
req(con2, "/fecon2")
