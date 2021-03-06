# Presence [![GoDoc](https://godoc.org/github.com/cihangir/presence?status.svg)](https://godoc.org/github.com/cihangir/presence) [![Build Status](https://travis-ci.org/cihangir/presence.svg?branch=master)](https://travis-ci.org/cihangir/presence)

Presence in Redis
------------------

One of the most important and resource intensive operation in chat systems is handling user presence.
Because in order to start the chat one should provide the presence information.
And keeping this data up-to-date is another challenge.

There are 3 main parts for handling user presence,
* Notification
* Persistency
* Data source

#### Starting with number one:
If we calculate the message count with a standart architecture for sending notifications to the online users:
count = (online user count) * (channel participant count) * (events), since every part is a
multiplication even increasing by one will cause huge effects.

#### Number two:
Transient presence data should be commited and updated regularly, after an inactive duration presence data should be set to offline
either by the application or by the persistency layer

#### Number three:
Online clients should notify the presence system with their unique identifier repeatedly. We can call it as heartbeat.


"This package aims to handle the persistency layer for presence"

For usage see examples below or go to godoc page.

## Install and Usage

Install the package with:

```bash
go get github.com/cihangir/presence
```

Import it with:

```go
import "github.com/cihangir/presence"
```


## Examples

#### Initialization of a new Redisence

```go

// create a presence system
backend, err := NewRedis(serverAddr, dbNumber, timeoutDuration)
if err != nil {
    return err
}

s, err := New(backend)
if err != nil {
    return err
}

```

#### Basic Operations

```go

// send online presence data to system - user log in
err = s.Online("id")
err = s.Online("id2")


// send offline presence data to system - user log out
err = s.Offline("id")
err = s.Offline("id2")

// get status of some ids
status, err := s.Status([]string{"id20", "id21"}...)
if err != nil {
    return err
}

for _, st := range status {
    if st.Status != Offline {
        //....
    }
}

```

#### Listening for events

```go

go func() {
    s.Online("id")
    time.Sleep(time.Second * 1)
    s.Online("id")
    s.Online("id2")
}()

// start listening to them
for event := range s.ListenStatusChanges() {
    switch event.Status {
    case Online:
        // ....
    case Offline:
        // ....
    }
}

```


# Redis configuration
To get the events from the redis database we should uptade the redis config with the following data

`redis-cli config set notify-keyspace-events Ex$`

Or
set in redis.conf
`notify-keyspace-events "Ex$"`

for more info http://redis.io/topics/notifications

## License

The MIT License (MIT) - see LICENSE for more details
