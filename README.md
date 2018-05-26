# sdbx-discord-dj-bot
샌드박스 디스코드 서버에 쓰이는 음악봇입니다.

# 설치
적절한 버전의 고 환경이 갖추어져 있다면, 아래의 명령어로 DJ봇을 설치 가능합니다.
```
go get github.com/sunho/sdbx-discord-dj-bot/djbot
```
설치를 마친 후에는 다음 명령어를 실행하여 `tokens.txt`을 생성하십시오.
```
djbot -initial
```
`tokens.txt`안의 `discord_token`을 디스코드 봇 토큰으로 `youtube_api_key`를 유튜브 api 키로  `bot_owner_id`를 당신의 [디스코드 id](https://support.discordapp.com/hc/en-us/articles/206346498-Where-can-I-find-my-User-Server-Message-ID-)로 대체하십시오.
그 뒤 아래의 명령어를 실행하면 DJ봇이 구동됩니다.
```
djbot
```
djbot이 음악을 재생하기 위해서는 ffmpeg와 [dca](https://github.com/bwmarrin/dca) 그리고 youtube-dl이 PATH안에 들어있어야 합니다. dca같은 경우는
```
go get github.com/bwmarrin/dca/cmd/dca
```
로 설치할 수 있고 나머지는 인터넷에 찾아보시면 쉽게 설치하실 수 있습니다.

# 명령어
모든 명령어 앞에는 `!!`을 붙이셔야 합니다.
!!help로 명령어의 설명을 보실 수 있습니다.

# 설정변수
| 설정변수 | 설명 |
| :-- | --: |
| CertainChannelInputOnly | 특정 텍스트 채널에서만 명령어를 입력받도록 합니다 |
| CertainChannel | 위의 특정 텍스트 채널은 이 변수에 그 채널의 디스코드 id를 넣어 설정할 수 있습니다 |
| ConnectToCertainVoiceChannelOnly | 특정 보이스 채널에만 접속할 수 있도록 합니다 |
| VoiceChannel | 위의 특정 보이스 채널은 이 변수에 그 채널의 디스코드 id를 넣어 설정할 수 있습니다 |
| VoteToSkip | 이것이 false 이면 s명령어가 투표를 진행하지 않고 바로 곡을 넘깁니다 |
