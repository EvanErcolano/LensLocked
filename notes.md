
# Misc ports and mapping
Ports 0-1023 are well known ports used by System processes

- port 22 for ssh, secure login, sftp, and port forwarding
- port 80 for unencrypted traffic for HTTP protocol
- port 443 HTTP over TSL/SSL which is *HTTPS*

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

### What are cookies?
    - Data stored on the user's computer
    - Used for authenticating web sessions
    - Usually have session id info in the cookie
    - Cookies are linked to a specific domain (facebook.com)
    - Your browser will take care of sending cookies linked to the website with each request
    - The server can edit the cookies and the user can also edit the cookies (foreshadows tampering with cookies)

### How do we secure our cookies?
    1. Cookie tampering
        - Editing cookies
        - Preventable by Obfuscating the value, signing it then hashing it, sessions ...
        - We will generate a gibberish remember token and store it in the db
        - we will create one and associate with a user and that will be sent back and forth
        - we can then lookup the user by that value to figure out their session / who it is
        - we don't store the sessionid we store the hash of it, jsut like the passwords
            - this prevents issues where if our db leaks attackers could fake being users via the cookie data stored
    2. A database leak that allows users to create fake cookies
        - Use db data to create cookies
    3. Cross site scripting (XSS)
        - Letting users inject JS into your site
    4. Cookie theft (via packet sniffing or physical access)
        - Stealing cookies and pretending they are yours
        - Preventable by having SSL
    5. Cross site request forgery (CSRF)
        - Sending web requests to other servers on behalf of a user w/out them knowing

4 and 5 will be covered when we prepare to deploy to production as they are solved via SSL
