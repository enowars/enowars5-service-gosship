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
    * [Database](#database)
- [Vulnerabilities](#vulnerabilities)
    * [Vulnerability 1 (private rooms)](#vulnerability-1-private-rooms)
        + [Exploit](#exploit)
        + [Fix](#fix)
    * [Vulnerability 2 (direct messages)](#vulnerability-2-direct-messages)
        + [Exploit](#exploit-1)
        + [Fix](#fix-1)

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

Additionally, the service provides a gRPC admin interface to send messages to a specific room and fetch all users' direct messages. This interface should only be accessible by the checker (the public key of the checker is embedded in the service), so it can verify the persistence of flags. An example of the admin interface can be found here: [checker/cmd/checker-test/main.go](../checker/cmd/checker-test/main.go).

## Database
The embedded key-value database [badger](https://github.com/dgraph-io/badger) is used to persist users, messages, rooms, and the SSH private key of the server. All database entries (except the config) expire after a certain time (messages: 30min, users and rooms: 1h; [service/pkg/database/database.go](../service/pkg/database/database.go#L103)).

# Vulnerabilities
The service has two different flag stores and one vulnerability each. The fixed version of the service can be found on the [fixed](https://github.com/enowars/enowars5-service-gosship/compare/fixed) branch.

## Vulnerability 1 (private rooms)
The first flag store is in the messages of password-protected rooms, and the vulnerability linked to this flag store is that users can join private rooms without knowing the correct password.

### Exploit
Let's assume that the flag is currently stored in a password-protected room called `private`, and the attacker does not know the password to join the room. To exploit this vulnerability, the attacker needs to create a new room with the same (case-insensitive) name (e.g., `Private`). When creating a new room, the creator will automatically join the room. Hence, this will also update the current room of the user in the database. Updating the current room contains a bug that saves the lowercase room name in the database ([service/pkg/chat/user.go](../service/pkg/chat/user.go#L120)). So if the attacker leaves the service and rejoins, they are automatically in the lowercase name of the created room (`private`)  and able to read the previous messages in that room.
A proof-of-concept exploit script for this vulnerability can be found in the checker folder: [checker/cmd/private-room-exploit/main.go](../checker/cmd/private-room-exploit/main.go)

### Fix
To fix this vulnerability the `strings.ToLower` function call needs to be removed in [service/pkg/chat/user.go](../service/pkg/chat/user.go#L120).
```diff
 func (u *User) UpdateCurrentRoom(room string) error {
-       u.CurrentRoom = strings.ToLower(room)
+       u.CurrentRoom = room
        return u.DBUpdate()
 }
````

## Vulnerability 2 (direct messages)
The second flag store is in the direct messages between two users, and the connected vulnerability lies in creating the admin session tokens.

### Exploit
To access the flags in the direct messages the `DumpDirectMessages` method from the admin interface is required. That means the attacker needs to get access to the admin interface first.
In the `GenerateRandomSessionToken` function ([service/pkg/rpc/admin/auth/auth.go](../service/pkg/rpc/admin/auth/auth.go#L36)) a supposedly random 32-byte session token is generated. A race condition in the function when calling the set function will lead to a 32-byte session token where every byte has the same value ([this is a common go mistake](https://github.com/golang/go/wiki/CommonMistakes#using-reference-to-loop-iterator-variable)). Thus, only 256 different session tokens are possible. The attacker easily finds a valid token using brute-forced.
With access to the admin interface, the attacker invokes the `DumpDirectMessages` method and retrieves the flag.
A proof-of-concept exploit script for this vulnerability can be found in the checker folder: [checker/cmd/session-exploit/main.go](../checker/cmd/session-exploit/main.go)

### Fix
To fix this vulnerability the `GenerateRandomSessionToken` function ([service/pkg/rpc/admin/auth/auth.go](../service/pkg/rpc/admin/auth/auth.go#L36)) needs to generate truly random session tokens.
```diff
 func GenerateRandomSessionToken() string {
        token := make([]byte, sessionTokenSize)
-       var wg sync.WaitGroup
-       wg.Add(sessionTokenSize)
-       for i := 0; i < sessionTokenSize; i++ {
-               go set(&wg, &token[i], &i)
-       }
-       wg.Wait()
+       _, _ = rand.Read(token)
        tokenHash := sha256.New()
        tokenHash.Write(token)
```
