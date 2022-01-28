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

### Dev Session

The Nil UUID `00000000-0000-0000-0000-000000000000` is pre created with credits.

http://localhost:8080/sessions/00000000-0000-0000-0000-000000000000/streams/MakinBeans/outputlist.m3u8

## Installing Dependencies

```sh
go get -u github.com/go-chi/chi/v5
```

## Create files for HLS 

Install ffmpeg `sudo apt install ffmpeg`

```sh
ffmpeg -i MakinBeans.mp3 -c:a libmp3lame -b:a 128k -map 0:0 -f segment -segment_time 10 -segment_list outputlist.m3u8 -segment_format mpegts output%03d.ts
```

# Generating Lightning GRPC fileso

```
curl -o src/lightning.proto -s https://raw.githubusercontent.com/lightningnetwork/lnd/master/lnrpc/lightning.proto
```