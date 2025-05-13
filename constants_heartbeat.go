package roselite

import (
	"net/url"
	"strconv"
	"time"

	"github.com/guregu/null/v6"
)

type Heartbeat struct {
	Status            HeartbeatStatus `json:"status"`
	Latency           int64           `json:"latency"`
	AdditionalMessage null.String     `json:"additional_message,omitempty"`
	HttpProtocol      null.String     `json:"http_protocol,omitempty"`
	TLSVersion        null.String     `json:"tls_version,omitempty"`
	TLSCipherName     null.String     `json:"tls_cipher_name,omitempty"`
	TLSExpiryDate     null.Time       `json:"tls_expiry_date,omitempty"`
}

func HeartbeatFromQuery(query url.Values) Heartbeat {
	status := HeartbeatStatusFromString(query.Get("status"))
	latency, _ := strconv.ParseInt(query.Get("ping"), 10, 64)
	message := query.Get("msg")
	httpProtocol := query.Get("http_protocol")
	tlsVersion := query.Get("tls_version")
	tlsCipherName := query.Get("tls_cipher")
	tlsExpiry := query.Get("tls_expiry")
	parsedTlsExpiryDate, _ := strconv.ParseInt(tlsExpiry, 10, 64)
	tlsExpiryDate := time.Unix(parsedTlsExpiryDate, 0)

	return Heartbeat{
		Status:            status,
		Latency:           latency,
		AdditionalMessage: null.NewString(message, message != ""),
		HttpProtocol:      null.NewString(httpProtocol, httpProtocol != ""),
		TLSVersion:        null.NewString(tlsVersion, tlsVersion != ""),
		TLSCipherName:     null.NewString(tlsCipherName, tlsCipherName != ""),
		TLSExpiryDate:     null.NewTime(tlsExpiryDate, !tlsExpiryDate.IsZero()),
	}
}

func (h Heartbeat) ToQuery() url.Values {
	query := url.Values{}
	query.Set("status", h.Status.String())
	query.Set("ping", strconv.FormatInt(h.Latency, 10))
	if h.AdditionalMessage.Valid {
		query.Set("message", h.AdditionalMessage.ValueOrZero())
	}
	if h.HttpProtocol.Valid {
		query.Set("http_protocol", h.HttpProtocol.ValueOrZero())
	}
	if h.TLSVersion.Valid {
		query.Set("tls_version", h.TLSVersion.ValueOrZero())
	}
	if h.TLSCipherName.Valid {
		query.Set("tls_cipher", h.TLSCipherName.ValueOrZero())
	}
	if h.TLSExpiryDate.Valid {
		query.Set("tls_expiry", strconv.FormatInt(h.TLSExpiryDate.Time.Unix(), 10))
	}

	return query
}
