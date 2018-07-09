/* ワードカウントをするプログラム*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

func main() {

	var (
		input []string
		mutex = &sync.Mutex{}
		wg    sync.WaitGroup
	)

	// データ入力
	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		// 何も入力しないでリターンされたら入力処理を終える。
		if stdin.Text() == "" {
			break
		}
		input = append(input, stdin.Text())
	}

	m := make([]map[string]int, len(input))
	// map処理
	// ワードカウント
	for i, a := range input {
		wg.Add(1) // 同期カウンタを上げる
		// 並列実行（ゴルーチン）
		go func(m []map[string]int, a string, i int) {
			txt := strings.Fields(a)
			m[i] = map[string]int{}
			for _, b := range txt {
				mutex.Lock() // 排他制御
				m[i][b]++
				mutex.Unlock() // ロック解除
			}
			wg.Done() // 同期カウンタを下げる
		}(m, a, i)
	}
	wg.Wait() // 同期カウンタが0になるまで待つ

	s := make(map[string][]int)
	// シャッフル
	// 同じキーのものをまとめる。
	for i := range m {
		for j, b := range m[i] {
			// 並列実行（ゴルーチン）
			wg.Add(1)
			go func(s map[string][]int, j string, b int) {
				s[j] = append(s[j], b)
				wg.Done()
			}(s, j, b)
		}
	}
	wg.Wait()

	// reduce処理
	// Mapから意味のある出力をする。
	for i, a := range s {
		// 並列実行（ゴルーチン）
		wg.Add(1)
		go func(i string, a []int) {
			var count int
			for _, b := range a {
				count += b
			}
			fmt.Println(i, count)
			count = 0
			wg.Done()
		}(i, a)
	}
	wg.Wait()
}
