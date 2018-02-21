package main

import (
	"net/url"
	"strings"
	"net/http"
	"log"
	"encoding/json"
	"context"
	"encoding/pem"
	"crypto/x509"
	"time"
	"encoding/base64"
	"bytes"
	"crypto/sha1"
	"io"
	"io/ioutil"
	"crypto/rsa"
	"crypto"
)

// Request Functions
func (this *EchoRequest) VerifyTimestamp() bool {
	reqTimestamp, _ := time.Parse("2006-01-02T15:04:05Z", this.Request.Timestamp)
	if time.Since(reqTimestamp) < time.Duration(150)*time.Second {
		return true
	}

	return false
}

func (this *EchoRequest) VerifyAppID(myAppID string) bool {
	if this.Session.Application.ApplicationID == myAppID ||
		this.Context.System.Application.ApplicationID == myAppID {
		return true
	}

	return false
}

func (this *EchoRequest) GetSessionID() string {
	return this.Session.SessionID
}

func (this *EchoRequest) GetUserID() string {
	return this.Session.User.UserID
}

func (this *EchoRequest) GetRequestType() string {
	return this.Request.Type
}

func (this *EchoRequest) GetIntentName() string {
	if this.GetRequestType() == "IntentRequest" {
		return this.Request.Intent.Name
	}

	return this.GetRequestType()
}

func (this *EchoRequest) GetSlotValue(slotName string) (string, []byte) {
	if _, ok := this.Request.Intent.Slots[slotName]; ok {
		return this.Request.Intent.Slots[slotName].Value, nil
	}

	return "", []byte("Slot name not found.")
}

func (this *EchoRequest) AllSlots() map[string]EchoSlot {
	return this.Request.Intent.Slots
}

// Locale returns the locale specified in the request.
func (this *EchoRequest) Locale() string {
	return this.Request.Locale
}

func GetEchoRequest(r *http.Request) *EchoRequest {
	return r.Context().Value("echoRequest").(*EchoRequest)
}

func HTTPError(w http.ResponseWriter, logMsg string, err string, errCode int) {
	if logMsg != "" {
		log.Println(logMsg)
	}

	http.Error(w, err, errCode)
}

// Decode the JSON request and verify it.
func verifyJSON(w http.ResponseWriter, r *http.Request) bool {
	var echoReq *EchoRequest
	err := json.NewDecoder(r.Body).Decode(&echoReq)
	if err != nil {
		HTTPError(w, err.Error(), "Bad Request", 400)
		return false
	}

	// Check the timestamp
	if !echoReq.VerifyTimestamp() && r.URL.Query().Get("_dev") == "" {
		HTTPError(w, "Request too old to continue (>150s).", "Bad Request", 400)
		return false
	}

	// Check the app id
	if !echoReq.VerifyAppID("amzn1.ask.skill.6d962c3c-24ca-44cc-aca0-7713f23019f5") {
		HTTPError(w, "Echo AppID mismatch!", "Bad Request", 400)
		return false
	}

	r = r.WithContext(context.WithValue(r.Context(), "echoRequest", echoReq))

	return true
}

// Run all mandatory Amazon security checks on the request.
func validateRequest(w http.ResponseWriter, r *http.Request) bool {
	// Check for debug bypass flag
	devFlag := r.URL.Query().Get("_dev")

	isDev := devFlag != ""

	if !isDev {
		isRequestValid := IsValidAlexaRequest(w, r)
		if !isRequestValid {
			return false
		}
	}

	return true
}

// IsValidAlexaRequest handles all the necessary steps to validate that an incoming http.Request has actually come from
// the Alexa service. If an error occurs during the validation process, an http.Error will be written to the provided http.ResponseWriter.
// The required steps for request validation can be found on this page:
// https://developer.amazon.com/public/solutions/alexa/alexa-skills-kit/docs/developing-an-alexa-skill-as-a-web-service#hosting-a-custom-skill-as-a-web-service
func IsValidAlexaRequest(w http.ResponseWriter, r *http.Request) bool {
	certURL := r.Header.Get("SignatureCertChainUrl")

	// Verify certificate URL
	if !verifyCertURL(certURL) {
		HTTPError(w, "Invalid cert URL: "+certURL, "Not Authorized", 401)
		return false
	}

	// Fetch certificate data
	certContents, err := readCert(certURL)
	if err != nil {
		HTTPError(w, string(err), "Not Authorized", 401)
		return false
	}

	// Decode certificate data
	block, _ := pem.Decode(certContents)
	if block == nil {
		HTTPError(w, "Failed to parse certificate PEM.", "Not Authorized", 401)
		return false
	}

	cert, err1 := x509.ParseCertificate(block.Bytes)
	if err1 != nil {
		HTTPError(w, err1.Error(), "Not Authorized", 401)
		return false
	}

	// Check the certificate date
	if time.Now().Unix() < cert.NotBefore.Unix() || time.Now().Unix() > cert.NotAfter.Unix() {
		HTTPError(w, "Amazon certificate expired.", "Not Authorized", 401)
		return false
	}

	// Check the certificate alternate names
	foundName := false
	for _, altName := range cert.Subject.Names {
		if altName.Value == "echo-api.amazon.com" {
			foundName = true
		}
	}

	if !foundName {
		HTTPError(w, "Amazon certificate invalid.", "Not Authorized", 401)
		return false
	}

	// Verify the key
	publicKey := cert.PublicKey
	encryptedSig, _ := base64.StdEncoding.DecodeString(r.Header.Get("Signature"))

	// Make the request body SHA1 and verify the request with the public key
	var bodyBuf bytes.Buffer
	hash := sha1.New()
	_, err2 := io.Copy(hash, io.TeeReader(r.Body, &bodyBuf))
	if err2 != nil {
		HTTPError(w, err2.Error(), "Internal Error", 500)
		return false
	}
	//log.Println(bodyBuf.String())
	r.Body = ioutil.NopCloser(&bodyBuf)

	err3 := rsa.VerifyPKCS1v15(publicKey.(*rsa.PublicKey), crypto.SHA1, hash.Sum(nil), encryptedSig)
	if err3 != nil {
		HTTPError(w, "Signature match failed.", "Not Authorized", 401)
		return false
	}

	return true
}

func readCert(certURL string) ([]byte, []byte) {
	cert, err := http.Get(certURL)
	if err != nil {
		return nil, []byte("Could not download Amazon cert file.")
	}
	defer cert.Body.Close()
	certContents, err := ioutil.ReadAll(cert.Body)
	if err != nil {
		return nil, []byte("Could not read Amazon cert file.")
	}

	return certContents, nil
}

func verifyCertURL(path string) bool {
	link, _ := url.Parse(path)

	if link.Scheme != "https" {
		return false
	}

	if link.Host != "s3.amazonaws.com" && link.Host != "s3.amazonaws.com:443" {
		return false
	}

	if !strings.HasPrefix(link.Path, "/echo.api/") {
		return false
	}

	return true
}
