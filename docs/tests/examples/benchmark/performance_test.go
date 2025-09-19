package benchmark

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// BenchmarkPasswordHashing 基準測試：密碼雜湊效能
// 測試 bcrypt 密碼雜湊的效能表現
func BenchmarkPasswordHashing(b *testing.B) {
	password := "testPassword123"

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkPasswordHashingDifferentCosts 測試不同成本的密碼雜湊效能
func BenchmarkPasswordHashingDifferentCosts(b *testing.B) {
	password := "testPassword123"
	costs := []int{bcrypt.MinCost, bcrypt.DefaultCost, 14}

	for _, cost := range costs {
		b.Run(fmt.Sprintf("Cost_%d", cost), func(b *testing.B) {
			b.ResetTimer()
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				_, err := bcrypt.GenerateFromPassword([]byte(password), cost)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkPasswordVerification 基準測試：密碼驗證效能
func BenchmarkPasswordVerification(b *testing.B) {
	password := "testPassword123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
		if err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkStringOperations 基準測試：字串操作效能
func BenchmarkStringOperations(b *testing.B) {
	// 測試字串拼接
	b.Run("StringConcatenation", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := "user_" + fmt.Sprintf("%d", i) + "@example.com"
			_ = result
		}
	})

	// 測試 fmt.Sprintf
	b.Run("FmtSprintf", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := fmt.Sprintf("user_%d@example.com", i)
			_ = result
		}
	})
}

// BenchmarkMapOperations 基準測試：Map 操作效能
func BenchmarkMapOperations(b *testing.B) {
	// 測試 map 寫入
	b.Run("MapWrite", func(b *testing.B) {
		m := make(map[string]int)
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("key_%d", i)
			m[key] = i
		}
	})

	// 測試 map 讀取
	b.Run("MapRead", func(b *testing.B) {
		m := make(map[string]int)
		for i := 0; i < 1000; i++ {
			key := fmt.Sprintf("key_%d", i)
			m[key] = i
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("key_%d", i%1000)
			_ = m[key]
		}
	})
}

// BenchmarkConcurrentOperations 基準測試：並發操作效能
func BenchmarkConcurrentOperations(b *testing.B) {
	// 測試並發 map 讀取（使用 sync.Map）
	b.Run("SyncMapConcurrentRead", func(b *testing.B) {
		var sm sync.Map

		// 預填資料
		for i := 0; i < 1000; i++ {
			key := fmt.Sprintf("key_%d", i)
			sm.Store(key, i)
		}

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				key := fmt.Sprintf("key_%d", randInt(1000))
				_, _ = sm.Load(key)
			}
		})
	})

	// 測試並發 map 寫入（使用 sync.Map）
	b.Run("SyncMapConcurrentWrite", func(b *testing.B) {
		var sm sync.Map

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			i := 0
			for pb.Next() {
				key := fmt.Sprintf("key_%d", i)
				sm.Store(key, i)
				i++
			}
		})
	})

	// 測試並發計數器
	b.Run("ConcurrentCounter", func(b *testing.B) {
		var counter int64
		var mu sync.Mutex

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				mu.Lock()
				counter++
				mu.Unlock()
			}
		})
	})
}

// BenchmarkMemoryAllocation 基準測試：記憶體分配效能
func BenchmarkMemoryAllocation(b *testing.B) {
	// 測試 slice 分配
	b.Run("SliceAllocation", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			slice := make([]int, 100)
			_ = slice
		}
	})

	// 測試 slice 預分配
	b.Run("SlicePreallocation", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			slice := make([]int, 0, 100)
			for j := 0; j < 100; j++ {
				slice = append(slice, j)
			}
		}
	})

	// 測試 map 分配
	b.Run("MapAllocation", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			m := make(map[string]int, 100)
			for j := 0; j < 100; j++ {
				key := fmt.Sprintf("key_%d", j)
				m[key] = j
			}
		}
	})
}

