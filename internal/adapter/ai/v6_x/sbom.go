package v6_x

import (
	"context"
	"fmt"
	"net/http"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/common"
	"github.com/POSIdev-community/aictl/internal/core/domain/project"
	"github.com/POSIdev-community/aictl/internal/core/domain/scantype"
	"github.com/POSIdev-community/aictl/pkg/clientai/v6_x"
)

func mapProjectTargetType(t *v6_x.ProjectTargetBasedType) project.Type {
	if t == nil {
		return project.TypeSource
	}
	if *t == v6_x.SbomBased {
		return project.TypeSbom
	}

	return project.TypeSource
}

func projectFromModel(id uuid.UUID, name string, target *v6_x.ProjectTargetBasedType) project.Project {
	return project.NewProjectWithType(id, name, mapProjectTargetType(target))
}

func (a *ClientAI6x) CreateSbomProject(ctx context.Context, projectName string) (*uuid.UUID, error) {
	projectURL := "http://localhost"
	patternMatchingEnabled := false
	searchForConfigurationFlawsEnabled := false
	searchForVulnerableComponentsEnabled := false
	searchWithScaEnabled := true
	secretDetectionEnabled := false
	searchForMaliciousCodeEnabled := false
	staticCodeAnalysisEnabled := false
	preferredAgentsOnly := false
	preferredAgents := []uuid.UUID{}
	priority := v6_x.PriorityMedium
	targetType := v6_x.SbomBased
	languages := []v6_x.LegacyProgrammingLanguageGroup{}

	projectBaseModel := v6_x.PostApiProjectsBaseJSONRequestBody{
		Name:                   &projectName,
		ProjectUrl:             &projectURL,
		ProjectTargetBasedType: &targetType,
		WhiteBox: &v6_x.WhiteBoxSettingsModel{
			PatternMatchingEnabled:               &patternMatchingEnabled,
			SearchForConfigurationFlawsEnabled:   &searchForConfigurationFlawsEnabled,
			SearchForVulnerableComponentsEnabled: &searchForVulnerableComponentsEnabled,
			SearchWithScaEnabled:                 &searchWithScaEnabled,
			SecretDetectionEnabled:               &secretDetectionEnabled,
			SearchForMaliciousCodeEnabled:        &searchForMaliciousCodeEnabled,
			StaticCodeAnalysisEnabled:            &staticCodeAnalysisEnabled,
		},
		Id:        &uuid.UUID{},
		Languages: &languages,
		PreferredAgentsSettings: &v6_x.PreferredAgentsSettings{
			PreferredAgents:     &preferredAgents,
			PreferredAgentsOnly: preferredAgentsOnly,
		},
		Priority: &priority,
	}

	createProjectResponse, err := a.PostApiProjectsBaseWithResponse(ctx, projectBaseModel, a.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("create sbom project request error: %w", err)
	}

	statusCode := createProjectResponse.StatusCode()
	body := string(createProjectResponse.Body)
	errorModel := createProjectResponse.JSON400
	if err = CheckResponseByModel(statusCode, body, errorModel); err != nil {
		return nil, err
	}

	projectId, err := uuid.Parse(body)
	if err != nil {
		return nil, err
	}

	return &projectId, nil
}

func (a *ClientAI6x) GetProjectByName(ctx context.Context, projectName string) (*project.Project, error) {
	response, err := a.GetApiProjectsNameNameWithResponse(ctx, projectName, a.AddJWTToHeader)
	if err != nil {
		return nil, fmt.Errorf("ai adapter get project by name request: %w", err)
	}

	statusCode := response.StatusCode()
	body := string(response.Body)
	errorModel := response.JSON400

	if statusCode == http.StatusBadRequest && errorModel != nil && *errorModel.ErrorCode == v6_x.ApiErrorTypePROJECTNOTFOUND {
		return nil, nil
	}

	if err = CheckResponseByModel(statusCode, body, errorModel); err != nil {
		return nil, err
	}
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("ai adapter get project by name: unexpected status %d", statusCode)
	}

	model := *response.JSON200
	p := projectFromModel(*model.Id, *model.Name, model.ProjectTargetBasedType)

	return &p, nil
}

