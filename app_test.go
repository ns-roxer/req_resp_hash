package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"testing"
)

func Test_normalizeUrl(t *testing.T) {
	type args struct {
		u string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: `Has "http://" prefix - do nothing`,
			args: args{u: "http://google.com"},
			want: "http://google.com",
		},
		{
			name: `Has "https://" prefix - do nothing`,
			args: args{u: "https://google.com"},
			want: "https://google.com",
		},
		{
			name: `Has no scheme prefix - add "https://"`,
			args: args{u: "google.com"},
			want: "https://google.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizeUrl(tt.args.u); got != tt.want {
				t.Errorf("normalizeUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_responseHasherApp_CalcUrlHashes(t *testing.T) {
	type fields struct {
		parallelWorkersNum uint
		httpClient         SimpleHttpClient
		semaphore          chan struct{}
	}
	type args struct {
		urls []string
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		successResult []string
		wantErr       bool
	}{
		{
			name: "Happy path",
			fields: fields{
				parallelWorkersNum: 1,
				httpClient:         &httpClientStub{},
				semaphore:          make(chan struct{}, 1),
			},
			args: args{
				urls: []string{"https://google.com", "https://yandex.com", "https://reddit.com"},
			},
			successResult: []string{
				fmt.Sprintf("%s %x", "https://google.com", md5.Sum([]byte("response from https://google.com"))),
				fmt.Sprintf("%s %x", "https://yandex.com", md5.Sum([]byte("response from https://yandex.com"))),
				fmt.Sprintf("%s %x", "https://reddit.com", md5.Sum([]byte("response from https://reddit.com"))),
			},
			wantErr: false,
		},
		{
			name: "Happy path with multiple concurrent http requests",
			fields: fields{
				parallelWorkersNum: 5,
				httpClient:         &httpClientStub{},
				semaphore:          make(chan struct{}, 5),
			},
			args: args{
				urls: []string{"https://google.com", "https://yandex.com", "https://reddit.com"},
			},
			successResult: []string{
				fmt.Sprintf("%s %x", "https://google.com", md5.Sum([]byte("response from https://google.com"))),
				fmt.Sprintf("%s %x", "https://yandex.com", md5.Sum([]byte("response from https://yandex.com"))),
				fmt.Sprintf("%s %x", "https://reddit.com", md5.Sum([]byte("response from https://reddit.com"))),
			},
			wantErr: false,
		},
		{
			name: "Happy path with empty request",
			fields: fields{
				parallelWorkersNum: 5,
				httpClient:         &httpClientStub{},
				semaphore:          make(chan struct{}, 5),
			},
			args: args{
				urls: []string{},
			},
			successResult: nil,
			wantErr:       false,
		},
		{
			name: "Error occurred",
			fields: fields{
				parallelWorkersNum: 1,
				httpClient:         &httpClientStub{},
				semaphore:          make(chan struct{}, 1),
			},
			args:    args{urls: []string{"unexpected_url"}},
			wantErr: true,
		},
		{
			name: "Error occurred with multiple concurrent http requests",
			fields: fields{
				parallelWorkersNum: 5,
				httpClient:         &httpClientStub{},
				semaphore:          make(chan struct{}, 5),
			},
			args: args{
				urls: []string{"https://google.com", "https://yandex.com", "https://reddit.com", "unexpected_url"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &responseHasherApp{
				parallelWorkersNum: tt.fields.parallelWorkersNum,
				httpClient:         tt.fields.httpClient,
				semaphore:          tt.fields.semaphore,
			}
			got, err := app.CalcUrlHashes(tt.args.urls)

			if tt.wantErr && (err != nil) {
				return
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Test %q, for CalcUrlHashes() error = %v, wantErr %v", tt.name, err, tt.wantErr)
				return
			}
			if tt.successResult == nil && got == nil {
				return
			}
			for _, wantResponse := range tt.successResult {
				isPresented := false
				for _, u := range got {
					if wantResponse == u {
						isPresented = true
						break
					}
				}
				if !isPresented {
					t.Errorf("Test %q, for CalcUrlHashes() response %q is not presented", tt.name, wantResponse)
				}
			}
		})
	}
}

type httpClientStub struct {
}

func (s *httpClientStub) GetContentFromUrl(url string) ([]byte, error) {
	if url == "https://google.com" {
		return []byte("response from https://google.com"), nil
	}
	if url == "https://yandex.com" {
		return []byte("response from https://yandex.com"), nil
	}
	if url == "https://reddit.com" {
		return []byte("response from https://reddit.com"), nil
	}
	return nil, errors.New("unexpected error")
}
