package analyze

import (
	"fmt"
	"testing"
	"time"
)

func TestGetScheduleData(t *testing.T) {
	schedule, err := GetScheduleData()
	if err != nil {
		t.Fatalf("Failed to get schedule: %v", err)
	}
	if len(schedule) == 0 {
		t.Fatal("Schedule data is empty")
	}
	t.Logf("Schedule data length: %d bytes", len(schedule))
	fmt.Printf("%s\n", schedule)
}

func TestGetSchedule(t *testing.T) {
	schedule, err := GetSchedule()
	if err != nil {
		t.Fatalf("Failed to get schedule: %v", err)
	}
	if schedule == nil {
		t.Fatal("Schedule data is nil")
	}
	t.Logf("Schedule data: %+v", schedule)
}

func TestExportScheduleToFile(t *testing.T) {
	filename := fmt.Sprintf("schedule_%s.csv", time.Now().Format("20060102_150405"))
	if err := ExportScheduleToFile(filename); err != nil {
		t.Fatalf("Failed to export schedule to file: %v", err)
	}
	t.Logf("Schedule exported to file: %s", filename)
}
