package registry

import (
	"crypto/sha256"
	_ "embed"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

const digestHeader = "docker-content-digest"

var (
	//go:embed manifest.json
	manifest       []byte
	ManifestDigest string
	//go:embed blob.tar.gz
	blob       []byte
	BlobDigest string
	//go:embed config.json
	config       []byte
	ConfigDigest string
)

func digest(data []byte) string {
	return fmt.Sprintf("sha256:%x", sha256.Sum256(data))
}

func init() {
	ManifestDigest = digest(manifest)
	BlobDigest = digest(blob)
	ConfigDigest = digest(config)
}

func manifestHandler(rw http.ResponseWriter, req *http.Request) {
	respDigest := ManifestDigest
	if reqDigest := req.PathValue("digest"); strings.HasPrefix(reqDigest, "sha256:") {
		respDigest = reqDigest
	}
	rw.Header().Set(digestHeader, respDigest)
	rw.Header().Set("Content-Type", "application/vnd.docker.distribution.manifest.v2+json")
	rw.Write(manifest)
}

func blobHandler(rw http.ResponseWriter, req *http.Request) {
	switch req.PathValue("digest") {
	case ConfigDigest:
		rw.Header().Set(digestHeader, ConfigDigest)
		rw.Write(config)
	case BlobDigest:
		rw.Header().Set(digestHeader, BlobDigest)
		rw.Write(blob)
	default:
		rw.WriteHeader(http.StatusNotFound)
	}
}

func v2Handler(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Docker-Distribution-API-Version", "registry/2.0")
	rw.WriteHeader(http.StatusOK)
}

func Run(addr *string) {
	http.DefaultServeMux.Handle("/v2/busybox/manifests/{digest}", http.HandlerFunc(manifestHandler))
	http.DefaultServeMux.Handle("/v2/busybox/blobs/{digest}", http.HandlerFunc(blobHandler))
	http.DefaultServeMux.Handle("/v2/", http.HandlerFunc(v2Handler))

	listener, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("Could not listen on %s: %v", *addr, err)
	}
	log.Printf("serving a registry on %s ...", *addr)
	http.Serve(listener, http.DefaultServeMux)
}
