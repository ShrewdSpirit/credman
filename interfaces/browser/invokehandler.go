package browser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	m "github.com/ShrewdSpirit/credman/interfaces/browser/methods"
	"github.com/ShrewdSpirit/ezsec"
	"github.com/labstack/echo/v4"
)

var methods = make(map[string]interface{})
var clientKeys = make(map[string][]byte)
var rsaPublicKey string
var rsaPrivateKey string

func initInvokeHandler() {
	priv, pub, err := ezsec.GenerateRSAKeyPair(2048)
	if err != nil {
		log.Fatalf("RSA keys generation failed: %s\n", err.Error())
	}

	pubBytes, err := ezsec.RSAPublicKeyToBytes(pub)
	if err != nil {
		log.Fatalf("Failed to convert public key: %s\n", err.Error())
	}

	privBytes := ezsec.RSAPrivateKeyToBytes(priv)
	rsaPublicKey = string(pubBytes)
	rsaPrivateKey = string(privBytes)

	methods["handshake_getkey"] = handshakeGetKey{}
	methods["handshake_setkey"] = handshakeSetKey{}
	methods["getinfo"] = m.GetInfo{}
	methods["test"] = m.Test{}
}

func invokeErrorResult(clientId string, encrypt bool, ctx echo.Context, message string, err error) error {
	msg := fmt.Sprintf("%s: %s", message, err.Error())
	server.Logger.Error("INVOKE ERROR:", msg)

	if encrypt {
		data, err := invokeEncryptBody([]byte(msg), clientId)
		if err != nil {
			server.Logger.Error(err)
			return err
		}

		return ctx.Blob(http.StatusBadRequest, echo.MIMEOctetStream, data)
	}

	return ctx.String(http.StatusBadRequest, msg)
}

func invokeSuccessJsonResult(clientId string, encrypt bool, ctx echo.Context, r interface{}) error {
	if encrypt {
		dataJson, err := json.Marshal(r)
		if err != nil {
			server.Logger.Error(err)
			return err
		}

		data, err := invokeEncryptBody(dataJson, clientId)
		if err != nil {
			server.Logger.Error(err)
			return err
		}

		return ctx.Blob(http.StatusOK, echo.MIMEOctetStream, data)
	}

	return ctx.JSON(http.StatusOK, r)
}

func invokeDecryptBody(body []byte, clientId string) ([]byte, error) {
	clientKey, ok := clientKeys[clientId]
	if !ok {
		return nil, errors.New("Client can't use encryption")
	}

	decoded, err := ezsec.Base64Decode(string(body))
	if err != nil {
		return nil, err
	}

	data, err := ezsec.CFBDecrypt(ezsec.ShaTypeSha512, decoded, clientKey)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func invokeEncryptBody(body []byte, clientId string) ([]byte, error) {
	clientKey, ok := clientKeys[clientId]
	if !ok {
		return nil, errors.New("Client can't use encryption")
	}

	data, err := ezsec.CFBEncrypt(ezsec.ShaTypeSha512, body, clientKey)
	if err != nil {
		return nil, err
	}

	return data, nil
}

type invokeData struct {
	Method string `json:"method"`
	Args   string `json:"args"`
}

func invokeHandler(ctx echo.Context) error {
	encrypted := ctx.Request().Header.Get("Encrypted") == "true"
	clientId := ctx.Request().Header.Get("Client-Id")

	body, err := ioutil.ReadAll(ctx.Request().Body)
	if err != nil {
		return invokeErrorResult(clientId, encrypted, ctx, "Failed to read request body", err)
	}

	if encrypted {
		body, err = invokeDecryptBody(body, clientId)
		if err != nil {
			return invokeErrorResult(clientId, encrypted, ctx, "Decryption error", err)
		}
	}

	var data invokeData
	if err := json.Unmarshal(body, &data); err != nil {
		return invokeErrorResult(clientId, encrypted, ctx, "Invalid request JSON", err)
	}

	methodStructType, ok := methods[data.Method]
	if !ok {
		return invokeErrorResult(clientId, encrypted, ctx, "Invalid method", errors.New("Method "+data.Method+" not found"))
	}

	methodInstance := reflect.New(reflect.TypeOf(methodStructType))
	if err := json.Unmarshal([]byte(data.Args), methodInstance.Interface()); err != nil {
		return invokeErrorResult(clientId, encrypted, ctx, "Invalid args", errors.New("Invalid arguments to method "+data.Method))
	}

	result, err := methodInstance.Interface().(m.MethodInterface).Do()
	if err != nil {
		return invokeErrorResult(clientId, encrypted, ctx, "Call error", err)
	}

	return invokeSuccessJsonResult(clientId, encrypted, ctx, result)
}

type handshakeGetKey struct{}

func (s handshakeGetKey) Do() (m.MethodResult, error) {
	return m.MethodResult{
		"publickey": rsaPublicKey,
	}, nil
}

type handshakeSetKey struct {
	ClientId  string `json:"p0"`
	ClientKey []byte `json:"p1"`
}

func (s handshakeSetKey) Do() (m.MethodResult, error) {
	priv, err := ezsec.RSABytesToPrivateKey([]byte(rsaPrivateKey))
	if err != nil {
		return nil, err
	}

	key, err := ezsec.RSADecrypt(ezsec.ShaTypeSha512, s.ClientKey, priv)
	if err != nil {
		return nil, err
	}

	clientKeys[s.ClientId] = key

	return nil, nil
}
