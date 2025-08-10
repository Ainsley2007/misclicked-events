package sqlite

import (
	"database/sql"
	"reflect"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func setupActivityDB(t *testing.T) (*sql.DB, ActivityDataSource) {
	t.Helper()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}
	create := `
CREATE TABLE activity (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    type TEXT NOT NULL,
    hiscore_names TEXT NOT NULL
);`
	if _, err := db.Exec(create); err != nil {
		t.Fatalf("failed to create activity table: %v", err)
	}
	ds := NewActivityDataSource(db)
	return db, ds
}

func TestAddActivity(t *testing.T) {
	_, ds := setupActivityDB(t)

	activity := &ActivityModel{
		Name:         "Zulrah",
		Type:         "boss",
		HiscoreNames: "Zulrah",
	}

	err := ds.AddActivity(activity)
	if err != nil {
		t.Fatalf("AddActivity error: %v", err)
	}

	activities, err := ds.GetActivityListByType("boss")
	if err != nil {
		t.Fatalf("GetActivityListByType error: %v", err)
	}

	if len(activities) != 1 {
		t.Fatalf("expected 1 activity, got %d", len(activities))
	}

	expected := &ActivityModel{
		ID:           1,
		Name:         "Zulrah",
		Type:         "boss",
		HiscoreNames: "Zulrah",
	}

	if !reflect.DeepEqual(activities[0], expected) {
		t.Errorf("expected %+v, got %+v", expected, activities[0])
	}
}

func TestAddActivityWithMultipleHiscoreNames(t *testing.T) {
	_, ds := setupActivityDB(t)

	activity := &ActivityModel{
		Name:         "Wildy",
		Type:         "boss",
		HiscoreNames: "Artio,Callisto,Cal'varion,Vet'ion,Venenatis,Spindel",
	}

	err := ds.AddActivity(activity)
	if err != nil {
		t.Fatalf("AddActivity error: %v", err)
	}

	activities, err := ds.GetActivityListByType("boss")
	if err != nil {
		t.Fatalf("GetActivityListByType error: %v", err)
	}

	if len(activities) != 1 {
		t.Fatalf("expected 1 activity, got %d", len(activities))
	}

	expected := &ActivityModel{
		ID:           1,
		Name:         "Wildy",
		Type:         "boss",
		HiscoreNames: "Artio,Callisto,Cal'varion,Vet'ion,Venenatis,Spindel",
	}

	if !reflect.DeepEqual(activities[0], expected) {
		t.Errorf("expected %+v, got %+v", expected, activities[0])
	}
}

func TestAddActivityUpsert(t *testing.T) {
	_, ds := setupActivityDB(t)

	activity1 := &ActivityModel{
		Name:         "Zulrah",
		Type:         "boss",
		HiscoreNames: "Zulrah",
	}

	err := ds.AddActivity(activity1)
	if err != nil {
		t.Fatalf("First AddActivity error: %v", err)
	}

	activity2 := &ActivityModel{
		Name:         "Zulrah",
		Type:         "boss",
		HiscoreNames: "Zulrah,Zulrah (hard mode)",
	}

	err = ds.AddActivity(activity2)
	if err != nil {
		t.Fatalf("Second AddActivity error: %v", err)
	}

	activities, err := ds.GetActivityListByType("boss")
	if err != nil {
		t.Fatalf("GetActivityListByType error: %v", err)
	}

	if len(activities) != 1 {
		t.Fatalf("expected 1 activity after upsert, got %d", len(activities))
	}

	expected := &ActivityModel{
		ID:           1,
		Name:         "Zulrah",
		Type:         "boss",
		HiscoreNames: "Zulrah,Zulrah (hard mode)",
	}

	if !reflect.DeepEqual(activities[0], expected) {
		t.Errorf("expected %+v, got %+v", expected, activities[0])
	}
}

