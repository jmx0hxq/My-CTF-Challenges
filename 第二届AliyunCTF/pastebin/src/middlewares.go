package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

func needAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil || token.Value == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		username, err := verifyJWTToken(token.Value)
		if err != nil || username == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func adminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("token")
		if err != nil || token.Value == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		username, err := verifyJWTToken(token.Value)
		if err != nil || username != "admin" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func commonMiddleware(DBCon *sql.DB) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			mContext := GetCtx(ctx)

			reqContext := &ReqContext{
				Context: mContext,
				DBCon:   DBCon,
			}

			*r = *r.WithContext(context.WithValue(ctx, reqContextKey{}, reqContext))

			next.ServeHTTP(w, r)
		})
	}

}

func addCSPHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nonce, err := generateNonce()
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Security-Policy", fmt.Sprintf("default-src 'self'; style-src 'nonce-%s' ;script-src 'nonce-%s'; connect-src data:; img-src data:;", nonce, nonce))
		ctx := GetReqCtx(r.Context())
		ctx.Nonce = nonce
		next.ServeHTTP(w, r)
	})
}

func addCacheHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "-1")
		next.ServeHTTP(w, r)
	})
}

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		log.Printf("%s %s %v", r.Method, r.URL.Path, duration)
	})
}

type bufferResponseWriter struct {
	http.ResponseWriter
	buf bytes.Buffer
}

func (bw *bufferResponseWriter) Write(b []byte) (int, error) {
	return bw.buf.Write(b)
}

func secureFlagMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path, err := url.PathUnescape(r.URL.Path)
		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		if path == "/flag" {
			bw := &bufferResponseWriter{ResponseWriter: w}
			next.ServeHTTP(bw, r)
			log.Print(strings.Contains(bw.buf.String(), "admin"))

			var modifiedBody []byte
			if !strings.Contains(bw.buf.String(), "admin") {
				re := regexp.MustCompile(`aliyunctf{.*}`)
				modifiedBody = re.ReplaceAll(bw.buf.Bytes(), []byte(""))
				w.Write(modifiedBody)
			} else {
				w.Write(bw.buf.Bytes())
			}

		} else {
			next.ServeHTTP(w, r)
		}

	})
}

// Handle with the gzip request
type GzipResponseWriter struct {
	w *gzip.Writer
	http.ResponseWriter
}

func (w *GzipResponseWriter) WriteHeader(c int) {
	w.Header().Del("Content-Length")
	w.ResponseWriter.WriteHeader(c)
}

func (w GzipResponseWriter) Write(p []byte) (int, error) {
	if w.Header().Get("Content-Type") == "" {
		w.Header().Set("Content-Type", http.DetectContentType(p))
	}
	w.Header().Del("Content-Length")
	return w.w.Write(p)
}

func gzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		defer gz.Close()

		gzipWriter := &GzipResponseWriter{gz, w}
		gzipWriter.Header().Set("Content-Encoding", "gzip")
		gzipWriter.Header().Set("Vary", "Accept-Encoding")
		next.ServeHTTP(gzipWriter, r)
		_ = gzipWriter.w.Close()
	})
}
