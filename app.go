package main

import (
	"crypto/md5"
	"fmt"
	"strings"
	"sync"
)

type responseHasherApp struct {
	parallelWorkersNum uint
	httpClient         SimpleHttpClient
	semaphore          chan struct{}
}

type SimpleHttpClient interface {
	GetContentFromUrl(url string) ([]byte, error)
}

func NewApp(parallelWorkersNum uint, httpClient SimpleHttpClient) *responseHasherApp {
	return &responseHasherApp{
		parallelWorkersNum: parallelWorkersNum,
		httpClient:         httpClient,
		semaphore:          make(chan struct{}, parallelWorkersNum),
	}
}

func (app *responseHasherApp) CalcUrlHashes(urls []string) ([]string, error) {
	if len(urls) == 0 {
		return nil, nil
	}
	var err error
	processedUrls := make([]string, 0, len(urls))
	urlsWithHash := make(chan string)
	errors := make(chan error)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case u := <-urlsWithHash:
				processedUrls = append(processedUrls, u)
			case err = <-errors:
				return
			default:
				continue
			}
			if len(processedUrls) == len(urls) {
				return
			}
		}
	}()

	for _, u := range urls {
		wg.Add(1)
		app.semaphore <- struct{}{}
		go func(url string) {
			defer func() { <-app.semaphore }()
			wg.Done()
			content, err := app.httpClient.GetContentFromUrl(url)
			if err != nil {
				errors <- err
				return
			}
			hashSum := md5.Sum(content)
			urlsWithHash <- fmt.Sprintf("%s %x", url, hashSum)
		}(u)
	}
	wg.Wait()

	if err != nil {
		return nil, err
	}
	return processedUrls, nil
}

func normalizeUrl(u string) string {
	if strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://") {
		return u
	}

	return "https://" + u
}
