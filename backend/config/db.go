package config


import (

   "fmt"

   "time"

   "example.com/se/entity"

   "gorm.io/driver/sqlite"

   "gorm.io/gorm"

)


var db *gorm.DB


func DB() *gorm.DB {

   return db

}


func ConnectionDB() {

   database, err := gorm.Open(sqlite.Open("vote_system.db?cache=shared"), &gorm.Config{})

   if err != nil {

       panic("failed to connect database")

   }

   fmt.Println("connected database")

   db = database

}


func SetupDatabase() {


   db.AutoMigrate(

       &entity.Users{},

       &entity.Genders{},

       &entity.Candidates{},

       &entity.Elections{},

       &entity.Votes{},

   )


   GenderMale := entity.Genders{Gender: "Male"}

   GenderFemale := entity.Genders{Gender: "Female"}


   db.FirstOrCreate(&GenderMale, &entity.Genders{Gender: "Male"})

   db.FirstOrCreate(&GenderFemale, &entity.Genders{Gender: "Female"})


   hashedPassword, _ := HashPassword("1")

   BirthDay, _ := time.Parse("2006-01-02", "1988-11-12")

   User := &entity.Users{

       FirstName: "Thuwanon",

       LastName:  "Deethala",

       Email:     "admin@gmail.com",

       Age:       99,

       Password:  hashedPassword,

       Role:    "admin",

       BirthDay:  BirthDay,

       GenderID:  1,

   }

   db.FirstOrCreate(User, &entity.Users{

       Email: "admin@gmail.com",

   })


}