func (a *ClientAI6x) UpdateSbom(ctx context.Context, projectId uuid.UUID, sbomPath string) error {
	sessionResp, err := a.PostApiStoreUploadSessionWithResponse(ctx, a.AddJWTToHeader)
	if err != nil {
		return fmt.Errorf("create upload session request: %w", err)
	}
	if err = CheckResponseByModel(sessionResp.StatusCode(), string(sessionResp.Body), sessionResp.JSON400); err != nil {
		return fmt.Errorf("create upload session: %w", err)
	}
	if sessionResp.JSON200 == nil || sessionResp.JSON200.Id == nil {
		return fmt.Errorf("create upload session: empty response")
	}
	uploadId := *sessionResp.JSON200.Id

	if err = a.uploadSessionAddFile(ctx, uploadId, sbomPath); err != nil {
		_, _ = a.PostApiStoreUploadSessionUploadIdCancelWithResponse(ctx, uploadId, a.AddJWTToHeader)

		return err
	}

	replaceResp, err := a.PostApiStoreUploadSessionUploadIdProjectProjectIdSbomWithResponse(ctx, uploadId, projectId, a.AddJWTToHeader)
	if err != nil {
		_, _ = a.PostApiStoreUploadSessionUploadIdCancelWithResponse(ctx, uploadId, a.AddJWTToHeader)

		return fmt.Errorf("replace sbom request: %w", err)
	}
	if err = CheckResponseByModel(replaceResp.StatusCode(), string(replaceResp.Body), replaceResp.JSON400); err != nil {
		_, _ = a.PostApiStoreUploadSessionUploadIdCancelWithResponse(ctx, uploadId, a.AddJWTToHeader)

		return fmt.Errorf("replace sbom: %w", err)
	}

	return nil
}

func (a *ClientAI6x) uploadSessionAddFile(ctx context.Context, uploadId uuid.UUID, filePath string) error {
	body, contentType, err := common.PrepareMultipartBody(ctx, filePath, true)
	if err != nil {
		return err
	}
	defer func() { _ = body.Close() }()

	client, ok := a.ClientInterface.(*v6_x.Client)
	if !ok {
		return fmt.Errorf("upload session add: unexpected client type")
	}

	req, err := v6_x.NewPostApiStoreUploadSessionUploadIdAddRequest(client.Server, uploadId)
	if err != nil {
		return fmt.Errorf("upload session add request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Body = body
	req.Header.Set("Content-Type", contentType)

	if err = a.AddJWTToHeader(ctx, req); err != nil {
		return fmt.Errorf("upload session add auth: %w", err)
	}

	resp, err := client.Client.Do(req)
	if err != nil {
		return fmt.Errorf("upload session add: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	parsed, err := v6_x.ParsePostApiStoreUploadSessionUploadIdAddResponse(resp)
	if err != nil {
		return fmt.Errorf("upload session add parse: %w", err)
	}
	if err = CheckResponseByModel(parsed.StatusCode(), string(parsed.Body), parsed.JSON400); err != nil {
		return fmt.Errorf("upload session add: %w", err)
	}

	return nil
}

func (a *ClientAI6x) StartScanSbom(ctx context.Context, projectId uuid.UUID, scanLabel string) (uuid.UUID, error) {
	response, err := a.GetApiProjectsProjectIdWithResponse(ctx, projectId, a.AddJWTToHeader)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%w: get project: %w", project.ErrCannotStartSbomScan, err)
	}
	if err = CheckResponseByModel(response.StatusCode(), string(response.Body), response.JSON400); err != nil {
		return uuid.UUID{}, fmt.Errorf("%w: get project: %w", project.ErrCannotStartSbomScan, err)
	}
	if response.JSON200 == nil || response.JSON200.VirtualBranchId == nil {
		return uuid.UUID{}, fmt.Errorf("%w: virtual branch not found", project.ErrCannotStartSbomScan)
	}

	branchId := uuid.UUID(*response.JSON200.VirtualBranchId)
	scanId, err := a.StartScanBranch(ctx, branchId, scanLabel, scantype.Incremental)
	if err != nil {
		return uuid.UUID{}, err
	}

	return scanId, nil
}
