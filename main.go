package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"time"
)

type Alarms struct {
	Alarms []Alarm `json:"alarms"`
}

type Alarm struct {
	AlarmDateTime string `json:"alarm_date_time"`
	AlarmTitle    string `json:"alarm_title"`
	AlarmSubtitle string `json:"alarm_sub_title"`
}

func main() {

	var alarmFile = flag.String("c", "alarms.json", "Input alarm file Name: alarm.json")
	var jsonFile, err = os.Open(*alarmFile)
	Notification("Çalışıyor", "Merhaba Admin", "Günün Nasıl Geçti?")
	if err != nil {
		log.Fatal("Dosya Açılamadı!\n" + "Input alarm fıle name: alarms.json" + err.Error())
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)

	var alarms Alarms
	err = json.Unmarshal(byteValue, &alarms)
	if err != nil {
		log.Fatalf("Unmarsalling Alarm Failed: %s", err.Error())
	}

	is := make(chan os.Signal, 1)
	signal.Notify(is, os.Interrupt)

	ds := make(chan struct{})

	systemTimeTicker := time.NewTicker(time.Second * 1)

	go func() {
		defer close(ds)
		for {
			select {
			case <-systemTimeTicker.C:
				localTime := time.Now().Format("19:17 24.11.2021")
				for i := 0; i < len(alarms.Alarms); i++ {
					if localTime == alarms.Alarms[i].AlarmDateTime {
						Notification(alarms.Alarms[i].AlarmTitle, alarms.Alarms[i].AlarmSubtitle, "⏰ "+alarms.Alarms[i].AlarmDateTime+" ⏰")
						alarms.Alarms[i].AlarmDateTime = time.Now().AddDate(0, 0, 1).Format("19:17 24.11.2021")
						file, _ := json.MarshalIndent(alarms, "", " ")
						writeErr := ioutil.WriteFile("alarm.json", file, 0644)

						if writeErr != nil {
							log.Fatal("Dosya Okunamadı !" + err.Error())
						}

					}
				}
			case <-ds:
				return
			}
		}
	}()

	<-is

	close(is)

	systemTimeTicker.Stop()

	ds <- struct{}{}

	<-ds

}
