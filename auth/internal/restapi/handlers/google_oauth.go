package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"log/slog"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/h4x4d/parking_net/auth/internal/models"
	"github.com/h4x4d/parking_net/auth/internal/restapi/operations"
	"github.com/h4x4d/parking_net/auth/internal/utils"
)

func (h *Handler) GoogleLoginHandler(params operations.GetAuthGoogleLoginParams) middleware.Responder {
	var responder middleware.Responder
	defer utils.CatchPanic(&responder)

	_, span := h.tracer.Start(context.Background(), "google_login")
	defer span.End()

	traceID := fmt.Sprintf("%s", span.SpanContext().TraceID())

	keycloakURL := os.Getenv("KEYCLOAK_FRONTEND_URL")
	if keycloakURL == "" {
		keycloakSubdomain := os.Getenv("KEYCLOAK_SUBDOMAIN")
		if keycloakSubdomain != "" {
			keycloakURL = fmt.Sprintf("https://%s", keycloakSubdomain)
		} else {
			keycloakURL = fmt.Sprintf("http://keycloak:%s", os.Getenv("KEYCLOAK_PORT"))
		}
	}
	realm := os.Getenv("KEYCLOAK_REALM")
	if realm == "" {
		realm = "parking-users"
	}

	redirectURI := os.Getenv("GOOGLE_OAUTH_REDIRECT_URI")
	if redirectURI == "" {
		backendDomain := os.Getenv("BACKEND_DOMAIN")
		if backendDomain != "" {
			redirectURI = fmt.Sprintf("https://%s/auth/google/callback", backendDomain)
		} else {
			authServiceURL := os.Getenv("AUTH_SERVICE_URL")
			if authServiceURL != "" {
				redirectURI = fmt.Sprintf("%s/auth/google/callback", authServiceURL)
			} else {
				redirectURI = fmt.Sprintf("http://localhost:%s/auth/google/callback", os.Getenv("AUTH_REST_PORT"))
			}
		}
	}

	keycloakAuthURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/auth", keycloakURL, realm)
	authURL, err := url.Parse(keycloakAuthURL)
	if err != nil {
		slog.Error(
			"failed to parse keycloak auth URL",
			slog.String("method", "GET"),
			slog.String("trace_id", traceID),
			slog.String("error", err.Error()),
		)
		errCode := int64(http.StatusInternalServerError)
		responder = operations.NewGetAuthGoogleLoginInternalServerError().WithPayload(&models.Error{
			ErrorMessage:    "Internal server error",
			ErrorStatusCode: &errCode,
		})
		return responder
	}

	clientID := os.Getenv("KEYCLOAK_CLIENT")
	if clientID == "" {
		clientID = "parking-auth"
	}

	q := authURL.Query()
	q.Set("client_id", clientID)
	q.Set("redirect_uri", redirectURI)
	q.Set("response_type", "code")
	q.Set("scope", "openid profile email")
	q.Set("kc_idp_hint", "google")
	authURL.RawQuery = q.Encode()

	slog.Info(
		"redirecting to google oauth",
		slog.String("method", "GET"),
		slog.String("trace_id", traceID),
		slog.String("redirect_url", authURL.String()),
	)

	responder = middleware.ResponderFunc(func(w http.ResponseWriter, p runtime.Producer) {
		http.Redirect(w, params.HTTPRequest, authURL.String(), http.StatusFound)
	})
	return responder
}

