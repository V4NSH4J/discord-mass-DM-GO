package client

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"strconv"
	"strings"

	utls "github.com/Danny-Dasilva/utls"
	"github.com/andybalholm/brotli"
)

const (
	chrome  = "chrome"  //chrome User agent enum
	firefox = "firefox" //firefox User agent enum
)

func parseUserAgent(userAgent string) string {
	switch {
	case strings.Contains(strings.ToLower(userAgent), "chrome"):
		return chrome
	case strings.Contains(strings.ToLower(userAgent), "firefox"):
		return firefox
	default:
		return chrome
	}

}

// DecompressBody unzips compressed data
func DecompressBody(Body []byte, encoding []string, content []string) (parsedBody string) {
	if len(encoding) > 0 {
		if encoding[0] == "gzip" {
			unz, err := gUnzipData(Body)
			if err != nil {
				return string(Body)
			}
			return string(unz)
		} else if encoding[0] == "deflate" {
			unz, err := enflateData(Body)
			if err != nil {
				return string(Body)
			}
			return string(unz)
		} else if encoding[0] == "br" {
			unz, err := unBrotliData(Body)
			if err != nil {
				return string(Body)
			}
			return string(unz)
		}
	} else if len(content) > 0 {
		decodingTypes := map[string]bool{
			"image/svg+xml":   true,
			"image/webp":      true,
			"image/jpeg":      true,
			"image/png":       true,
			"application/pdf": true,
		}
		if decodingTypes[content[0]] {
			return base64.StdEncoding.EncodeToString(Body)
		}
	}
	parsedBody = string(Body)
	return parsedBody

}

func gUnzipData(data []byte) (resData []byte, err error) {
	gz, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return []byte{}, err
	}
	defer gz.Close()
	respBody, err := ioutil.ReadAll(gz)
	return respBody, err
}
func enflateData(data []byte) (resData []byte, err error) {
	zr, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return []byte{}, err
	}
	defer zr.Close()
	enflated, err := ioutil.ReadAll(zr)
	return enflated, err
}
func unBrotliData(data []byte) (resData []byte, err error) {
	br := brotli.NewReader(bytes.NewReader(data))
	respBody, err := ioutil.ReadAll(br)
	return respBody, err
}

// StringToSpec creates a ClientHelloSpec based on a JA3 string
func StringToSpec(ja3 string, userAgent string) (*utls.ClientHelloSpec, error) {
	parsedUserAgent := parseUserAgent(userAgent)
	extMap := genMap()
	tokens := strings.Split(ja3, ",")

	version := tokens[0]
	ciphers := strings.Split(tokens[1], "-")
	extensions := strings.Split(tokens[2], "-")
	curves := strings.Split(tokens[3], "-")
	if len(curves) == 1 && curves[0] == "" {
		curves = []string{}
	}
	pointFormats := strings.Split(tokens[4], "-")
	if len(pointFormats) == 1 && pointFormats[0] == "" {
		pointFormats = []string{}
	}
	// parse curves
	var targetCurves []utls.CurveID
	targetCurves = append(targetCurves, utls.CurveID(utls.GREASE_PLACEHOLDER)) //append grease for Chrome browsers
	for _, c := range curves {
		cid, err := strconv.ParseUint(c, 10, 16)
		if err != nil {
			return nil, err
		}
		targetCurves = append(targetCurves, utls.CurveID(cid))
		// if cid != uint64(utls.CurveP521) {
		// CurveP521 sometimes causes handshake errors
		// }
	}
	extMap["10"] = &utls.SupportedCurvesExtension{Curves: targetCurves}

	// parse point formats
	var targetPointFormats []byte
	for _, p := range pointFormats {
		pid, err := strconv.ParseUint(p, 10, 8)
		if err != nil {
			return nil, err
		}
		targetPointFormats = append(targetPointFormats, byte(pid))
	}
	extMap["11"] = &utls.SupportedPointsExtension{SupportedPoints: targetPointFormats}

	// set extension 43
	vid64, err := strconv.ParseUint(version, 10, 16)
	if err != nil {
		return nil, err
	}
	vid := uint16(vid64)
	// extMap["43"] = &utls.SupportedVersionsExtension{
	// 	Versions: []uint16{
	// 		utls.VersionTLS12,
	// 	},
	// }

	// build extenions list
	var exts []utls.TLSExtension
	//Optionally Add Chrome Grease Extension
	if parsedUserAgent == chrome {
		exts = append(exts, &utls.UtlsGREASEExtension{})
	}
	for _, e := range extensions {
		te, ok := extMap[e]
		if !ok {
			// eAsint, err := strconv.Atoi(e)
			// if err != nil {
			// 	return nil, err
			// }
			// te = &utls.GenericExtension{Id: uint16(eAsint)}
			continue
			// return nil, raiseExtensionError(e)
		}
		// //Optionally add Chrome Grease Extension
		if e == "21" && parsedUserAgent == chrome {
			exts = append(exts, &utls.UtlsGREASEExtension{})
		}
		exts = append(exts, te)
	}
	//Add this back in if user agent is chrome and no padding extension is given
	// if parsedUserAgent == chrome {
	// 	exts = append(exts, &utls.UtlsGREASEExtension{})
	// 	exts = append(exts, &utls.UtlsPaddingExtension{GetPaddingLen: utls.BoringPaddingStyle})
	// }
	// build SSLVersion
	// vid64, err := strconv.ParseUint(version, 10, 16)
	// if err != nil {
	// 	return nil, err
	// }

	// build CipherSuites
	var suites []uint16
	//Optionally Add Chrome Grease Extension
	if parsedUserAgent == chrome {
		suites = append(suites, utls.GREASE_PLACEHOLDER)
	}
	for _, c := range ciphers {
		cid, err := strconv.ParseUint(c, 10, 16)
		if err != nil {
			return nil, err
		}
		suites = append(suites, uint16(cid))
	}
	_ = vid
	return &utls.ClientHelloSpec{
		// TLSVersMin:         vid,
		// TLSVersMax:         vid,
		CipherSuites:       suites,
		CompressionMethods: []byte{0},
		Extensions:         exts,
		GetSessionID:       sha256.Sum256,
	}, nil
}

