package ai

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/POSIdev-community/aictl/internal/adapter/ai/common/notify"
	"github.com/POSIdev-community/aictl/internal/core/apperror"
	"github.com/POSIdev-community/aictl/internal/core/domain/scanstage"
	"github.com/POSIdev-community/aictl/pkg/logger"
)

type accessTokenRefresher interface {
	RefreshAccessToken(ctx context.Context) error
}

type disconnectReason int

const (
	disconnectCtxDone disconnectReason = iota
	disconnectRefreshToken
	disconnectTransient
)

// WatchScanStage subscribes to PT AI notification hub and forwards progress updates
// for the given scan result. The subscription reconnects with exponential backoff on
// transient failures. On NeedRefreshToken (or 401) the JWT is refreshed and the
// websocket is reopened with the new access token.
//
// Fatal authentication failures are sent once on errc before both channels are closed.
// errc is closed when the watch loop exits for any reason.
func (a *Adapter) WatchScanStage(ctx context.Context, scanID uuid.UUID) (<-chan scanstage.ScanStage, <-chan error, error) {
	if a.baseClient == nil || a.baseClient.GetAccessToken() == "" {
		return nil, nil, fmt.Errorf("watch scan stage: adapter is not initialized")
	}

	client, err := notify.NewClient(notify.Options{
		BaseURL:     a.cfg.UriString(),
		AccessToken: a.baseClient.GetAccessToken(),
		HTTPClient:  a.baseClient.JwtHttpClient,
		TLSSkip:     a.cfg.TLSSkip(),
	})
	if err != nil {
		return nil, nil, fmt.Errorf("create notify client: %w", err)
	}

	out := make(chan scanstage.ScanStage, 8)
	errc := make(chan error, 1)
	go a.watchScanLoop(ctx, client, scanID, out, errc)

	return out, errc, nil
}

func (a *Adapter) watchScanLoop(
	ctx context.Context,
	client *notify.Client,
	scanID uuid.UUID,
	out chan<- scanstage.ScanStage,
	errc chan<- error,
) {
	defer close(out)
	defer close(errc)

	log := logger.FromContext(ctx)
	backoff := notify.ReconnectMinDelay

	for {
		if ctx.Err() != nil {
			return
		}

		client.SetAccessToken(a.baseClient.GetAccessToken())

		subCtx, cancel := context.WithCancel(ctx)
		messages, subErrc, err := client.Subscribe(subCtx)
		if err != nil {
			cancel()
			if ctx.Err() != nil {
				return
			}

			if notify.IsAuthError(err) {
				log.Debugf("notification auth failed, refreshing access token: %v", err)
				if refreshErr := a.refreshAccessToken(ctx); refreshErr != nil {
					log.Debugf("refresh access token failed: %v", refreshErr)
					if isAuthenticationError(refreshErr) {
						errc <- refreshErr

						return
					}
				} else {
					backoff = notify.ReconnectMinDelay
					continue
				}
			}

			log.Debugf("notification subscribe failed: %v; retry in %s", err, backoff)
			if !notify.Sleep(ctx, backoff) {
				return
			}
			backoff = notify.NextBackoff(backoff)

			continue
		}

		connectedAt := time.Now()

		reason := a.consumeSubscription(ctx, cancel, messages, subErrc, scanID, out)
		cancel()
		drainMessages(messages)

		if reason == disconnectCtxDone || ctx.Err() != nil {
			return
		}

		if reason == disconnectRefreshToken {
			log.Debugf("NeedRefreshToken received, refreshing access token")
			if err := a.refreshAccessToken(ctx); err != nil {
				log.Debugf("refresh access token after NeedRefreshToken: %v", err)
				if isAuthenticationError(err) {
					errc <- err

					return
				}
			}
			backoff = notify.ReconnectMinDelay

			continue
		}

		if time.Since(connectedAt) >= notify.StableConnectionForBackoffReset {
			backoff = notify.ReconnectMinDelay
		}

		if !notify.Sleep(ctx, backoff) {
			return
		}
		backoff = notify.NextBackoff(backoff)
	}
}

func (a *Adapter) consumeSubscription(
	ctx context.Context,
	cancel context.CancelFunc,
	messages <-chan notify.Message,
	errc <-chan error,
	scanID uuid.UUID,
	out chan<- scanstage.ScanStage,
) disconnectReason {
	log := logger.FromContext(ctx)

	for {
		select {
		case <-ctx.Done():
			return disconnectCtxDone
		case err, ok := <-errc:
			if ok && err != nil && ctx.Err() == nil {
				log.Debugf("notification subscription ended: %v", err)
			}
			if ctx.Err() != nil {
				return disconnectCtxDone
			}

			return disconnectTransient
		case msg, ok := <-messages:
			if !ok {
				if ctx.Err() != nil {
					return disconnectCtxDone
				}

				return disconnectTransient
			}

			if msg.NeedRefreshToken {
				cancel()

				return disconnectRefreshToken
			}

			stage, matched := scanStageFromNotification(msg, scanID)
			if !matched {
				continue
			}

			select {
			case out <- stage:
			case <-ctx.Done():
				return disconnectCtxDone
			}
		}
	}
}

func (a *Adapter) refreshAccessToken(ctx context.Context) error {
	refresher, ok := a.activeClient.(accessTokenRefresher)
	if !ok {
		return fmt.Errorf("active client does not support access token refresh")
	}

	if err := refresher.RefreshAccessToken(ctx); err != nil {
		return fmt.Errorf("refresh access token: %w", err)
	}

	return nil
}

func isAuthenticationError(err error) bool {
	var authErr *apperror.AuthenticationError

	return errors.As(err, &authErr)
}

func drainMessages(messages <-chan notify.Message) {
	for range messages {
	}
}

func scanStageFromNotification(msg notify.Message, scanID uuid.UUID) (scanstage.ScanStage, bool) {
	if msg.ScanProgress != nil && msg.ScanProgress.ScanResultID == scanID {
		stage := scanstage.ScanStage{}
		if msg.ScanProgress.Progress != nil {
			if msg.ScanProgress.Progress.Stage != nil {
				stage.Stage = *msg.ScanProgress.Progress.Stage
			}
			if msg.ScanProgress.Progress.SubStage != nil {
				stage.SubStage = *msg.ScanProgress.Progress.SubStage
			}
			if msg.ScanProgress.Progress.Value != nil {
				stage.Value = *msg.ScanProgress.Progress.Value
			}
		}
		if stage.Stage == "" {
			return scanstage.ScanStage{}, false
		}

		return stage, true
	}

	if msg.ScanCompleted != nil && msg.ScanCompleted.ScanResultID == scanID {
		stage := scanstage.ScanStage{Stage: "Done"}
		if msg.ScanCompleted.Stage != nil && *msg.ScanCompleted.Stage != "" {
			stage.Stage = *msg.ScanCompleted.Stage
		}

		return stage, true
	}

	return scanstage.ScanStage{}, false
}
