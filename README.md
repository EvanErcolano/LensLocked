# Lens Locked

Photo gallery application written in go.


# What are the differences between?

go run

go build .

# Misc ports and mapping
Ports 0-1023 are well known ports used by System processes

- port 22 for ssh, secure login, sftp, and port forwarding
- port 80 for unencrypted traffic for HTTP protocol
- port 443 HTTP over TSL/SSL which is *HTTPS*

# What is a web request?
- *URL* = {protocol}  {domain}        {path}
    -   {https}     {google.com}    {/cats}
- *Headers*
    - metadata about the request that aren't rendered
        - user agent - browser automatically sends tells which browser and computer we're using
        - cookies 
- *Body*
    - data to be sent