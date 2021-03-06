package vcsclient

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"testing"
)

func TestRepository_FileSystem_Open(t *testing.T) {
	setup()
	defer teardown()

	repoPath := "a.b/c"
	repo_, _ := vcsclient.Repository(repoPath)
	repo := repo_.(*repository)
	want := []byte("c")
	entry := &TreeEntry{
		Contents: want,
	}

	var called bool
	mux.HandleFunc(urlPath(t, RouteRepoTreeEntry, repo, map[string]string{"CommitID": "abcd", "Path": "f"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, entry)
	})

	fs, err := repo.FileSystem("abcd")
	if err != nil {
		t.Errorf("Repository.FileSystem returned error: %v", err)
		return
	}

	f, err := fs.Open("f")
	if err != nil {
		t.Errorf("FileSystem.Open returned error: %v", err)
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		t.Fatal(err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !bytes.Equal(data, want) {
		t.Errorf("FileSystem.Open returned data %+v, want %+v", data, want)
	}
}

func TestRepository_FileSystem_Lstat(t *testing.T) {
	setup()
	defer teardown()

	repoPath := "a.b/c"
	repo_, _ := vcsclient.Repository(repoPath)
	repo := repo_.(*repository)
	entry := &TreeEntry{Name: "f"}
	want, _ := entry.Stat()

	var called bool
	mux.HandleFunc(urlPath(t, RouteRepoTreeEntry, repo, map[string]string{"CommitID": "abcd", "Path": "f"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, entry)
	})

	fs, err := repo.FileSystem("abcd")
	if err != nil {
		t.Errorf("Repository.FileSystem returned error: %v", err)
		return
	}

	fi, err := fs.Lstat("f")
	if err != nil {
		t.Errorf("FileSystem.Lstat returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(fi, want) {
		t.Errorf("FileSystem.Lstat returned %+v, want %+v", fi, want)
	}
}

func TestRepository_FileSystem_Stat(t *testing.T) {
	setup()
	defer teardown()

	repoPath := "a.b/c"
	repo_, _ := vcsclient.Repository(repoPath)
	repo := repo_.(*repository)
	entry := &TreeEntry{Name: "f"}
	want, _ := entry.Stat()

	var called bool
	mux.HandleFunc(urlPath(t, RouteRepoTreeEntry, repo, map[string]string{"CommitID": "abcd", "Path": "f"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, entry)
	})

	fs, err := repo.FileSystem("abcd")
	if err != nil {
		t.Errorf("Repository.FileSystem returned error: %v", err)
		return
	}

	fi, err := fs.Stat("f")
	if err != nil {
		t.Errorf("FileSystem.Stat returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(fi, want) {
		t.Errorf("FileSystem.Stat returned %+v, want %+v", fi, want)
	}
}

func TestRepository_FileSystem_ReadDir(t *testing.T) {
	setup()
	defer teardown()

	repoPath := "a.b/c"
	repo_, _ := vcsclient.Repository(repoPath)
	repo := repo_.(*repository)
	entries := []*TreeEntry{{Name: "d/a"}, {Name: "d/b"}}
	fi0, _ := entries[0].Stat()
	fi1, _ := entries[1].Stat()
	want := []os.FileInfo{fi0, fi1}

	var called bool
	mux.HandleFunc(urlPath(t, RouteRepoTreeEntry, repo, map[string]string{"CommitID": "abcd", "Path": "d"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, &TreeEntry{Name: "d", Entries: entries})
	})

	fs, err := repo.FileSystem("abcd")
	if err != nil {
		t.Errorf("Repository.FileSystem returned error: %v", err)
		return
	}

	fis, err := fs.ReadDir("d")
	if err != nil {
		t.Errorf("FileSystem.ReadDir returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(fis, want) {
		t.Errorf("FileSystem.ReadDir returned %+v, want %+v", fis, want)
	}
}

func TestRepository_FileSystem_Get(t *testing.T) {
	setup()
	defer teardown()

	repoPath := "a.b/c"
	repo_, _ := vcsclient.Repository(repoPath)
	repo := repo_.(*repository)
	want := &TreeEntry{Name: "f", Contents: []byte("c")}

	var called bool
	mux.HandleFunc(urlPath(t, RouteRepoTreeEntry, repo, map[string]string{"CommitID": "abcd", "Path": "f"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")

		writeJSON(w, want)
	})

	fs, err := repo.FileSystem("abcd")
	if err != nil {
		t.Errorf("Repository.FileSystem returned error: %v", err)
		return
	}

	e, err := fs.(*repositoryFS).Get("f")
	if err != nil {
		t.Errorf("FileSystem.Stat returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(e, want) {
		t.Errorf("FileSystem.Get returned %+v, want %+v", e, want)
	}
}

func TestRepository_FileSystem_GetFileWithOptions(t *testing.T) {
	setup()
	defer teardown()

	repoPath := "a.b/c"
	repo_, _ := vcsclient.Repository(repoPath)
	repo := repo_.(*repository)
	want := &FileWithRange{
		TreeEntry: &TreeEntry{Name: "f", Contents: []byte("c")},
		FileRange: FileRange{
			StartByte: 123, EndByte: 456,
			StartLine: 2, EndLine: 4,
		},
	}

	var called bool
	mux.HandleFunc(urlPath(t, RouteRepoTreeEntry, repo, map[string]string{"CommitID": "abcd", "Path": "f"}), func(w http.ResponseWriter, r *http.Request) {
		called = true
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"StartByte": "123",
			"EndByte":   "456",
		})

		writeJSON(w, want)
	})

	fs, err := repo.FileSystem("abcd")
	if err != nil {
		t.Errorf("Repository.FileSystem returned error: %v", err)
		return
	}

	e, err := fs.(*repositoryFS).GetFileWithOptions("f", GetFileOptions{FileRange: FileRange{StartByte: 123, EndByte: 456}})
	if err != nil {
		t.Errorf("FileSystem.Stat returned error: %v", err)
	}

	if !called {
		t.Fatal("!called")
	}

	if !reflect.DeepEqual(e, want) {
		t.Errorf("FileSystem.Get returned %+v, want %+v", e, want)
	}
}
