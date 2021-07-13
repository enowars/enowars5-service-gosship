# Service Documentation
```
              _____ _____ _    _ _
             / ____/ ____| |  | (_)
   __ _  ___| (___| (___ | |__| |_ _ __
  / _` |/ _ \\___ \\___ \|  __  | | '_ \
 | (_| | (_) |___) |___) | |  | | | |_) |
  \__, |\___/_____/_____/|_|  |_|_| .__/
   __/ |                          | |
  |___/                           |_|

```

- [Introduction](#introduction)
- [Vulnerabilities and Exploits](#vulnerabilities-and-exploits)
    * [Vulnerability 1 (private rooms)](#vulnerability-1--private-rooms-)
        + [Exploit](#exploit)
    * [Vulnerability 2 (direct messages)](#vulnerability-2--direct-messages-)
        + [Exploit](#exploit-1)
- [Lessons Learned](#lessons-learned)

# Introduction
goSSHip is a SSH chat service written in Go and inspired by the blog post ["Why aren't we using SSH for everything?" by Andrey Petrov](https://shazow.net/posts/ssh-how-does-it-even/).
The service allows users to log in with their default SSH client and chat with other people connected to the service. Additionally, users can call special commands to create private (password protected) or public rooms or send direct messages. The following commands are available in goSSHip:

```
+---------------------------+-----------+------------------------------------------+
| COMMAND                   | ALIASES   | HELP                                     |
+---------------------------+-----------+------------------------------------------+
| /dm [user] [msg]          |           | send a direct message to a user          |
+---------------------------+-----------+------------------------------------------+
| /help                     | /h, /?    | show the help for all available commands |
+---------------------------+-----------+------------------------------------------+
| /exit                     | /quit, /q | leave the chat                           |
+---------------------------+-----------+------------------------------------------+
| /info                     | /i        | info about the logged-in user            |
+---------------------------+-----------+------------------------------------------+
| /reply [msg]              | /r        | reply to your last direct message        |
+---------------------------+-----------+------------------------------------------+
| /history [user]           |           | show the direct message history          |
+---------------------------+-----------+------------------------------------------+
| /shrug                    |           | ¯\_(ツ)_/¯                               |
+---------------------------+-----------+------------------------------------------+
| /rename [new name]        |           | change your username                     |
+---------------------------+-----------+------------------------------------------+
| /create [room] <password> |           | create a new room                        |
+---------------------------+-----------+------------------------------------------+
| /join <room> <password>   | /j        | join a room                              |
+---------------------------+-----------+------------------------------------------+
| /users                    |           | list users on the server                 |
+---------------------------+-----------+------------------------------------------+
| /rooms                    |           | list rooms on the server                 |
+---------------------------+-----------+------------------------------------------+
```

# Vulnerabilities and Exploits
## Vulnerability 1 (private rooms)
### Exploit

## Vulnerability 2 (direct messages)
### Exploit

# Lessons Learned
