# 2048-AI

[![Build Status](https://travis-ci.org/xwjdsh/2048-ai.svg?branch=master)](https://travis-ci.org/xwjdsh/2048-ai)
[![Go Report Card](https://goreportcard.com/badge/github.com/xwjdsh/2048-ai)](https://goreportcard.com/report/github.com/xwjdsh/2048-ai)
[![](https://images.microbadger.com/badges/image/wendellsun/2048-ai.svg)](https://microbadger.com/images/wendellsun/2048-ai "Get your own image badge on microbadger.com")
[![DUB](https://img.shields.io/dub/l/vibe-d.svg)](https://github.com/xwjdsh/2048-ai/blob/master/LICENSE)

AI for the 2048 game, implements by expectimax search, powered by Go.

The web front-ends of 2048 game was forked from [gabrielecirulli/2048](https://github.com/gabrielecirulli/2048), respect and gratitude!

## Screenshot
![](https://raw.githubusercontent.com/xwjdsh/2048-ai/master/screenshot/2048-ai.jpg)

## How to run?
#### Via Go
```bash
go get github.com/xwjdsh/2048-ai
cd $GOPATH/src/github.com/xwjdsh/2048-ai
go build
./2048-ai
```

#### Via Docker
```shell
docker run -p 8080:8080 wendellsun/2048-ai
```

or build images by yourself:

```shell
# clone and enter project
docker build -t 2048-ai .
docker run -p 8080:8080 2048-ai
```

Then, you can access http://localhost:8080/ from the browser.

## Related

* [Expectimax search ](https://www.google.co.jp/url?sa=t&rct=j&q=&esrc=s&source=web&cd=3&cad=rja&uact=8&ved=0ahUKEwiVrsfmiojXAhWExbwKHa6GAuYQFgg3MAI&url=https%3A%2F%2Fweb.uvic.ca%2F~maryam%2FAISpring94%2FSlides%2F06_ExpectimaxSearch.pdf&usg=AOvVaw0pjG10MxUtkBvM-mvRNlew)
* [Model weight matrix](https://codemyroad.wordpress.com/2014/05/14/2048-ai-the-intelligent-bot/)
* Communication by websocket, implements by [gorilla/websocket](https://github.com/gorilla/websocket)

## Licence

[MIT License](https://github.com/xwjdsh/2048-ai/blob/master/LICENSE)