func TestRemoveActivity(t *testing.T) {
	_, ds := setupActivityDB(t)

	activity := &ActivityModel{
		Name:         "Zulrah",
		Type:         "boss",
		HiscoreNames: "Zulrah",
	}

	err := ds.AddActivity(activity)
	if err != nil {
		t.Fatalf("AddActivity error: %v", err)
	}

	activities, err := ds.GetActivityListByType("boss")
	if err != nil {
		t.Fatalf("GetActivityListByType error: %v", err)
	}

	if len(activities) != 1 {
		t.Fatalf("expected 1 activity before removal, got %d", len(activities))
	}

	err = ds.RemoveActivity("Zulrah")
	if err != nil {
		t.Fatalf("RemoveActivity error: %v", err)
	}

	activities, err = ds.GetActivityListByType("boss")
	if err != nil {
		t.Fatalf("GetActivityListByType error after removal: %v", err)
	}

	if len(activities) != 0 {
		t.Fatalf("expected 0 activities after removal, got %d", len(activities))
	}
}

func TestRemoveActivityNonExistent(t *testing.T) {
	_, ds := setupActivityDB(t)

	err := ds.RemoveActivity("NonExistent")
	if err != nil {
		t.Fatalf("RemoveActivity on non-existent activity should not error: %v", err)
	}
}

func TestGetActivityListByType(t *testing.T) {
	_, ds := setupActivityDB(t)

	bossActivities := []*ActivityModel{
		{Name: "Zulrah", Type: "boss", HiscoreNames: "Zulrah"},
		{Name: "Corp", Type: "boss", HiscoreNames: "Corporeal Beast"},
		{Name: "COX", Type: "boss", HiscoreNames: "Chambers of Xeric,Chambers of Xeric: Challenge Mode"},
	}

	skillActivities := []*ActivityModel{
		{Name: "Slayer", Type: "skill", HiscoreNames: "Slayer"},
		{Name: "Mining", Type: "skill", HiscoreNames: "Mining"},
	}

	for _, activity := range bossActivities {
		err := ds.AddActivity(activity)
		if err != nil {
			t.Fatalf("AddActivity error: %v", err)
		}
	}

	for _, activity := range skillActivities {
		err := ds.AddActivity(activity)
		if err != nil {
			t.Fatalf("AddActivity error: %v", err)
		}
	}

	bossList, err := ds.GetActivityListByType("boss")
	if err != nil {
		t.Fatalf("GetActivityListByType for boss error: %v", err)
	}

	if len(bossList) != 3 {
		t.Fatalf("expected 3 boss activities, got %d", len(bossList))
	}

	skillList, err := ds.GetActivityListByType("skill")
	if err != nil {
		t.Fatalf("GetActivityListByType for skill error: %v", err)
	}

	if len(skillList) != 2 {
		t.Fatalf("expected 2 skill activities, got %d", len(skillList))
	}

	nonExistentList, err := ds.GetActivityListByType("non_existent")
	if err != nil {
		t.Fatalf("GetActivityListByType for non-existent type error: %v", err)
	}

	if len(nonExistentList) != 0 {
		t.Fatalf("expected 0 activities for non-existent type, got %d", len(nonExistentList))
	}
}

func TestGetActivityListEmpty(t *testing.T) {
	_, ds := setupActivityDB(t)

	activities, err := ds.GetActivityListByType("boss")
	if err != nil {
		t.Fatalf("GetActivityListByType error: %v", err)
	}

	if len(activities) != 0 {
		t.Fatalf("expected 0 activities in empty database, got %d", len(activities))
	}
}

func TestAddActivityWithEmptyFields(t *testing.T) {
	_, ds := setupActivityDB(t)

	activity := &ActivityModel{
		Name:         "",
		Type:         "",
		HiscoreNames: "",
	}

	err := ds.AddActivity(activity)
	if err != nil {
		t.Fatalf("AddActivity with empty fields should not error: %v", err)
	}

	activities, err := ds.GetActivityListByType("")
	if err != nil {
		t.Fatalf("GetActivityListByType error: %v", err)
	}

	if len(activities) != 1 {
		t.Fatalf("expected 1 activity with empty type, got %d", len(activities))
	}
}

