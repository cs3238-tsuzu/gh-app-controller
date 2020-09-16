package appstransport

import (
	"crypto/rsa"
	"net/http"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/xerrors"
)

// NewAppsTransportFromBytes returns a new AppsTransport from a private key in []byte
func NewAppsTransportFromBytes(appID int64, privateKey []byte) (*AppsTransport, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(privateKey)
	if err != nil {
		return nil, xerrors.Errorf("could not parse private key: %s", err)
	}
	return NewAppsTransport(appID, key), nil
}

// NewAppsTransport returns a new AppsTransport from rsa.PrivateKey
func NewAppsTransport(appID int64, key *rsa.PrivateKey) *AppsTransport {
	return &AppsTransport{
		Transport: http.DefaultTransport,

		appID: appID,
		key:   key,
	}
}

// AppsTransport returns a transport that inserts JWT into header
type AppsTransport struct {
	Transport http.RoundTripper

	appID int64
	key   *rsa.PrivateKey
}

// RoundTrip - http.RoundTripper
func (at *AppsTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	now := time.Now()
	st := &jwt.StandardClaims{
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(time.Minute).Unix(),
		Issuer:    strconv.FormatInt(at.appID, 10),
	}
	tk := jwt.NewWithClaims(jwt.SigningMethodRS256, st)

	token, err := tk.SignedString(at.key)
	if err != nil {
		return nil, xerrors.Errorf("signing failed: %s", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Add("Accept", "application/vnd.github.machine-man-preview+json")

	resp, err := at.Transport.RoundTrip(req)

	if err != nil {
		return nil, xerrors.Errorf("roundtrip failed: %w", err)
	}

	return resp, err
}
