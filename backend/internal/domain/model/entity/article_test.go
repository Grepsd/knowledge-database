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
			"no error",
			args{
				validTitle, "http://go.com/lol"},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArticleURL, err := NewArticleURL(tt.args.url)
			if err != nil {
				t.Error("invalid URL")
			}
			_, err = CreateArticle(NewArticleID(uuid.NewV4()), NewArticleTitle(tt.args.title), tArticleURL)
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
	testTag, err := NewTag(NewTagID(uuid.NewV4()), NewTagName("valie I"))
	testTag2, err := NewTag(NewTagID(uuid.NewV4()), NewTagName("valid II"))
	testArticleURL, err := NewArticleURL("http://go.com")
	if err != nil {
		t.Error("invalid URL")
	}
	testArticle, err := NewArticle(
		NewArticleID(uuid.NewV4()),
		NewArticleTitle("valid"),
		testArticleURL,
		NewArticleReadDateTime(time.Time{}),
		NewArticleSavedDateTime(time.Now()),
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
	testArticleURL, err := NewArticleURL("http://go.com")
	testArticle, err := CreateArticle(NewArticleID(uuid.NewV4()), NewArticleTitle("valid"), testArticleURL)
	if err != nil {
		t.Error("cannot create test article")
	}
	testTag, err := NewTag(NewTagID(uuid.NewV4()), NewTagName("test tag"))
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

	testTag, err := NewTag(NewTagID(uuid.NewV4()), NewTagName("valid tag I"))
	if err != nil {
		t.Error("cannot create test tag")
	}
	testTag2, err := NewTag(NewTagID(uuid.NewV4()), NewTagName("valid tag II"))
	if err != nil {
		t.Error("cannot create test tag")
	}
	testArticleURL, _ := NewArticleURL("http://go.com")
	testArticleWithoutTags, err := CreateArticle(NewArticleID(uuid.NewV4()), NewArticleTitle("valid title"), testArticleURL)
	if err != nil {
		t.Error("cannot create test article")
	}
	testTags := NewArticleTags()
	testTags = append(testTags, testTag)
	testArticleWithOneTag, err := NewArticle(
		NewArticleID(uuid.NewV4()),
		NewArticleTitle("valid"),
		testArticleURL,
		NewArticleReadDateTime(time.Time{}),
		NewArticleSavedDateTime(time.Now()),
		[]*Tag{testTag},
	)
	if err != nil {
		t.Error("cannot create test article")
	}
	testTags2 := append(testTags, testTag2)
	testArticleWithTwoTags, err := NewArticle(
		NewArticleID(uuid.NewV4()),
		NewArticleTitle("valid"),
		testArticleURL,
		NewArticleReadDateTime(time.Time{}),
		NewArticleSavedDateTime(time.Now()),
		testTags2,
	)
	if err != nil {
		t.Error("cannot create test article")
	}
	t.Run("test-no-tag", func(t *testing.T) {
		if cnt := testArticleWithoutTags.CountTags(); cnt != 0 {
			t.Errorf("expected 0, got %d", cnt)
		}
	})
	t.Run("test-one-tag", func(t *testing.T) {
		if cnt := testArticleWithOneTag.CountTags(); cnt != 1 {
			t.Errorf("expected 1, got %d", cnt)
		}
	})
	t.Run("test-two-tag", func(t *testing.T) {
		if cnt := testArticleWithTwoTags.CountTags(); cnt != 2 {
			t.Errorf("expected 2, got %d", cnt)
		}
	})
}

func TestArticle_RemoveTag(t *testing.T) {
	testTag, err := NewTag(NewTagID(uuid.NewV4()), NewTagName("valid name I"))
	testTag2, err := NewTag(NewTagID(uuid.NewV4()), NewTagName("valid name II"))
	testTag3, err := NewTag(NewTagID(uuid.NewV4()), NewTagName("valid name III"))
	testArticleURL, _ := NewArticleURL("http://go.com")

	testTags1 := NewEmptyTags()
	testTags1 = append(testTags1, testTag)
	testArticle, err := NewArticle(
		NewArticleID(uuid.NewV4()),
		NewArticleTitle("valid"),
		testArticleURL,
		NewArticleReadDateTime(time.Time{}),
		NewArticleSavedDateTime(time.Now()),
		testTags1,
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
		if current != previous-1 {
			t.Errorf("tags count should be %d but is %d", previous-1, current)
		}
	})
}

