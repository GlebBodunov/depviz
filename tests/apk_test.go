package tests

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"reflect"
	"strings"
	"testing"

	"depviz/src/apk"
	"depviz/src/graph"
)

const sampleIndex = `C:Q1abc=
P:musl
V:1.2.5-r0
p:so:libc.musl-x86_64.so.1=1

P:busybox
V:1.36.1-r29
D:so:libc.musl-x86_64.so.1 musl>=1.2 !busybox-extras
p:cmd:sh=1.36.1-r29

P:curl
V:8.9.0-r0
D:ca-certificates so:libc.musl-x86_64.so.1 libcurl=8.9.0-r0

P:libcurl
V:8.9.0-r0
`

func TestDirectDeps(t *testing.T) {
	idx, err := apk.ParseIndex(strings.NewReader(sampleIndex))
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	cases := map[string][]string{
		"busybox": {"musl"},
		"curl":    {"ca-certificates", "musl", "libcurl"},
		"musl":    nil,
	}
	for name, want := range cases {
		got, err := idx.DirectDeps(name)
		if err != nil {
			t.Fatalf("%s: неожиданная ошибка: %v", name, err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s: получено %v, ожидалось %v", name, got, want)
		}
	}
}

func TestUnknownPackage(t *testing.T) {
	idx, _ := apk.ParseIndex(strings.NewReader(sampleIndex))
	if _, err := idx.DirectDeps("nonexistent"); err == nil {
		t.Error("ожидалась ошибка для несуществующего пакета")
	}
}

func TestReadArchive(t *testing.T) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	content := []byte(sampleIndex)
	header := &tar.Header{Name: "APKINDEX", Mode: 0o644, Size: int64(len(content))}
	if err := tw.WriteHeader(header); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write(content); err != nil {
		t.Fatal(err)
	}
	tw.Close()
	gz.Close()

	idx, err := apk.ReadArchive(&buf)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if _, ok := idx.Packages["curl"]; !ok {
		t.Error("пакет curl не найден после распаковки")
	}
}

func TestReadArchiveInvalid(t *testing.T) {
	if _, err := apk.ReadArchive(strings.NewReader("not an archive")); err == nil {
		t.Error("ожидалась ошибка для некорректного архива")
	}
}

func TestReverseApk(t *testing.T) {
	idx, _ := apk.ParseIndex(strings.NewReader(sampleIndex))
	g, err := graph.Build(graph.Reverse(idx), "musl", unlimitedDepth)
	if err != nil {
		t.Fatalf("неожиданная ошибка: %v", err)
	}
	if !reflect.DeepEqual(g.Reached(), []string{"busybox", "curl"}) {
		t.Errorf("от musl зависят %v, ожидалось [busybox curl]", g.Reached())
	}
}
