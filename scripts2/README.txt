IMPORTANT!

Files 'ca.cert' and 'ca.key' are required to sign the certificate in the process of creation.
They are not required to run the service as far as 'service.pem' and 'service.key' are good.

-----

Before cert/key creation you must edit 'certificate.conf' and set correct hostname(s) and IP(s).
File 'certificate.conf' is supposed to be used to create cert and key for docker image if it includes
  ...
  CN = s70xx
  ...
  DNS.1 = s70xx
  DNS.2 = localhost
  IP.1 = 172.16.70.xx
  ...

It specifies 'CN', 'DNS.1' (hostname), e.g. 's7026', and IP, e.g. '172.25.70.26'.
The created 'service.pem' and 'service.key' are supposed to be used with the docker image!

If 'CN' and 'DNS.1' are set to 'localhost' or hostname/IP of your computer, the 'certificate.conf' is supposed to create cert/key
for normal use [with your computer as a host]. Note that these 'service.pem' and 'service.key' won't work with Docker image.
