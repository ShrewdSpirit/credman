package browser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"

	"github.com/ShrewdSpirit/ezsec"

	m "github.com/ShrewdSpirit/credman/interfaces/browser/methods"
	"github.com/labstack/echo/v4"
)

var methods = make(map[string]interface{})
var clientKeys = make(map[string][]byte)
var rsaPublicKey string
var rsaPrivateKey string

type invokeData struct {
	Method string      `json:"method"`
	Args   interface{} `json:"args"`
}

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

	methods["getinfo"] = m.GetInfo
	methods["handshake_getkey"] = handshakeGetKey
	methods["handshake_setkey"] = handshakeSetKey
}

func invokeErrorResult(clientId string, encrypt bool, ctx echo.Context, message string, err error) error {
	msg := fmt.Sprintf("%s: %s", message, err.Error())
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

	data, err := ezsec.CFBDecrypt(ezsec.ShaTypeSha512, body, clientKey)
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

	methodFunc, ok := methods[data.Method]
	if !ok {
		return invokeErrorResult(clientId, encrypted, ctx, "Invalid method", errors.New("Method "+data.Method+" not found"))
	}

	methodRef := reflect.ValueOf(methodFunc)
	methodType := methodRef.Type()
	numParams := methodType.NumIn()

	argsRef := reflect.ValueOf(data.Args)

	if argsRef.Len() != numParams {
		return invokeErrorResult(clientId, encrypted, ctx, "Call error", errors.New("Invalid number of arguments to method "+data.Method))
	}

	params := make([]reflect.Value, numParams)
	for i := 0; i < numParams; i++ {
		p := methodType.In(i)
		argRef := argsRef.Index(i)
		argType := reflect.TypeOf(argRef)

		if p != argType {
			if argType.Kind() == reflect.Slice {
				fmt.Printf("Slice of %s\n", argRef.Type().Elem().Kind().String())
			}
			return invokeErrorResult(clientId, encrypted, ctx, "Call error", errors.New(
				fmt.Sprintf("Invalid argument[%d] type for method %s. Expected %s got %s: %v",
					i, data.Method, p.Kind().String(), argType.Kind().String(), argRef.Interface())))
		}

		params[i] = argRef
	}

	results := methodRef.Call(params)
	if len(results) > 0 {
		if err, ok := results[len(results)-1].Interface().(error); ok && err != nil {
			return invokeErrorResult(clientId, encrypted, ctx, "Call error", err)
		}

		if _, ok := results[0].Interface().(error); !ok {
			return invokeSuccessJsonResult(clientId, encrypted, ctx, results[0].Interface())
		}
	}

	return invokeSuccessJsonResult(clientId, encrypted, ctx, map[string]string{})
}

func handshakeGetKey() map[string]string {
	return map[string]string{
		"publickey": rsaPublicKey,
	}
}

func handshakeSetKey(clientId string, clientKey []float64) error {
	priv, err := ezsec.RSABytesToPrivateKey([]byte(rsaPrivateKey))
	if err != nil {
		return err
	}

	convertedKey := make([]byte, len(clientKey))
	for i := 0; i < len(clientKey); i++ {
		convertedKey[i] = byte(clientKey[i])
	}

	key, err := ezsec.RSADecrypt(ezsec.ShaTypeSha512, convertedKey, priv)
	if err != nil {
		return err
	}

	clientKeys[clientId] = key

	return nil
}