func (h *Handler) GoogleCallbackHandler(params operations.GetAuthGoogleCallbackParams) middleware.Responder {
	var responder middleware.Responder
	defer utils.CatchPanic(&responder)

	ctx, span := h.tracer.Start(context.Background(), "google_callback")
	defer span.End()

	traceID := fmt.Sprintf("%s", span.SpanContext().TraceID())

	if params.Error != nil && *params.Error != "" {
		slog.Error(
			"google oauth error",
			slog.String("method", "GET"),
			slog.String("trace_id", traceID),
			slog.String("error", *params.Error),
		)
		errCode := int64(http.StatusUnauthorized)
		responder = operations.NewGetAuthGoogleCallbackUnauthorized().WithPayload(&models.Error{
			ErrorMessage:    "OAuth authentication failed",
			ErrorStatusCode: &errCode,
		})
		return responder
	}

	if params.Code == nil || *params.Code == "" {
		slog.Error(
			"missing authorization code",
			slog.String("method", "GET"),
			slog.String("trace_id", traceID),
		)
		errCode := int64(http.StatusBadRequest)
		responder = operations.NewGetAuthGoogleCallbackUnauthorized().WithPayload(&models.Error{
			ErrorMessage:    "Missing authorization code",
			ErrorStatusCode: &errCode,
		})
		return responder
	}

	keycloakURL := os.Getenv("KEYCLOAK_FRONTEND_URL")
	if keycloakURL == "" {
		keycloakSubdomain := os.Getenv("KEYCLOAK_SUBDOMAIN")
		if keycloakSubdomain != "" {
			keycloakURL = fmt.Sprintf("https://%s", keycloakSubdomain)
		} else {
			keycloakURL = fmt.Sprintf("http://keycloak:%s", os.Getenv("KEYCLOAK_PORT"))
		}
	}
	realm := os.Getenv("KEYCLOAK_REALM")
	if realm == "" {
		realm = "parking-users"
	}

	redirectURI := os.Getenv("GOOGLE_OAUTH_REDIRECT_URI")
	if redirectURI == "" {
		backendDomain := os.Getenv("BACKEND_DOMAIN")
		if backendDomain != "" {
			redirectURI = fmt.Sprintf("https://%s/auth/google/callback", backendDomain)
		} else {
			authServiceURL := os.Getenv("AUTH_SERVICE_URL")
			if authServiceURL != "" {
				redirectURI = fmt.Sprintf("%s/auth/google/callback", authServiceURL)
			} else {
				redirectURI = fmt.Sprintf("http://localhost:%s/auth/google/callback", os.Getenv("AUTH_REST_PORT"))
			}
		}
	}

	clientID := os.Getenv("KEYCLOAK_CLIENT")
	if clientID == "" {
		clientID = "parking-auth"
	}
	clientSecret := os.Getenv("KEYCLOAK_CLIENT_SECRET")

	keycloakInternalURL := os.Getenv("KEYCLOAK_URL")
	if keycloakInternalURL == "" {
		keycloakInternalURL = fmt.Sprintf("http://keycloak:%s", os.Getenv("KEYCLOAK_PORT"))
	}
	tokenURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", keycloakInternalURL, realm)

	token, err := h.exchangeCodeForToken(ctx, tokenURL, *params.Code, redirectURI, clientID, clientSecret)
	if err != nil {
		slog.Error(
			"failed to exchange code for token",
			slog.String("method", "GET"),
			slog.String("trace_id", traceID),
			slog.String("error", err.Error()),
		)
		errCode := int64(http.StatusInternalServerError)
		responder = operations.NewGetAuthGoogleCallbackInternalServerError().WithPayload(&models.Error{
			ErrorMessage:    "Failed to authenticate",
			ErrorStatusCode: &errCode,
		})
		return responder
	}

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000"
	}

	redirectURL := fmt.Sprintf("%s/auth/callback?token=%s", frontendURL, token)

	slog.Info(
		"google oauth success",
		slog.String("method", "GET"),
		slog.String("trace_id", traceID),
	)

	responder = middleware.ResponderFunc(func(w http.ResponseWriter, p runtime.Producer) {
		http.Redirect(w, params.HTTPRequest, redirectURL, http.StatusFound)
	})
	return responder
}

func (h *Handler) exchangeCodeForToken(ctx context.Context, tokenURL, code, redirectURI, clientID, clientSecret string) (string, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to exchange code: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		slog.Error("token exchange failed",
			slog.Int("status_code", resp.StatusCode),
			slog.String("response_body", string(bodyBytes)))
		return "", fmt.Errorf("token exchange failed with status: %d", resp.StatusCode)
	}

	var tokenResponse struct {
		AccessToken  string `json:"access_token"`
		TokenType    string `json:"token_type"`
		RefreshToken string `json:"refresh_token,omitempty"`
		ExpiresIn    int    `json:"expires_in,omitempty"`
		Scope        string `json:"scope,omitempty"`
	}

	bodyBytes, _ := io.ReadAll(resp.Body)
	if err := json.Unmarshal(bodyBytes, &tokenResponse); err != nil {
		slog.Error("failed to decode token response",
			slog.String("error", err.Error()),
			slog.String("status", fmt.Sprintf("%d", resp.StatusCode)),
			slog.String("body", string(bodyBytes)))
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	if tokenResponse.AccessToken == "" {
		slog.Error("token response has empty access_token")
		return "", fmt.Errorf("access_token is empty in response")
	}

	slog.Info("token exchange successful",
		slog.String("token_type", tokenResponse.TokenType),
		slog.String("token_prefix", tokenResponse.AccessToken[:min(30, len(tokenResponse.AccessToken))]),
		slog.Int("expires_in", tokenResponse.ExpiresIn))

	return tokenResponse.AccessToken, nil
}
