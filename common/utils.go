// utils.go
package common

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func GenerateTransactionID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return fmt.Sprintf("txn_%x_%d", b, time.Now().Unix())
}

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// 模拟的HTTP服务端，用于演示
func StartMockServers() {
	// 启动库存服务模拟
	go startMockInventoryService()

	// 启动订单服务模拟
	go startMockOrderService()
}

func startMockInventoryService() {
	http.HandleFunc("/api/inventory/check", func(w http.ResponseWriter, r *http.Request) {
		// 模拟随机失败（30%概率失败）
		if time.Now().UnixNano()%10 < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "库存不足"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success": true, "message": "库存验证通过"}`))
	})

	log.Println("库存服务模拟器启动在 :8081")
	http.ListenAndServe(":8081", nil)
}

func startMockOrderService() {
	http.HandleFunc("/api/order/create", func(w http.ResponseWriter, r *http.Request) {
		// 模拟随机失败（20%概率失败）
		if time.Now().UnixNano()%10 < 2 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error": "创建订单失败"}`))
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"success": true, "order_id": "order_123456"}`))
	})

	log.Println("订单服务模拟器启动在 :8082")
	http.ListenAndServe(":8082", nil)
}
