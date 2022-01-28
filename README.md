# Stream Wall

Music examples borrowed from https://ableandthewolf.com/ check them out and send sats.

## Help

Turn on the server locally.

```
-> % go run main.go
Starting server on 8080
2022/01/26 23:04:46 Serving songs on HTTP port: 8080
```

Create a session

```
curl http://localhost:8080/sessions -X POST
{"id":"64926bb7-9546-4ed0-a496-02407ff8e3cc","credits":0}
```

Open the stream corresponding to your session


Open `http://localhost:8080/sessions/64926bb7-9546-4ed0-a496-02407ff8e3cc/streams/MakinBeans/outputlist.m3u8` in VLC which will play the song.

## Installing Dependencies

```sh
go get -u github.com/go-chi/chi/v5
```

## Create files for HLS 

Install ffmpeg `sudo apt install ffmpeg`

```sh
ffmpeg -i MakinBeans.mp3 -c:a libmp3lame -b:a 128k -map 0:0 -f segment -segment_time 10 -segment_list outputlist.m3u8 -segment_format mpegts output%03d.ts
```