func TestActivityIDAutoIncrement(t *testing.T) {
	_, ds := setupActivityDB(t)

	activities := []*ActivityModel{
		{Name: "Activity1", Type: "boss", HiscoreNames: "Boss1"},
		{Name: "Activity2", Type: "boss", HiscoreNames: "Boss2"},
		{Name: "Activity3", Type: "boss", HiscoreNames: "Boss3"},
	}

	for _, activity := range activities {
		err := ds.AddActivity(activity)
		if err != nil {
			t.Fatalf("AddActivity error: %v", err)
		}
	}

	activityList, err := ds.GetActivityListByType("boss")
	if err != nil {
		t.Fatalf("GetActivityListByType error: %v", err)
	}

	if len(activityList) != 3 {
		t.Fatalf("expected 3 activities, got %d", len(activityList))
	}

	for i, activity := range activityList {
		expectedID := int64(i + 1)
		if activity.ID != expectedID {
			t.Errorf("expected ID %d for activity %d, got %d", expectedID, i+1, activity.ID)
		}
	}
}

func TestRemoveActivityByNameCaseSensitive(t *testing.T) {
	_, ds := setupActivityDB(t)

	activity := &ActivityModel{
		Name:         "Zulrah",
		Type:         "boss",
		HiscoreNames: "Zulrah",
	}

	err := ds.AddActivity(activity)
	if err != nil {
		t.Fatalf("AddActivity error: %v", err)
	}

	err = ds.RemoveActivity("zulrah")
	if err != nil {
		t.Fatalf("RemoveActivity with different case should not error: %v", err)
	}

	activities, err := ds.GetActivityListByType("boss")
	if err != nil {
		t.Fatalf("GetActivityListByType error: %v", err)
	}

	if len(activities) != 1 {
		t.Fatalf("expected 1 activity after case-sensitive removal attempt, got %d", len(activities))
	}

	err = ds.RemoveActivity("Zulrah")
	if err != nil {
		t.Fatalf("RemoveActivity with correct case should not error: %v", err)
	}

	activities, err = ds.GetActivityListByType("boss")
	if err != nil {
		t.Fatalf("GetActivityListByType error: %v", err)
	}

	if len(activities) != 0 {
		t.Fatalf("expected 0 activities after correct case removal, got %d", len(activities))
	}
}

func TestAddActivityUniqueConstraint(t *testing.T) {
	_, ds := setupActivityDB(t)

	activity1 := &ActivityModel{
		Name:         "Zulrah",
		Type:         "boss",
		HiscoreNames: "Zulrah",
	}

	err := ds.AddActivity(activity1)
	if err != nil {
		t.Fatalf("First AddActivity error: %v", err)
	}

	activity2 := &ActivityModel{
		Name:         "Zulrah",
		Type:         "skill",
		HiscoreNames: "Different hiscore name",
	}

	err = ds.AddActivity(activity2)
	if err != nil {
		t.Fatalf("Second AddActivity should not error due to upsert: %v", err)
	}

	activities, err := ds.GetActivityListByType("boss")
	if err != nil {
		t.Fatalf("GetActivityListByType for boss error: %v", err)
	}

	if len(activities) != 0 {
		t.Fatalf("expected 0 boss activities after upsert to skill, got %d", len(activities))
	}

	skillActivities, err := ds.GetActivityListByType("skill")
	if err != nil {
		t.Fatalf("GetActivityListByType for skill error: %v", err)
	}

	if len(skillActivities) != 1 {
		t.Fatalf("expected 1 skill activity after upsert, got %d", len(skillActivities))
	}

	expected := &ActivityModel{
		ID:           1,
		Name:         "Zulrah",
		Type:         "skill",
		HiscoreNames: "Different hiscore name",
	}

	if !reflect.DeepEqual(skillActivities[0], expected) {
		t.Errorf("expected %+v, got %+v", expected, skillActivities[0])
	}
}