// BenchmarkJSONOperations 基準測試：JSON 操作效能
func BenchmarkJSONOperations(b *testing.B) {
	type User struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
		NickName string `json:"nickName"`
		Phone    string `json:"phone"`
		Status   string `json:"status"`
		RoleID   int    `json:"roleId"`
	}

	user := User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
		NickName: "測試使用者",
		Phone:    "0912345678",
		Status:   "2",
		RoleID:   1,
	}

	// 測試 JSON 序列化
	b.Run("JSONMarshal", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			_, err := json.Marshal(user)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	// 測試 JSON 反序列化
	jsonData, _ := json.Marshal(user)
	b.Run("JSONUnmarshal", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			var u User
			err := json.Unmarshal(jsonData, &u)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkDatabaseOperations 基準測試：模擬資料庫操作效能
func BenchmarkDatabaseOperations(b *testing.B) {
	// 模擬資料庫連接池
	type DBPool struct {
		connections chan *Connection
		maxConns    int
	}

	type Connection struct {
		ID   int
		Used bool
	}

	newDBPool := func(maxConns int) *DBPool {
		pool := &DBPool{
			connections: make(chan *Connection, maxConns),
			maxConns:    maxConns,
		}

		for i := 0; i < maxConns; i++ {
			pool.connections <- &Connection{ID: i}
		}

		return pool
	}

	getConnection := func(pool *DBPool) *Connection {
		return <-pool.connections
	}

	releaseConnection := func(pool *DBPool, conn *Connection) {
		conn.Used = false
		pool.connections <- conn
	}

	// 測試連接池效能
	b.Run("ConnectionPool", func(b *testing.B) {
		pool := newDBPool(10)
		b.ResetTimer()

		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				conn := getConnection(pool)
				// 模擬資料庫操作
				time.Sleep(time.Microsecond)
				releaseConnection(pool, conn)
			}
		})
	})

	// 測試查詢操作（模擬）
	b.Run("QueryOperation", func(b *testing.B) {
		data := make(map[int]string)
		for i := 0; i < 10000; i++ {
			data[i] = fmt.Sprintf("user_%d", i)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			id := i % 10000
			_ = data[id]
		}
	})
}

// BenchmarkCacheOperations 基準測試：快取操作效能
func BenchmarkCacheOperations(b *testing.B) {
	// 簡單的記憶體快取實作
	type MemoryCache struct {
		data map[string]interface{}
		mu   sync.RWMutex
	}

	newMemoryCache := func() *MemoryCache {
		return &MemoryCache{
			data: make(map[string]interface{}),
		}
	}

	set := func(cache *MemoryCache, key string, value interface{}) {
		cache.mu.Lock()
		defer cache.mu.Unlock()
		cache.data[key] = value
	}

	get := func(cache *MemoryCache, key string) (interface{}, bool) {
		cache.mu.RLock()
		defer cache.mu.RUnlock()
		value, exists := cache.data[key]
		return value, exists
	}

	// 測試快取寫入
	b.Run("CacheWrite", func(b *testing.B) {
		cache := newMemoryCache()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("key_%d", i)
			set(cache, key, i)
		}
	})

	// 測試快取讀取
	b.Run("CacheRead", func(b *testing.B) {
		cache := newMemoryCache()
		// 預填資料
		for i := 0; i < 1000; i++ {
			key := fmt.Sprintf("key_%d", i)
			set(cache, key, i)
		}

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			key := fmt.Sprintf("key_%d", i%1000)
			_, _ = get(cache, key)
		}
	})

	// 測試並發快取讀取
	b.Run("ConcurrentCacheRead", func(b *testing.B) {
		cache := newMemoryCache()
		// 預填資料
		for i := 0; i < 1000; i++ {
			key := fmt.Sprintf("key_%d", i)
			set(cache, key, i)
		}

		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				key := fmt.Sprintf("key_%d", randInt(1000))
				_, _ = get(cache, key)
			}
		})
	})
}

// BenchmarkStringComparison 基準測試：字串比較效能
func BenchmarkStringComparison(b *testing.B) {
	str1 := "go-admin-test-string-for-comparison"
	str2 := "go-admin-test-string-for-comparison"
	str3 := "different-string"

	// 測試相等字串比較
	b.Run("EqualStrings", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = str1 == str2
		}
	})

	// 測試不等字串比較
	b.Run("DifferentStrings", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = str1 == str3
		}
	})

	// 測試字串長度比較
	b.Run("StringLength", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = len(str1) == len(str2)
		}
	})
}

// 輔助函數

// randInt 產生隨機整數
func randInt(max int) int {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(max)))
	return int(n.Int64())
}

// 需要的 import（為了讓範例完整）
import (
	"encoding/json"
)

// BenchmarkAll 執行所有基準測試的便利函數
func BenchmarkAll(b *testing.B) {
	benchmarks := []struct {
		name string
		fn   func(*testing.B)
	}{
		{"PasswordHashing", BenchmarkPasswordHashing},
		{"PasswordVerification", BenchmarkPasswordVerification},
	}

	for _, benchmark := range benchmarks {
		b.Run(benchmark.name, benchmark.fn)
	}
}

// 效能測試報告範例
/*
執行基準測試：
go test -bench=. -benchmem ./test/benchmark/...

範例輸出：
BenchmarkPasswordHashing-8              100      15234567 ns/op      64 B/op       1 allocs/op
BenchmarkPasswordVerification-8        1000       1523456 ns/op       0 B/op       0 allocs/op
BenchmarkStringOperations/StringConcatenation-8    50000000    34.2 ns/op      32 B/op       2 allocs/op
BenchmarkStringOperations/FmtSprintf-8            10000000   152.0 ns/op      64 B/op       3 allocs/op

說明：
- 第一個數字：執行次數
- ns/op：每次操作的納秒數
- B/op：每次操作分配的位元組數
- allocs/op：每次操作的記憶體分配次數
*/