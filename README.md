[![Build Status](https://api.travis-ci.org/dazorni/9tac.svg?branch=master)](https://travis-ci.org/dazorni/9tac)
[![Coverage Status](https://coveralls.io/repos/github/dazorni/9tac/badge.svg?branch=master)](https://coveralls.io/github/dazorni/9tac?branch=master)

# 9tac

# Installation Guide with Docker

Pull the latest 9tac tag from the docker hub

```
docker pull dazorni/9tac
```

Pull the latest mongo database and set the same to mongo so we can link the environment later

```
docker pull mongo
```

Start the mongo database with a mounted folder

```
docker run --name mongo -v /home/data/db:/data/db -d mongo
```

Start the application and link to the mongo container

```
docker run -it -e "PORT=5000" -p 80:5000 --rm --link mongo:mongo dazorni/9tac
# or detached
docker run -e "PORT=5000" -p 80:5000 -d --link mongo:mongo dazorni/9tac
```

Now you can access the game on your docker ip at port ``80``
