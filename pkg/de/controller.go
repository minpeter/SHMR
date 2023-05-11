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

	return resp.ID, nil
}

func Remove(id string) error {
	// 도커 클라이언트 초기화
	ctx := context.Background()
	cli, err := set()
	if err != nil {
		return err
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
			list[container.ID] = container.State
		}
	}

	return list, nil
}
