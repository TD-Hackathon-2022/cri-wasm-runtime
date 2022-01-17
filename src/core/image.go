package core

import (
	"context"
	"crypto/md5"
	"fmt"
	"github.com/sirupsen/logrus"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1alpha2"
)

// This file implements methods in ImageManagerService.

// ListImages lists existing images.
func (ds *templateService) ListImages(
	_ context.Context,
	r *runtimeapi.ListImagesRequest,
) (*runtimeapi.ListImagesResponse, error) {
	// todo something here
	//logrus.Infof("list image, image count : %d", len(ds.imageCache))
	items := make([]*runtimeapi.Image, 0, len(ds.imageCache))
	for id, cache := range ds.imageCache {
		item := &runtimeapi.Image{
			Id:          id,
			RepoTags:    cache.imageStatus.GetRepoTags(),
			RepoDigests: cache.imageStatus.GetRepoDigests(),
			Size_:       12345,
			Spec:        cache.image,
		}
		items = append(items, item)
	}

	defer logrus.Infof("end list image")
	return &runtimeapi.ListImagesResponse{Images: items}, nil
}

// ImageStatus returns the status of the image, returns nil if the image doesn't present.
func (ds *templateService) ImageStatus(
	_ context.Context,
	r *runtimeapi.ImageStatusRequest,
) (*runtimeapi.ImageStatusResponse, error) {
	logrus.Infof("image status, image: %s, image count : %d", r.GetImage().GetImage(), len(ds.imageCache))
	defer logrus.Infof("end status image, image: %s", r.GetImage().GetImage())
	imageId := fmt.Sprintf("%x", md5.Sum([]byte(r.GetImage().GetImage())))
	imageCache := ds.imageCache[imageId]
	if imageCache == nil {
		return &runtimeapi.ImageStatusResponse{Image: nil}, nil
	}
	return &runtimeapi.ImageStatusResponse{Image: imageCache.imageStatus}, nil
}

// PullImage pulls an image with authentication config.
func (ds *templateService) PullImage(
	_ context.Context,
	r *runtimeapi.PullImageRequest,
) (*runtimeapi.PullImageResponse, error) {
	imageSpec := r.GetImage()
	logrus.Infof("pull image, image: %s, image count : %d", r.GetImage().GetImage(), len(ds.imageCache))
	defer logrus.Infof("end pull image, image: %s", r.GetImage().GetImage())

	imageMockId := fmt.Sprintf("%x", md5.Sum([]byte(imageSpec.GetImage())))
	imageCache := &imageCacheModel{
		id:            imageMockId,
		name:          imageSpec.GetImage(),
		image:         imageSpec,
		sandboxConfig: ds.sandboxCache[r.GetSandboxConfig().GetMetadata().GetUid()].config,
		imageStatus: &runtimeapi.Image{
			Id:    imageMockId,
			Spec:  imageSpec,
			Size_: 12345,
		},
	}
	ds.imageCache[imageMockId] = imageCache
	return &runtimeapi.PullImageResponse{ImageRef: imageMockId}, nil
}

// RemoveImage removes the image.
func (ds *templateService) RemoveImage(
	_ context.Context,
	r *runtimeapi.RemoveImageRequest,
) (*runtimeapi.RemoveImageResponse, error) {
	logrus.Infof("remove image, image: %s, image count : %d", r.GetImage().GetImage(), len(ds.imageCache))
	defer logrus.Infof("end remove image, image: %s", r.GetImage().GetImage())
	imageId := fmt.Sprintf("%x", md5.Sum([]byte(r.GetImage().GetImage())))
	delete(ds.imageCache, imageId)
	return &runtimeapi.RemoveImageResponse{}, nil
}
