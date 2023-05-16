# SHMR (self-hosted multi runner)

한개의 머신에서 여러 러너를 구동할 수 있도록 도와줍니다.

## 다운로드

_도커와 curl이 미리 설치되어 있어야 합니다._

```sh
curl -L https://raw.githubusercontent.com/minpeter/SHMR/main/install.sh | sh
```

## self runner 추가 (add)

새로운 컨테이너에서 러너를 등록하고 구동시킴

```

shmr add -url [github repo url] -token [self action runner add token]

```

## 실행 중인 runner 확인 (list)

현재 해당 머신에서 실행 중인 러너 id 확인

```

shmr list

```

## 실행 중인 runner 삭제 (remove)

더 이상 사용하지 않을 runner 삭제

```

shmr remove -id [runner id] -token [runner remove token]

```

## 나를 위한 메모

새로운 릴리즈를 배포할 떄 git tags를 만들어야함

```sh
git commit ...
git tag -a 0.0.0-extra -m "someting"
git push --tags
```
