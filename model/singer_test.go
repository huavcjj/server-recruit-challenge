package model_test

import (
	"github.com/pulse227/server-recruit-challenge-sample/model"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestSinger(t *testing.T) {
	singerID := model.SingerID(1)
	singer := model.Singer{
		ID:   singerID,
		Name: "test_singer_name",
	}

	assert.Equal(t, singerID, singer.ID)
	assert.Equal(t, "test_singer_name", singer.Name)
	assert.NoError(t, singer.Validate())
}

func TestSinger_Validate(t *testing.T) {
	validSinger := model.Singer{Name: "Valid Name"}
	assert.Equal(t, nil, validSinger.Validate())

	emptyNameSinger := model.Singer{Name: ""}
	err := emptyNameSinger.Validate()
	assert.ErrorIs(t, err, model.ErrInvalidParam)

	longNameSinger := model.Singer{Name: strings.Repeat("a", 256)}
	err = longNameSinger.Validate()
	assert.ErrorIs(t, err, model.ErrInvalidParam)
}