func genMap() (extMap map[string]utls.TLSExtension) {
	extMap = map[string]utls.TLSExtension{
		"0": &utls.SNIExtension{},
		"5": &utls.StatusRequestExtension{},
		// These are applied later
		// "10": &tls.SupportedCurvesExtension{...}
		// "11": &tls.SupportedPointsExtension{...}
		"13": &utls.SignatureAlgorithmsExtension{
			SupportedSignatureAlgorithms: []utls.SignatureScheme{
				utls.ECDSAWithP256AndSHA256,
				utls.ECDSAWithP384AndSHA384,
				utls.ECDSAWithP521AndSHA512,
				utls.PSSWithSHA256,
				utls.PSSWithSHA384,
				utls.PSSWithSHA512,
				utls.PKCS1WithSHA256,
				utls.PKCS1WithSHA384,
				utls.PKCS1WithSHA512,
				utls.ECDSAWithSHA1,
				utls.PKCS1WithSHA1,
			},
		},
		"16": &utls.ALPNExtension{
			AlpnProtocols: []string{"h2", "http/1.1"},
		},
		"17": &utls.GenericExtension{Id: 17}, // status_request_v2
		"18": &utls.SCTExtension{},
		"21": &utls.UtlsPaddingExtension{GetPaddingLen: utls.BoringPaddingStyle},
		"22": &utls.GenericExtension{Id: 22}, // encrypt_then_mac
		"23": &utls.UtlsExtendedMasterSecretExtension{},
		"27": &utls.CompressCertificateExtension{
			Algorithms: []utls.CertCompressionAlgo{utls.CertCompressionBrotli},
		},
		"28": &utls.FakeRecordSizeLimitExtension{}, //Limit: 0x4001
		"35": &utls.SessionTicketExtension{},
		"34": &utls.GenericExtension{Id: 34},
		"41": &utls.GenericExtension{Id: 41}, //FIXME pre_shared_key
		"43": &utls.SupportedVersionsExtension{Versions: []uint16{
			utls.GREASE_PLACEHOLDER,
			utls.VersionTLS13,
			utls.VersionTLS12,
			utls.VersionTLS11,
			utls.VersionTLS10}},
		"44": &utls.CookieExtension{},
		"45": &utls.PSKKeyExchangeModesExtension{Modes: []uint8{
			utls.PskModeDHE,
		}},
		"49": &utls.GenericExtension{Id: 49}, // post_handshake_auth
		"50": &utls.GenericExtension{Id: 50}, // signature_algorithms_cert
		"51": &utls.KeyShareExtension{KeyShares: []utls.KeyShare{
			{Group: utls.CurveID(utls.GREASE_PLACEHOLDER), Data: []byte{0}},
			{Group: utls.X25519},

			// {Group: utls.CurveP384}, known bug missing correct extensions for handshake
		}},
		"30032": &utls.GenericExtension{Id: 0x7550, Data: []byte{0}}, //FIXME
		"13172": &utls.NPNExtension{},
		"17513": &utls.ApplicationSettingsExtension{
			SupportedALPNList: []string{
				"h2",
			},
		},
		"65281": &utls.RenegotiationInfoExtension{
			Renegotiation: utls.RenegotiateOnceAsClient,
		},
	}
	return

}

func PrettyStruct(data interface{}) (string, error) {
	val, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}
	return string(val), nil
}
