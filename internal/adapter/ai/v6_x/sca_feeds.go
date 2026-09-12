package v6_x

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/common"
	"github.com/POSIdev-community/aictl/internal/core/domain/scafeeds"
	"github.com/POSIdev-community/aictl/internal/core/domain/validation"
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

func (a *ClientAI6x) GetScaFeeds(ctx context.Context, statuses []scafeeds.Status) ([]scafeeds.Package, error) {
	pkgType := v6_x.ScaFeeds
	resp, err := a.GetPackagesWithResponse(ctx, &v6_x.GetPackagesParams{Type: &pkgType}, a.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("get sca feeds request: %w", err)
	}

	if err = CheckResponseByModel(resp.StatusCode(), string(resp.Body), resp.JSON400); err != nil {
		return nil, fmt.Errorf("get sca feeds: %w", err)
	}

	if resp.JSON200 == nil {
		return nil, fmt.Errorf("get sca feeds: empty response")
	}

	want := make(map[scafeeds.Status]struct{}, len(statuses))
	for _, s := range statuses {
		want[s] = struct{}{}
	}

	out := make([]scafeeds.Package, 0, len(*resp.JSON200))
	for _, p := range *resp.JSON200 {
		mapped := mapPackage(p)
		if len(want) > 0 {
			if _, ok := want[mapped.Status]; !ok {
				continue
			}
		}
		out = append(out, mapped)
	}

	return out, nil
}

func (a *ClientAI6x) DownloadScaFeeds(ctx context.Context, version string) (io.ReadCloser, string, error) {
	pkg, err := a.findScaFeedsByVersion(ctx, version)
	if err != nil {
		return nil, "", err
	}

	resp, err := a.DownloadPackage(ctx, pkg.Id, a.AddJWTToHeader)
	if err != nil {
		return nil, "", fmt.Errorf("download sca feeds request: %w", err)
	}

	if err = CheckResponse(resp, "sca feeds"); err != nil {
		_ = resp.Body.Close()

		return nil, "", fmt.Errorf("download sca feeds: %w", err)
	}

	fileName := pkg.FileName
	if fileName == "" {
		fileName = scafeeds.DefaultArchiveName(version)
	}

	total := pkg.FileSize
	if total <= 0 && resp.ContentLength > 0 {
		total = resp.ContentLength
	}

	log := logger.FromContext(ctx)
	onProgress := common.NewProgressReporter(log, true, "downloading sca feeds")

	return common.ProgressReadCloser(resp.Body, total, onProgress), fileName, nil
}

func (a *ClientAI6x) RollbackScaFeeds(ctx context.Context) (scafeeds.Package, error) {
	pkgs, err := a.GetScaFeeds(ctx, nil)
	if err != nil {
		return scafeeds.Package{}, err
	}

	var current *scafeeds.Package
	for i := range pkgs {
		if pkgs[i].Status == scafeeds.StatusCurrent {
			current = &pkgs[i]
			break
		}
	}
	if current == nil {
		return scafeeds.Package{}, validation.NewError(scafeeds.ErrNoCurrent)
	}

	resp, err := a.RollbackPackageWithResponse(ctx, current.Id, a.AddJWTToHeader)
	if err != nil {
		return scafeeds.Package{}, fmt.Errorf("rollback sca feeds request: %w", err)
	}

	errorModel := resp.JSON400
	switch resp.StatusCode() {
	case http.StatusNotFound:
		errorModel = resp.JSON404
	case http.StatusConflict:
		errorModel = resp.JSON409
	}
	if err = CheckResponseByModel(resp.StatusCode(), string(resp.Body), errorModel); err != nil {
		return scafeeds.Package{}, fmt.Errorf("rollback sca feeds: %w", err)
	}

	if resp.JSON200 == nil {
		return scafeeds.Package{}, fmt.Errorf("rollback sca feeds: empty response")
	}

	return mapPackage(*resp.JSON200), nil
}

func (a *ClientAI6x) findScaFeedsByVersion(ctx context.Context, version string) (scafeeds.Package, error) {
	pkgs, err := a.GetScaFeeds(ctx, nil)
	if err != nil {
		return scafeeds.Package{}, err
	}

	for _, p := range pkgs {
		if p.Version == version {
			return p, nil
		}
	}

	return scafeeds.Package{}, validation.NewError(fmt.Sprintf("sca-feeds with version '%s' not found", version))
}

func mapPackage(p v6_x.Package) scafeeds.Package {
	return scafeeds.Package{
		Id:             p.Id,
		PackageType:    scafeeds.Type(p.PackageType),
		Version:        p.Version,
		Status:         scafeeds.Status(p.Status),
		FileName:       p.FileName,
		FileSize:       p.FileSize,
		FileHash:       p.FileHash,
		TriggeredBy:    scafeeds.Trigger(p.TriggeredBy),
		UploadedAt:     p.UploadedAt,
		UploadedBy:     mapPackageUser(p.UploadedBy),
		LastModifiedAt: p.LastModifiedAt,
		LastModifiedBy: mapPackageUser(p.LastModifiedBy),
	}
}

func mapPackageUser(u *v6_x.PackageUser) *scafeeds.PackageUser {
	if u == nil {
		return nil
	}

	return &scafeeds.PackageUser{
		Id:        uuid.UUID(u.Id),
		UserName:  u.UserName,
		Email:     u.Email,
		TokenName: u.TokenName,
	}
}
