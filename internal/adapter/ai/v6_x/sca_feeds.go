package v6_x

import (
	"context"
	"fmt"
	"net/http"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/common"
	"github.com/POSIdev-community/aictl/pkg/clientai/v6_x"
	"github.com/POSIdev-community/aictl/pkg/logger"
)

func (a *ClientAI6x) UpdateScaFeeds(ctx context.Context, path, version string) error {
	meta, err := common.ComputePackageFileMeta(path)
	if err != nil {
		return fmt.Errorf("package metadata: %w", err)
	}

	log := logger.FromContext(ctx)
	onProgress := common.NewProgressReporter(log, true, "uploading sca feeds")

	body, contentType, err := common.PreparePackageMultipartBody(ctx, path, version, meta, onProgress)
	if err != nil {
		return fmt.Errorf("prepare package body: %w", err)
	}
	defer func() { _ = body.Close() }()

	resp, err := a.UploadPackageWithBodyWithResponse(ctx, v6_x.ScaFeeds, contentType, body, a.AddJWTToHeader)
	if err != nil {
		return fmt.Errorf("upload sca feeds request: %w", err)
	}

	errorModel := resp.JSON400
	if resp.StatusCode() == http.StatusConflict {
		errorModel = resp.JSON409
	}
	if err = CheckResponseByModel(resp.StatusCode(), string(resp.Body), errorModel); err != nil {
		return fmt.Errorf("upload sca feeds: %w", err)
	}

	return nil
}
