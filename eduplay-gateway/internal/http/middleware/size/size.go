package middleware

import (
	// "errors"
	// "io"
	// "mime/multipart"
	"net/http"
)

func MaxBodySize(limit int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, limit)
			next.ServeHTTP(w, r)
		})
	}
}

// func LimitMultipartParts(maxPartSize int64) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			if r.Header.Get("Content-Type") == "" {
// 				next.ServeHTTP(w, r)
// 				return
// 			}

// 			mr, err := r.MultipartReader()
// 			if err != nil {
// 				next.ServeHTTP(w, r)
// 				return
// 			}

// 			pr, pw := io.Pipe()

// 			go func() {
// 				defer func() {
// 					if err := pw.Close(); err != nil {
// 						panic(err)
// 					}
// 				}()

// 				writer := multipart.NewWriter(pw)
// 				defer func() {
// 					if err := writer.Close(); err != nil {
// 						panic(err)
// 					}
// 				}()

// 				for {
// 					part, err := mr.NextPart()
// 					if err == io.EOF {
// 						return
// 					}
// 					if err != nil {
// 						pw.CloseWithError(err)
// 						return
// 					}

// 					newPart, err := writer.CreatePart(part.Header)
// 					if err != nil {
// 						pw.CloseWithError(err)
// 						return
// 					}

// 					limited := io.LimitReader(part, maxPartSize+1)

// 					n, err := io.Copy(newPart, limited)
// 					if err != nil {
// 						pw.CloseWithError(err)
// 						return
// 					}

// 					if n > maxPartSize {
// 						pw.CloseWithError(errors.New("multipart part too large"))
// 						return
// 					}
// 				}
// 			}()

// 			r.Body = pr
// 			r.Header.Set("Content-Type", writerContentType(r))

// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }

// func writerContentType(r *http.Request) string {
// 	return r.Header.Get("Content-Type")
// }
