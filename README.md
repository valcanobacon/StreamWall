# Stream Wall

Music examples borrowed from https://ableandthewolf.com/ check them out and send sats.

## Help

Turn on the server locally.

```
-> % go run main.go
Starting server on 8080
2022/01/26 23:04:46 Serving songs on HTTP port: 8080
2022/01/26 23:05:14 /MakinBeans/outputlist.m3u8
2022/01/26 23:05:14 /MakinBeans/output000.ts
2022/01/26 23:05:14 /MakinBeans/output001.ts
2022/01/26 23:05:14 /MakinBeans/output002.ts
2022/01/26 23:05:14 /MakinBeans/output003.ts
2022/01/26 23:05:22 /MakinBeans/output004.ts
2022/01/26 23:05:32 /MakinBeans/output005.ts
2022/01/26 23:05:42 /MakinBeans/output006.ts
2022/01/26 23:05:52 /MakinBeans/output007.ts
2022/01/26 23:06:02 /MakinBeans/output008.ts
2022/01/26 23:06:12 /MakinBeans/output009.ts
2022/01/26 23:06:22 /MakinBeans/output010.ts
2022/01/26 23:06:32 /MakinBeans/output011.ts
2022/01/26 23:06:42 /MakinBeans/output012.ts
2022/01/26 23:06:52 /MakinBeans/output013.ts
2022/01/26 23:07:02 /MakinBeans/output014.ts
2022/01/26 23:07:12 /MakinBeans/output015.ts
2022/01/26 23:07:22 /MakinBeans/output016.ts
2022/01/26 23:07:32 /MakinBeans/output017.ts
2022/01/26 23:07:42 /MakinBeans/output018.ts
```

Open `http://localhost:8080/MakinBeans/outputlist.m3u8` in VLC which will play the song.

## Installing Dependencies

```sh
go get -u github.com/go-chi/chi/v5
```

## Create files for HLS 

Install ffmpeg `sudo apt install ffmpeg`

```sh
ffmpeg -i MakinBeans.mp3 -c:a libmp3lame -b:a 128k -map 0:0 -f segment -segment_time 10 -segment_list outputlist.m3u8 -segment_format mpegts output%03d.ts
```