func TestArticle_HasBeenRead(t *testing.T) {
	testArticleURL, _ := NewArticleURL("http://go.com")
	testArticleNotAlreadyRead, err := CreateArticle(NewArticleID(uuid.NewV4()), NewArticleTitle("valid"), testArticleURL)
	if err != nil {
		t.Error("cannot create test article")
	}
	readDateTime := NewArticleReadDateTime(time.Time{})
	testArticleAlreadyRead, err := NewArticle(
		NewArticleID(uuid.NewV4()),
		NewArticleTitle("valid"),
		testArticleURL,
		readDateTime,
		NewArticleSavedDateTime(time.Now()),
		NewEmptyTags(),
	)
	if err != nil {
		t.Error("cannot create test article")
	}
	_ = testArticleAlreadyRead.UpdateReadDateTime(time.Now())
	if err != nil {
		t.Error("failed to update read date time")
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
	testArticleURL, _ := NewArticleURL("http://go.com")
	readDateTime := NewArticleReadDateTime(time.Time{})
	testArticle, err := NewArticle(
		NewArticleID(uuid.NewV4()),
		NewArticleTitle("valid"),
		testArticleURL,
		readDateTime,
		NewArticleSavedDateTime(time.Now()),
		NewEmptyTags(),
	)
	if err != nil {
		t.Error("cannot create test article")
	}
	t.Run("ok", func(t *testing.T) {
		if err := testArticle.Read(time.Now()); err != nil {
			t.Errorf("did not expected error, got %s", err.Error())
		}
	})
	t.Run("already read", func(t *testing.T) {
		if err := testArticle.Read(time.Now()); err == nil {
			t.Error("expected error, got none")
		}
	})
}

func TestArticle_Id(t *testing.T) {
	articleId := NewArticleID(uuid.NewV4())
	testArticleURL, _ := NewArticleURL("http://go.com")
	testArticle, err := NewArticle(
		articleId,
		NewArticleTitle("valid"),
		testArticleURL,
		NewArticleReadDateTime(time.Time{}),
		NewArticleSavedDateTime(time.Now()),
		NewEmptyTags(),
	)
	if err != nil {
		t.Error("cannot create test article")
	}
	testArticleId := testArticle.Id()
	if testArticleId != articleId {
		t.Errorf("expected %s got %s", articleId.String(), testArticleId.String())
	}
}

func TestArticle_ReadDateTime(t *testing.T) {
	testTime := NewArticleReadDateTime(time.Time{})
	testArticleURL, _ := NewArticleURL("http://go.com")
	testTime.Update(time.Now())
	testArticle, err := NewArticle(
		NewArticleID(uuid.NewV4()),
		NewArticleTitle("valid"),
		testArticleURL,
		testTime,
		NewArticleSavedDateTime(time.Now()),
		NewEmptyTags(),
	)
	if err != nil {
		t.Error("cannot create test article")
		return
	}
	if returnedTime := testArticle.ReadDateTime(); returnedTime != testTime {
		t.Errorf("expected %s got %s", testTime.String(), returnedTime.String())
	}
}

func TestArticle_SavedDateTime(t *testing.T) {
	savedDateTime := NewArticleSavedDateTime(time.Now())
	readDateTime := NewArticleReadDateTime(time.Time{})
	testArticleURL, _ := NewArticleURL("http://go.com")
	readDateTime.Update(time.Now())
	testArticle, err := NewArticle(
		NewArticleID(uuid.NewV4()),
		NewArticleTitle("valid"),
		testArticleURL,
		readDateTime,
		savedDateTime,
		NewEmptyTags(),
	)
	if err != nil {
		t.Error("cannot create test article")
		return
	}
	if returnedTime := testArticle.SavedDateTime(); returnedTime != savedDateTime {
		t.Errorf("expected %s got %s", savedDateTime.String(), returnedTime.String())
	}
}

func TestArticle_Title(t *testing.T) {
	expected := NewArticleTitle("valid title")
	testArticleURL, _ := NewArticleURL("http://go.com")
	testArticle, err := CreateArticle(NewArticleID(uuid.NewV4()), expected, testArticleURL)
	if err != nil {
		t.Error("cannot create test article")
		return
	}
	if returned := testArticle.Title(); returned != expected {
		t.Errorf("expected %s got %s", returned, expected)
	}
}

func TestArticle_Url(t *testing.T) {
	expected, _ := NewArticleURL("http://go.com")
	testArticle, err := CreateArticle(NewArticleID(uuid.NewV4()), NewArticleTitle("valid title"), expected)
	if err != nil {
		t.Error("cannot create test article")
		return
	}
	if returned := testArticle.Url(); returned != expected {
		t.Errorf("expected %s got %s", returned, expected)
	}
}

func TestArticle_Tags(t *testing.T) {
	tag1, _ := NewTag(NewTagID(uuid.NewV4()), NewTagName("tag I"))
	tag2, _ := NewTag(NewTagID(uuid.NewV4()), NewTagName("tag II"))
	tag3, _ := NewTag(NewTagID(uuid.NewV4()), NewTagName("tag III"))
	testArticleURL, _ := NewArticleURL("http://go.com")
	savedDateTime := NewArticleSavedDateTime(time.Now())
	readDateTime := NewArticleReadDateTime(time.Time{})
	readDateTime.Update(time.Now())
	expected := []*Tag{tag1, tag2, tag3}
	testArticle, err := NewArticle(
		NewArticleID(uuid.NewV4()),
		NewArticleTitle("valid"),
		testArticleURL,
		readDateTime,
		savedDateTime,
		expected,
	)
	if err != nil {
		t.Error("cannot create test article")
		return
	}
	if returned := testArticle.Tags(); reflect.DeepEqual(returned, &expected) {
		t.Error("result gotten differs from expected result")
	}
}

func TestArticle_Strings(t *testing.T) {
	id := NewArticleID(uuid.NewV4())
	title := NewArticleTitle("valid")
	url, _ := NewArticleURL("http://go.com")
	readDateTime := NewArticleReadDateTime(time.Time{})
	readDateTime.Update(time.Now())
	savedDateTime := NewArticleSavedDateTime(time.Now())
	t.Run("valid id.String", func(t *testing.T) {
		if id.String() != id.value.String() {
			t.Errorf("expected %s got %s", id.String(), id.value.String())
		}
	})
	t.Run("valid title.String", func(t *testing.T) {
		if title.String() != title.value {
			t.Errorf("expected %s got %s", title.String(), title.value)
		}
	})
	t.Run("valid url.String", func(t *testing.T) {
		if url.String() != url.value {
			t.Errorf("expected %s got %s", url.String(), url.value)
		}
	})
	t.Run("valid readDateTime.String", func(t *testing.T) {
		if readDateTime.String() != readDateTime.value.String() {
			t.Errorf("expected %s got %s", readDateTime.String(), readDateTime.value.String())
		}
	})
	t.Run("valid savedDateTime.String", func(t *testing.T) {
		if savedDateTime.String() != savedDateTime.value.String() {
			t.Errorf("expected %s got %s", savedDateTime.String(), savedDateTime.value.String())
		}
	})
}

func TestNewArticleURL(t *testing.T) {
	type args struct {
		value string
	}
	validURL, invalidURL := "http://go.com", "sdfjifjasdf"
	tests := []struct {
		name    string
		url     string
		want    ArticleURL
		wantErr bool
	}{
		{
			"valid",
			validURL,
			ArticleURL{"http://go.com"},
			false,
		},
		{
			"invalid",
			invalidURL,
			ArticleURL{},
			true,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewArticleURL(tt.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewArticleURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewArticleURL() got = %v, want %v", got, tt.want)
			}
		})
	}
}
