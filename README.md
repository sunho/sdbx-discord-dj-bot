# sdbx-discord-dj-bot
샌드박스 디스코드 서버에 쓰이는 음악봇입니다.

# 설치
적절한 버전의 고 환경이 갖추어져 있다면, 아래의 명령어로 DJ봇을 설치 가능합니다.
```
go get github.com/ksunhokim123/sdbx-discord-dj-bot/cmd
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
# 명령어
모든 명령어 앞에는 `!!`을 붙이셔야 합니다.

| 명령어 | 설명 |
| :-- | --: |
| c | DJ봇을 보이스 채널에 접속시킵니다. |
| admin chid | 채널을 선택해 해당 채널의 디스코드 id를 알아냅니다. |
| admin disconnect | DJ봇을 보이스 채널에서 나가게합니다. |
| admin envget | 모든 설정변수의 이름과 값을 출력합니다. |
| admin envset | 설정 변수의 값을 변경시킵니다(변수이름,값) |
| admin fskip | 현재 재생중인 음악을 강제로 넘깁니다. |
| p | 유튜브 링크로부터 음악을 추가합니다(링크) |
| remove | 음악큐의 음악을 삭제합니다 (인덱스) |
| rremove | 음악큐의 음악들을 삭제합니다 (시작인덱스~끝인덱스) |
| s | 간단한 투표를 진행해 현재 음악을 넘길 지 말지를 결정합니다. |
| sr | 유튜브에서 음악을 검색해 추가합니다 (키워드) |
| start | 음악 재생을 시작합니다 |

# 설정변수
| 설정변수 | 설명 |
| :-- | --: |
| CertainChannelInputOnly | 특정 텍스트 채널에서만 명령어를 입력받도록 합니다 |
| CertainChannel | 위의 특정 텍스트 채널은 이 변수에 그 채널의 디스코드 id를 넣어 설정할 수 있습니다 |
| ConnectToCertainVoiceChannelOnly | 특정 보이스 채널에만 접속할 수 있도록 합니다 |
| VoiceChannel | 위의 특정 보이스 채널은 이 변수에 그 채널의 디스코드 id를 넣어 설정할 수 있습니다 |
| VoteToSkip | 이것이 false 이면 s명령어가 투표를 진행하지 않고 바로 곡을 넘깁니다 |
