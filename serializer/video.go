package serializer

import (
	"giligili/model"
)

// Video 视频序列化器
type Video struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	Info      string `json:"info"`
	Cover     string `json:"cover"`
	VideoUrl  string `json:"video_url"`
	User      User   `json:"user"`
	CreatedAt int64  `json:"created_at"`
}

// BuildVideo 序列化用户
func BuildVideo(video model.Video) Video {
	return Video{
		ID:        video.ID,
		Title:     video.Title,
		Info:      video.Info,
		Cover:     video.Cover,
		VideoUrl:  video.VideoUrl,
		User:      BuildUser(video.User),
		CreatedAt: video.CreatedAt.Unix(),
	}
}

// BuildVideos 序列化视频列表
func BuildVideos(items []model.Video) (videos []Video) {
	for _, item := range items {
		video := BuildVideo(item)

		videos = append(videos, video)
	}

	return videos
}

// BuildVideoResponse 序列化视频响应
func BuildVideoResponse(video model.Video) Response {
	return Response{
		Data: BuildVideo(video),
	}
}

// BuildVideosResponse 序列化视频列表响应
func BuildVideosResponse(videos []model.Video) Response {
	return Response{
		Data: BuildVideos(videos),
	}
}
