package metrics

import (
    "log"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/push"
)

// สร้าง metrics ต่าง ๆ
var (
    TotalUsers = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "online_voting_users_total",
        Help: "Total number of users in the online voting system",
    })

    TotalVotes = prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "online_voting_votes_total",
        Help: "Total number of votes cast",
    })
)

// ฟังก์ชันอัปเดตค่าของ metrics และ push ไปยัง Pushgateway
func PushMetrics(pushGatewayURL string, jobName string, userCount int, voteCount int) {
    TotalUsers.Set(float64(userCount))
    TotalVotes.Set(float64(voteCount))

    if err := push.New(pushGatewayURL, jobName).
        Collector(TotalUsers).
        Collector(TotalVotes).
        Push(); err != nil {
        log.Printf("Could not push metrics to Pushgateway: %v", err)
    } else {
        log.Println("Metrics pushed to Pushgateway successfully")
    }
}
