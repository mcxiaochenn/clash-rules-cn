package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// Result 表示一个拉取结果
type Result struct {
	Name    string
	Content []byte
	Error   error
}

// Fetcher 负责并发拉取上游数据
type Fetcher struct {
	client    *http.Client
	maxRetries int
}

// New 创建一个新的 Fetcher
func New(maxRetries int) *Fetcher {
	return &Fetcher{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		maxRetries: maxRetries,
	}
}

// FetchAll 并发拉取所有 URL，返回结果通道
func (f *Fetcher) FetchAll(urls map[string]string) <-chan Result {
	results := make(chan Result, len(urls))
	var wg sync.WaitGroup

	for name, url := range urls {
		wg.Add(1)
		go func(name, url string) {
			defer wg.Done()
			content, err := f.fetchWithRetry(url)
			results <- Result{
				Name:    name,
				Content: content,
				Error:   err,
			}
		}(name, url)
	}

	// 等待所有拉取完成后关闭通道
	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

// fetchWithRetry 带重试的 HTTP 拉取
func (f *Fetcher) fetchWithRetry(url string) ([]byte, error) {
	var lastErr error

	for i := 0; i <= f.maxRetries; i++ {
		if i > 0 {
			// 重试前等待，指数退避
			time.Sleep(time.Duration(i) * time.Second)
		}

		content, err := f.fetch(url)
		if err == nil {
			return content, nil
		}
		lastErr = err
	}

	return nil, fmt.Errorf("拉取失败（已重试 %d 次）: %w", f.maxRetries, lastErr)
}

// fetch 执行单次 HTTP 拉取
func (f *Fetcher) fetch(url string) ([]byte, error) {
	resp, err := f.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP 请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP 状态码异常: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	return body, nil
}
