package de

import (
	"context"
	"errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func set() (*client.Client, error) {
	// 도커 클라이언트 초기화
	// ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, errors.New("failed to initialize Docker client")
	}

	// 도커 사용 가능 여부 확인
	_, err = cli.Ping(context.Background())
	if err != nil {
		return nil, errors.New("failed to connect to Docker daemon")
	}

	// "minpeter/gh-action-runner:latest" 이미지의 존재 여부 확인
	imageName := "minpeter/gh-action-runner:latest"
	_, _, err = cli.ImageInspectWithRaw(context.Background(), imageName)
	if err != nil {
		// 이미지 풀 시도
		_, err = cli.ImagePull(context.Background(), imageName, types.ImagePullOptions{})
		if err != nil {
			return nil, errors.New("failed to pull image")
		}

	}

	// 모든 조건을 만족하면 true 반환
	return cli, nil
}

func New(url, token string) (id string, err error) {
	// 도커 클라이언트 초기화
	ctx := context.Background()
	cli, err := set()
	if err != nil {
		return "", err
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "minpeter/gh-action-runner:latest",
		Env: []string{
			"URL=" + url,
			"TOKEN=" + token,
		},
		// 도커 소켓을 공유하여 호스트의 도커 데몬을 사용
		Volumes: map[string]struct{}{
			"/var/run/docker.sock": {},
		},
		Tty: false,
	}, nil, nil, nil, "")
	if err != nil {
		return "", errors.New("failed to create container")
	}

	// 컨테이너 시작
	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		return "", errors.New("failed to start container")
	}

	// return resp.ID, nil
	return resp.ID[:12], nil
}

func Remove(id, token string) error {
	// 도커 클라이언트 초기화
	ctx := context.Background()
	cli, err := set()
	if err != nil {
		return err
	}

	// exec ./config.sh remove --token + token 실행 & 완료 대기
	exec, err := cli.ContainerExecCreate(ctx, id, types.ExecConfig{
		Cmd: []string{"./config.sh", "remove", "--token", token},
	})
	if err != nil {
		return errors.New("failed to create exec")
	}
	// 완료까지 대기
	err = cli.ContainerExecStart(ctx, exec.ID, types.ExecStartCheck{})
	if err != nil {
		return errors.New("failed to start exec")
	}

	statusCh, errCh := cli.ContainerWait(ctx, id, container.WaitConditionNotRunning)
	select {
	case err := <-errCh:
		if err != nil {
			return errors.New("failed to wait exec")
		}
	case <-statusCh:
	}

	// 컨테이너 삭제
	err = cli.ContainerRemove(ctx, id, types.ContainerRemoveOptions{
		Force: true,
	})
	if err != nil {
		return errors.New("failed to remove container")
	}
	return nil
}

// gh-action-runner 이미지로 실행중인 모든 컨테이너의 상태와 ID를 반환
func List() (map[string]string, error) {
	// 도커 클라이언트 초기화
	ctx := context.Background()
	cli, err := set()
	if err != nil {
		return nil, err
	}

	// 모든 컨테이너의 상태를 조회
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, errors.New("failed to list containers")
	}

	// gh-action-runner 이미지로 실행중인 컨테이너의 상태와 ID를 반환
	list := make(map[string]string)
	for _, container := range containers {
		if container.Image == "minpeter/gh-action-runner:latest" {
			list[container.ID[:12]] = container.State
		}
	}

	return list, nil
}
