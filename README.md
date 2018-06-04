# sdbx-discord-dj-bot
샌드박스 디스코드 서버에 쓰이는 음악봇입니다.

# 요구사항

디제이봇은 [도커](https://www.docker.com/community-edition) 와 [도커컴포즈](https://docs.docker.com/compose/install/#prerequisites) 를 요구합니다. 이를 설치하는 방법은 해당 링크를 클릭하여 찾으실 수 있습니다.

# 설치
먼저 깃 저장소를 복제합니다.
```
git clone https://github.com/sunho/sdbx-discord-dj-bot
cd sdbx-discord-dj-bot
```
그후 config.yaml을 [이곳](djbot/config.go)을 참고하여 작성한 후 아래 커맨드를 실행해주세요.
```
docker-compose up
``` 

