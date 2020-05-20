package main

import (
	"fmt"
	"github.com/alexflint/go-arg"
	"github.com/azon0320/jwtcli/signer"
	"github.com/azon0320/jwtcli/utils"
	"github.com/pascaldekloe/jwt"
	"os"
	"regexp"
	"strconv"
	"time"
)

type argStruct struct {
	Alg       string   `arg:"-A,required" help:"specify jwt sign alg"`
	ClaimSet  []string `arg:"-C,separate" help:"add claims kv"`
	Options   []string `arg:"-O,separate" help:"add extra options"`
	Aud       []string `arg:"-a,separate" help:"add audiences"`
	Expires   string   `arg:"env,--exp" help:"set expires at, e.g. +15<s|m|h|d|w|M>"`
	NotBefore string   `arg:"env,--nbf" help:"set not before at, e.g. +15<s|m|h|d|w|M>"`
	Issued    string   `arg:"env,--iat" help:"set issued at, e.g. -10<s|m|h|d|w|M>"`
	Issuer    string   `arg:"env,--iss" help:"set issuer"`
	Subject   string   `arg:"env,--sub" help:"set subject"`
	ID        string   `arg:"env,--jti" help:"set jwt id"`
}

func WithTime(timeStr string) time.Time {
	var now = time.Now()
	if len(timeStr) < 3 {
		return now
	}
	// e.g.  +15s
	var timeOperation = timeStr[:1]            // +
	var timeUnit = timeStr[len(timeStr)-1:]    // s
	var timeBody = timeStr[1 : len(timeStr)-1] // 15
	var timeBodyI, _ = strconv.Atoi(timeBody)  // if err, result = 0
	var timeBodyD time.Duration
	if timeOperation == "+" {
		timeBodyD = time.Duration(timeBodyI)
	} else {
		timeBodyD = time.Duration(-timeBodyI)
	}
	switch timeUnit {
	case "s":
		now = now.Add(timeBodyD * time.Second)
	case "m":
		now = now.Add(timeBodyD * time.Minute)
	case "h":
		now = now.Add(timeBodyD * time.Hour)
	case "d":
		now = now.Add(timeBodyD * time.Hour * 24)
	case "w":
		now = now.Add(timeBodyD * time.Hour * 24 * 7)
	case "M":
		now = now.Add(timeBodyD * time.Hour * 24 * 7 * 30)
	}
	return now
}

func parseLinesToMap(lines []string) map[string]interface{} {
	var lineMap = make(map[string]interface{}, 0)
	var equalsRegex = regexp.MustCompile("^([^=]*)=([^=]*)$")
	for _, line := range lines {
		var kvArgs = equalsRegex.FindStringSubmatch(line)
		if len(kvArgs) == 3 { // 1: k=v, 2:k, 3:v
			var key, value = kvArgs[1], kvArgs[2]
			var finalValue interface{}
			finalValue = value
			if len(value) > 0 {
				var suffix = value[len(value)-1:]
				switch suffix {
				case "i": // integer
					var num, err = strconv.Atoi(value[:len(value)-1])
					if err == nil {
						finalValue = num
					}
				case "f": // float64
					var flt, err = strconv.ParseFloat(value[:len(value)-1], 64)
					if err == nil {
						finalValue = flt
					}
				case "b": // boolean
					switch value {
					case "trueb":
						finalValue = true
					case "falseb":
						finalValue = false
					}
				}
			}
			lineMap[key] = finalValue
		}
	}
	return lineMap
}

func main() {
	var argStruct argStruct
	arg.MustParse(&argStruct)
	var parsedStruct struct {
		ClaimsSet *utils.Wrapper
		Options   *utils.Wrapper
	}
	parsedStruct.ClaimsSet = utils.WrapMap(parseLinesToMap(argStruct.ClaimSet))
	parsedStruct.Options = utils.WrapMap(parseLinesToMap(argStruct.Options))
	var claims = &jwt.Claims{}
	claims.Set = parsedStruct.ClaimsSet.Map
	claims.ID = argStruct.ID
	claims.Issuer = argStruct.Issuer
	claims.Subject = argStruct.Subject

	{
		// https://github.com/pascaldekloe/jwt/issues/12
		// https://github.com/dgrijalva/jwt-go/issues/184
		// https://tools.ietf.org/html/rfc7519#section-4.1
		if len(argStruct.Aud) == 1 {
			claims.Set["aud"] = argStruct.Aud[0]
		}else {
			claims.Audiences = argStruct.Aud
		}
	}

	if len(argStruct.Expires) > 0 {
		claims.Expires = jwt.NewNumericTime(WithTime(argStruct.Expires).Round(time.Second))
	}
	if len(argStruct.NotBefore) > 0 {
		claims.NotBefore = jwt.NewNumericTime(WithTime(argStruct.NotBefore).Round(time.Second))
	}
	claims.Issued = jwt.NewNumericTime(WithTime(argStruct.Issued).Round(time.Second))
	claims.ID = argStruct.ID
	token, err := signer.MakeJwt(argStruct.Alg, claims, parsedStruct.Options)
	if err != nil {
		fmt.Printf("error generating jwt: %s", err.Error())
		os.Exit(-1)
	} else {
		fmt.Println(token)
		os.Exit(0)
	}
}
