package users


import (

    "net/http"


    "github.com/gin-gonic/gin"
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/push"

    "example.com/se/config"

    "example.com/se/entity"
    "log"
)

const pushGatewayURL = "http://pushgateway:9091" // เปลี่ยนตาม URL ของคุณ

func sendUserUpdateMetric() {
    gauge := prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "users_updated_total",
        Help: "Total number of user updates.",
    })

    gauge.Set(1) // 1 ครั้ง ต่อการ update

    if err := push.New(pushGatewayURL, "user_job").
        Collector(gauge).
        Push(); err != nil {
        log.Printf("Could not push to Pushgateway: %v", err)
    }
}

func sendUserDeleteMetric() {
    gauge := prometheus.NewGauge(prometheus.GaugeOpts{
        Name: "users_deleted_total",
        Help: "Total number of user deletes.",
    })

    gauge.Set(1) // 1 ครั้ง ต่อการ delete

    if err := push.New(pushGatewayURL, "user_job").
        Collector(gauge).
        Push(); err != nil {
        log.Printf("Could not push to Pushgateway: %v", err)
    }
}


func GetAll(c *gin.Context) {


   var users []entity.Users


   db := config.DB()

   results := db.Preload("Gender").Find(&users)

   if results.Error != nil {

       c.JSON(http.StatusNotFound, gin.H{"error": results.Error.Error()})

       return

   }

   c.JSON(http.StatusOK, users)


}


func Get(c *gin.Context) {


   ID := c.Param("id")

   var user entity.Users


   db := config.DB()

   results := db.Preload("Gender").First(&user, ID)

   if results.Error != nil {

       c.JSON(http.StatusNotFound, gin.H{"error": results.Error.Error()})

       return

   }

   if user.ID == 0 {

       c.JSON(http.StatusNoContent, gin.H{})

       return

   }

   c.JSON(http.StatusOK, user)


}


func Update(c *gin.Context) {
    var user entity.Users
    UserID := c.Param("id")
    db := config.DB()
    result := db.First(&user, UserID)
    if result.Error != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "id not found"})
        return
    }

    if err := c.ShouldBindJSON(&user); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request, unable to map payload"})
        return
    }

    result = db.Save(&user)
    if result.Error != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Bad request"})
        return
    }

    // ส่ง metric ไป Pushgateway
    sendUserUpdateMetric()

    c.JSON(http.StatusOK, gin.H{"message": "Updated successful"})
}


func Delete(c *gin.Context) {
    id := c.Param("id")
    db := config.DB()
    if tx := db.Exec("DELETE FROM users WHERE id = ?", id); tx.RowsAffected == 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "id not found"})
        return
    }

    // ส่ง metric ไป Pushgateway
    sendUserDeleteMetric()

    c.JSON(http.StatusOK, gin.H{"message": "Deleted successful"})
}

