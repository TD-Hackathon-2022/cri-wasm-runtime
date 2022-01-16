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
	return &runtimeapi.ListImagesResponse{Images: nil}, nil
}

// ImageStatus returns the status of the image, returns nil if the image doesn't present.
func (ds *templateService) ImageStatus(
	_ context.Context,
	r *runtimeapi.ImageStatusRequest,
) (*runtimeapi.ImageStatusResponse, error) {
	logrus.Info("image status, image: %s", r.GetImage().GetImage())

	imageCache := ds.imageCache[r.GetImage().GetImage()]
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
	logrus.Info("pull image, image: %s", r.GetImage().GetImage())

	imageMockId := fmt.Sprintf("%x", md5.Sum([]byte(imageSpec.GetImage())))
	imageCache := &imageCacheModel{
		id:            imageMockId,
		name:          imageSpec.GetImage(),
		image:         imageSpec,
		sandboxConfig: ds.sandboxCache[r.GetSandboxConfig().GetMetadata().GetUid()].config,
		imageStatus: &runtimeapi.Image{
			Id:   imageMockId,
			Spec: imageSpec,
		},
	}
	ds.imageCache[r.GetImage().GetImage()] = imageCache
	return &runtimeapi.PullImageResponse{ImageRef: imageCache.id}, nil
}

// RemoveImage removes the image.
func (ds *templateService) RemoveImage(
	_ context.Context,
	r *runtimeapi.RemoveImageRequest,
) (*runtimeapi.RemoveImageResponse, error) {
	logrus.Info("remove image, image: %s", r.GetImage().GetImage())

	delete(ds.imageCache, r.GetImage().GetImage())
	return &runtimeapi.RemoveImageResponse{}, nil
}
