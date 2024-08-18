# Restaurant Reservation System
This is a restaurant reservation system in golang. This system provides APIs to deal with all things related to reservations.

## APIs available
This system has primarily 7 APIs to work with:

    1. / : Default API just to make sure server is up and running
    2. /add : Add a new reservation
    3. /view : View reservations by day
    4. /confirm: Confirm a reservation
    5. /cancel: Cancel a reservation
    6. /waitinglist : Dynamic waitlist management visualization through API
    7. /availability : View all available slots for reservations by day

Pseudo notification to waitlist members in order of requests and waitlist managemenet is inherently handled. These APIs, by design, handle multiple locations.

## Code Organization

    -- restaurant_reservation
        |-- api
            |-- api_handler.go
            |-- api_helper.go
            |-- auto_cancelation.go
            |-- waitlist_management.go
            |-- api_test.go
            |-- endpoint_test.go
        |-- utilities
            |- utility.go
        |-- reservation_system.go
        |-- go.mod
        |-- Readme.md
        |-- build
        |-- clean
        |-- test

## Getting Started
First thing we need to make sure is that the go version being used is 1.22.x or above. This system is built on 1.22.0 and it may or may not work on earlier versions.
Using a git clone command should allow you to clone this repo

    git clone <url>

After cloning the repo, change the diretory to restaurant_reservation. Then you can use the build script and build the code.

    ./build

At this moment, it should build and work seamlessly. If any of the tests fail, the build would fail. You can alter that behavior in the script but I ***highly recommend NOT to do so.***

## Test, Build and Clean (Testing and Coverage)

To build and run it in a simple way, use

    ./build
    ./reservation_system

Or

    go build -race -o reservation_system
    ./reservation_system

If you're planning to run it on OS other than MAC, feel free to use:

    GOOS=linux GOARCH=arm64 go build -race -o reservation_system

or

    GOOS=windows GOARCH=arm64 go build -race -o reservation_system

NOTE: GOOS and GOARCH selection might differ on your machine.

To directly run, use

    go run reservation_system.go

There are three bash files in the repo. It is suggested that you use "build" bash file, which would run all the tests first and then gives you a binary which you can run. You can always use 

    ./clean

to clean any generated files or binaries.

NOTE:- please make sure that you run chmod command to make all the bash files executable, if they're not already. You can do so as:

    chmod +x test

OR

    chmod 777 test

Similar statments for build and clean.

## Usage

### / default API

Go to a browser or Postman and run a simple GET request with localhost:7070/ or 127.0.0.1:7070/ and that should serve you the default message

    {"message":"The server is up. Use other endpoints for reservation"}

#### All other APIs could be documented but this is not the place for them to go. There should be a swagger doc soon to address those.

## Prerequisites

There are no special prerequisites for this code to run. Just need to include this under your src folder in $GOPATH, and it should work smoothly.

## Scripts Heirarchy

The 'build' script calls 'test' script by default as we want to implement TDD as our standard practice. Now, the 'test' script call the 'clean' script within itself to clean all / any temporary files generated during testing. Make sure if you add your test cases and those generate some temporary files, you add them in 'clean' script to get them removed. Also, you can use both 'test' and 'clean' script individually according to your needs.

NOTE :- The code has ample comments to make life easy.

