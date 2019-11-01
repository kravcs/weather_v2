# Get current weather data for several cities

JSON returns the data from cities specified in query parameter

# How to use it

Run in command line

```bash
$ mkdir myApp
$ cd myApp
$ git clone https://github.com/kravcs/U2VyZ2lpIEtyYXZjaGVua28gcmVjcnVpdG1lbnQgdGFzaw-.git .
$ docker-compose up
```

After application built, run the sample call

```bash
$ curl -X GET "http://localhost:8000/weather?city=Vienna,Moscow"
```