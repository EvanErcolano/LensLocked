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

## Chapter 11 - How do we remember users?
    - servers don't "remember" what you did 15 minutes ago
    - servers are stateless they handle each request independently
    - Stateless servers make it easier - we don't rely on local info
        - everything we need is stored within the request itself 
    
    - So how do we remember who a user is ?
        - We let the user tell us who they are every web request!
            - What if the user lies to us?
            - Users don't tell us
    Cookies solve this problem
    
    ## What are cookies?
    - Data stored on the user's computer
    - Used for authenticating web sessions
    - Usually have session id info in the cookie
    - Cookies are linked to a specific domain (facebook.com)
    - Your browser will take care of sending cookies linked to the website with each request
    - The server can edit the cookies and the user can also edit the cookies (foreshadows tampering with cookies)



