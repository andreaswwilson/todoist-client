package todoist

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestProject(t *testing.T) {
	logrus.SetLevel(logrus.ErrorLevel)
	token := os.Getenv("TODOIST_TOKEN")
	c := NewClient(token)
	ctx := context.Background()

	name := fmt.Sprint(time.Now().Unix())
	projectId := ""
	t.Run("create project", func(t *testing.T) {
		project := &CreateProject{
			Name: &name,
		}
		res, err := c.CreateProject(ctx, *project)
		projectId = res.ID
		assert.Nil(t, err, "expecting nil error")
		assert.NotNil(t, res, "expecting non-nil result")
	})

	t.Run("read project", func(t *testing.T) {
		res, err := c.GetProject(ctx, projectId)
		assert.Nil(t, err, "expecting nil error")
		assert.NotNil(t, res, "expecting non-nil result")
		assert.Equal(t, res.Name, name, fmt.Sprintf("expecting %s, got %s", name, res.Name))
	})

	t.Run("update project", func(t *testing.T) {
		name = "new name"
		project := &UpdateProject{
			ID:   &projectId,
			Name: &name,
		}
		res, err := c.UpdateProject(ctx, *project)
		assert.Nil(t, err, "expecting nil error")
		assert.Equal(t, res.Name, "new name", "expecting project name to be %s, got %s", &name, res.Name)

	})

	t.Run("delete project", func(t *testing.T) {
		err := c.DeleteProject(ctx, projectId)
		assert.Nil(t, err, "expecting nil error")
	})

}
