// Copyright 2023 Gravitational, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"net/http"
	nethttp "net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// TransportFormatter is a span formatter that may be provided to
// [otelhttp.WithSpanNameFormatter] to include the url path in the span
// names generated by an [otelhttp.Transport].
func TransportFormatter(_ string, r *nethttp.Request) string {
	return "HTTP " + r.Method + " " + r.URL.Path
}

// HandlerFormatter is a span formatter that may be provided to
// [otelhttp.WithSpanNameFormatter] to include the component and url path in the span
// names generated by [otelhttp.NewHandler].
func HandlerFormatter(operation string, r *nethttp.Request) string {
	return operation + " " + r.Method + " " + r.URL.Path
}

// NewTransport wraps the provided [nethttp.RoundTripper] with one
// that automatically adds spans for each http request.
//
// Note: special care has been taken to ensure that the returned
// [nethttp.RoundTripper] has a CloseIdleConnections method because
// the [otelhttp.Transport] does not implement it:
// https://github.com/open-telemetry/opentelemetry-go-contrib/issues/3543.
// Once the issue is resolved the wrapper may be discarded.
func NewTransport(rt nethttp.RoundTripper) nethttp.RoundTripper {
	return NewTransportWithInner(rt, rt)
}

// NewTransportWithInner wraps the provided [nethttp.RoundTripper] with one
// that automatically adds spans for each http request.
// The inner round tripper is used to close idle connections when
// rt.CloseIdleConnections isn't implemented for the rt provided.
//
// Note: special care has been taken to ensure that the returned
// [nethttp.RoundTripper] has a CloseIdleConnections method because
// the [otelhttp.Transport] does not implement it:
// https://github.com/open-telemetry/opentelemetry-go-contrib/issues/3543.
// Once the issue is resolved the wrapper may be discarded.
func NewTransportWithInner(rt nethttp.RoundTripper, inner http.RoundTripper) nethttp.RoundTripper {
	return &roundTripWrapper{
		RoundTripper: otelhttp.NewTransport(rt, otelhttp.WithSpanNameFormatter(TransportFormatter)),
		inner:        inner,
	}
}

type closeIdler interface {
	CloseIdleConnections()
}

type roundTripWrapper struct {
	nethttp.RoundTripper
	inner nethttp.RoundTripper
}

// Unwrap returns the inner round tripper.
func (w *roundTripWrapper) Unwrap() http.RoundTripper {
	return w.inner
}

// CloseIdleConnections ensures that the returned [nethttp.RoundTripper]
// has a CloseIdleConnections method. Since [otelhttp.Transport] does not implement
// this any [nethttp.Client.CloseIdleConnections] calls result in a noop instead
// of forwarding the request onto its wrapped [nethttp.RoundTripper].
//
// DELETE WHEN https://github.com/open-telemetry/opentelemetry-go-contrib/issues/3543
// has been resolved.
func (w *roundTripWrapper) CloseIdleConnections() {
	if c, ok := w.RoundTripper.(closeIdler); ok {
		c.CloseIdleConnections()
	} else if c, ok := w.inner.(closeIdler); ok {
		c.CloseIdleConnections()
	}
}
