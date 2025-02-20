package model_test

import (
	"github.com/pulse227/server-recruit-challenge-sample/model"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestAlbum(t *testing.T) {
	albumID := model.AlbumID(1)
	singerID := model.SingerID(1)
	singer := model.Singer{
		ID:   singerID,
		Name: "test_singer_name",
	}
	album := model.Album{
		ID:       albumID,
		Title:    "test_album_title",
		SingerID: singerID,
		Singer:   &singer,
	}
	assert.Equal(t, albumID, album.ID)
	assert.Equal(t, "test_album_title", album.Title)
	assert.Equal(t, singerID, album.SingerID)
	assert.Equal(t, singerID, album.Singer.ID)
	assert.Equal(t, "test_singer_name", album.Singer.Name)
	assert.NoError(t, album.Validate())
}

func TestAlbum_Validate(t *testing.T) {
	validAlbum := model.Album{Title: "Valid Title"}
	assert.Equal(t, nil, validAlbum.Validate())

	emptyTitleAlbum := model.Album{Title: ""}
	err := emptyTitleAlbum.Validate()
	assert.ErrorIs(t, err, model.ErrInvalidParam)

	longTitleAlbum := model.Album{Title: strings.Repeat("a", 256)}
	err = longTitleAlbum.Validate()
	assert.ErrorIs(t, err, model.ErrInvalidParam)
}
