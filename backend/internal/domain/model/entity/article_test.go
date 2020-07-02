package entity

import (
	uuid "github.com/satori/go.uuid"
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestCreateArticle(t *testing.T) {
	type args struct {
		title string
		url   string
	}
	validURL := "http://go.com/lol"
	validTitle := "valid title"
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"title too long",
			args{
				strings.Repeat("I", TitleMaxLength+1), validURL},
			true,
		},
		{
			"invalid url",
			args{
				"valid title", "olo"},
			true,
		},
		{
			"no error",
			args{
				validTitle, "http://go.com/lol"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := CreateArticle(tt.args.title, tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateArticle() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestNewEmptyTags(t *testing.T) {
	if ! reflect.DeepEqual(NewEmptyTags(), []*Tag{}) {
		t.Error("NewEmptyTags() doesn't return []*Tag{}")
	}
}

func TestArticle_HasTag(t *testing.T) {
	tagUUID := uuid.NewV4()
	testTag := &Tag{
		ID:   tagUUID,
		Name: "valid tag",
	}
	testTag2 := &Tag{
		ID:   uuid.NewV4(),
		Name: "valid tag II",
	}
	testArticle, err := NewArticle(
		uuid.NewV4(),
		"valid",
		"http://go.com",
		time.Now(),
		time.Time{},
		[]*Tag{testTag},
	)
	if err != nil {
		t.Error("cannot create test article")
	}
	t.Run("shouldn't have found", func(t *testing.T) {
		if testArticle.HasTag(testTag2) {
			t.Error("tag found but isn't here")
		}
	})
	t.Run("should have found", func(t *testing.T) {
		if !testArticle.HasTag(testTag) {
			t.Error("tag not found but is here")
		}
	})
}

func TestArticle_AddTag(t *testing.T) {
	testArticle, err := CreateArticle("valid", "http://go.com")
	if err != nil {
		t.Error("cannot create test article")
	}
	testTag, err := CreateTag("test tag")
	if err != nil {
		t.Error("cannot create test tag")
	}
	err = testArticle.AddTag(testTag)
	if err != nil {
		t.Errorf("method call failed, error : %s", err)
	}
	t.Run("valid", func(t *testing.T) {
		if !testArticle.HasTag(testTag) {
			t.Error("tag not present but it should be")
		}
	})
	t.Run("already has tag", func(t *testing.T) {
		if err = testArticle.AddTag(testTag); err == nil {
			t.Error("tag is already present, error should have occurred")
		}
	})
}

func TestArticle_CountTags(t *testing.T) {

	testTag, err := CreateTag("valid tag I")
	if err != nil {
		t.Error("cannot create test tag")
	}
	testTag2, err := CreateTag("valid tag II")
	if err != nil {
		t.Error("cannot create test tag")
	}
	testArticleWithoutTags, err := CreateArticle("valid title", "http://go.com/lol")
	if err != nil {
		t.Error("cannot create test article")
	}
	testArticleWithOneTag, err := NewArticle(
		uuid.NewV4(),
		"valid",
		"http://go.com",
		time.Now(),
		time.Time{},
		[]*Tag{testTag},
	)
	if err != nil {
		t.Error("cannot create test article")
	}
	testArticleWithTwoTags, err := NewArticle(
		uuid.NewV4(),
		"valid",
		"http://go.com",
		time.Now(),
		time.Time{},
		[]*Tag{testTag, testTag2},
	)
	if err != nil {
		t.Error("cannot create test article")
	}
	t.Run("test-no-tag", func(t *testing.T) {
		if cnt := testArticleWithoutTags.CountTags() ; cnt != 0 {
			t.Errorf("expected 0, got %d", cnt)
		}
	})
	t.Run("test-one-tag", func(t *testing.T) {
		if cnt := testArticleWithOneTag.CountTags() ; cnt != 1 {
			t.Errorf("expected 1, got %d", cnt)
		}
	})
	t.Run("test-two-tag", func(t *testing.T) {
		if cnt := testArticleWithTwoTags.CountTags() ; cnt != 2 {
			t.Errorf("expected 2, got %d", cnt)
		}
	})
}

func TestArticle_RemoveTag(t *testing.T) {
	tagUUID := uuid.NewV4()
	testTag := &Tag{
		ID:   tagUUID,
		Name: "valid tag",
	}
	testTag2 := &Tag{
		ID:   uuid.NewV4(),
		Name: "valid tag II",
	}
	testTag3 := &Tag{
		ID:   uuid.NewV4(),
		Name: "valid tag III",
	}
	testArticle, err := NewArticle(
		uuid.NewV4(),
		"valid",
		"http://go.com",
		time.Now(),
		time.Time{},
		[]*Tag{testTag},
	)
	if err != nil {
		t.Error("cannot create test article")
	}
	t.Run("1-tag-removing_non_existing_tag", func(t *testing.T) {
		err := testArticle.RemoveTag(testTag2)
		if err == nil {
			t.Error("error expected but did not occurred")
		}
	})
	t.Run("1-tag-removing_existing_tag", func(t *testing.T) {
		err := testArticle.RemoveTag(testTag)
		if err != nil {
			t.Error("error not expected but did occurred")
		}
		if testArticle.HasTag(testTag) {
			t.Error("no error occurred, but tag is still here")
		}
		if testArticle.CountTags() != 0 {
			t.Errorf("tags count should be %d but is %d", 0, testArticle.CountTags())
		}
	})
	err = testArticle.AddTag(testTag)
	if err != nil {
		t.Error("cannot create test tag")
	}
	err = testArticle.AddTag(testTag2)
	if err != nil {
		t.Error("cannot create test tag")
	}
	t.Run("2-tag-removing_non_existing_tag", func(t *testing.T) {
		err := testArticle.RemoveTag(testTag3)
		if err == nil {
			t.Error("error expected but did not occurred")
		}
	})
	t.Run("2-tag-removing_existing_tag", func(t *testing.T) {
		previous := testArticle.CountTags()
		err := testArticle.RemoveTag(testTag)
		if err != nil {
			t.Error("error not expected but did occurred")
		}
		if testArticle.HasTag(testTag) {
			t.Error("no error occurred, but tag is still here")
		}
		current := testArticle.CountTags()
		if current != previous - 1 {
			t.Errorf("tags count should be %d but is %d", previous - 1, current)
		}
	})
}

func TestArticle_HasBeenRead(t *testing.T) {
	testArticleNotAlreadyRead, err := CreateArticle("valid", "http://go.com")
	if err != nil {
		t.Error("cannot create test article")
	}
	testArticleAlreadyRead, err := NewArticle(
		uuid.NewV4(),
		"valid",
		"http://go.com",
		time.Now(),
		time.Now(),
		NewEmptyTags(),
	)
	if err != nil {
		t.Error("cannot create test article")
	}
	t.Run("not already read", func(t *testing.T) {
		if testArticleNotAlreadyRead.HasBeenRead() {
			t.Error("got true, expected false")
		}
	})
	t.Run("already read", func(t *testing.T) {
		if !testArticleAlreadyRead.HasBeenRead() {
			t.Error("got false, expected true")
		}
	})
}

func TestArticle_Read(t *testing.T) {

	testArticle, err := NewArticle(
		uuid.NewV4(),
		"valid",
		"http://go.com",
		time.Time{},
		time.Now(),
		NewEmptyTags(),
	)
	if err != nil {
		t.Error("cannot create test article")
	}
	t.Run("ok", func(t *testing.T) {
		if err := testArticle.Read(); err != nil {
			t.Errorf("did not expected error, got %s", err)
		}
	})
	t.Run("already read", func(t *testing.T) {
		if err := testArticle.Read(); err == nil {
			t.Error("expected error, got none")
		}
	})
}

func TestArticle_Id(t *testing.T) {
	testUUID := uuid.NewV4()
	testArticle, err := NewArticle(
		testUUID,
		"valid",
		"http://go.com",
		time.Time{},
		time.Now(),
		NewEmptyTags(),
	)
	if err != nil {
		t.Error("cannot create test article")
	}
	if testArticle.Id() != testUUID {
		t.Errorf("expected %s got %s", testUUID.String(), testArticle.Id().String())
	}
}

func TestArticle_ReadDateTime(t *testing.T) {
	testTime := time.Now()
	testArticle, err := NewArticle(
		uuid.NewV4(),
		"valid",
		"http://go.com",
		testTime,
		time.Now(),
		NewEmptyTags(),
	)
	if err != nil {
		t.Error("cannot create test article")
		return
	}
	if returnedTime := testArticle.ReadDateTime() ; returnedTime != testTime {
		t.Errorf("expected %s got %s", testTime.String(), returnedTime.String())
	}
}


func TestArticle_SavedDateTime(t *testing.T) {
	testTime := time.Now()
	testArticle, err := NewArticle(
		uuid.NewV4(),
		"valid",
		"http://go.com",
		time.Now(),
		testTime,
		NewEmptyTags(),
	)
	if err != nil {
		t.Error("cannot create test article")
		return
	}
	if returnedTime := testArticle.SavedDateTime() ; returnedTime != testTime {
		t.Errorf("expected %s got %s", testTime.String(), returnedTime.String())
	}
}

func TestArticle_Title(t *testing.T) {
	expected := "valid title"
	testArticle, err := CreateArticle(expected, "http://go.com")
	if err != nil {
		t.Error("cannot create test article")
		return
	}
	if returned := testArticle.Title() ; returned != expected {
		t.Errorf("expected %s got %s", returned, expected)
	}
}


func TestArticle_Url(t *testing.T) {
	expected := "http://go.com"
	testArticle, err := CreateArticle("valid title", expected)
	if err != nil {
		t.Error("cannot create test article")
		return
	}
	if returned := testArticle.Url() ; returned != expected {
		t.Errorf("expected %s got %s", returned, expected)
	}
}


func TestArticle_Tags(t *testing.T) {
	tag1, _ := CreateTag("tag 1")
	tag2, _ := CreateTag("tag 2")
	tag3, _ := CreateTag("tag 3")
	expected := []*Tag{tag1, tag2, tag3}
	testArticle, err := NewArticle(
		uuid.NewV4(),
		"valid",
		"http://go.com",
		time.Now(),
		time.Now(),
		expected,
	)
	if err != nil {
		t.Error("cannot create test article")
		return
	}
	if returned := testArticle.Tags() ; reflect.DeepEqual(returned, &expected) {
		t.Error("result gotten differs from expected result")
	}
}