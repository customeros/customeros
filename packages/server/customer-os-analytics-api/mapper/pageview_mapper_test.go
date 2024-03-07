package mapper

import (
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-analytics-api/graph/model"
	"github.com/openline-ai/openline-customer-os/packages/server/customer-os-analytics-api/repository/entity"
	"reflect"
	"testing"
)

var (
	pageViewEntity = entity.PageViewEntity{
		ID:             "someID",
		SessionID:      "someSessionID",
		OrderInSession: 1,
		EngagedTime:    123,
		Title:          "someTitme",
		Path:           "/some/path",
	}
	expectedPageView = model.PageView{
		ID:          pageViewEntity.ID,
		Path:        pageViewEntity.Path,
		Title:       pageViewEntity.Title,
		Order:       pageViewEntity.OrderInSession,
		EngagedTime: pageViewEntity.EngagedTime,
		SessionId:   pageViewEntity.SessionID,
	}
)

func TestMapPageViewEmpty(t *testing.T) {
	emptyExpectedPageView := model.PageView{}
	newPageViewEntity := entity.PageViewEntity{}
	pageView := MapPageView(&newPageViewEntity)
	if !reflect.DeepEqual(pageView, &emptyExpectedPageView) {
		t.Errorf("Expecting: %+v, received: %+v", emptyExpectedPageView, pageView)
	}
}

func TestMapPageView(t *testing.T) {
	pageView := MapPageView(&pageViewEntity)
	if !reflect.DeepEqual(pageView, &expectedPageView) {
		t.Errorf("Expecting: %+v, received: %+v", expectedPageView, pageView)
	}
}

func TestMapPageViewsEmpty(t *testing.T) {
	pageViews := MapPageViews(&entity.PageViewEntities{})
	if len(pageViews) != 0 {
		t.Errorf("Expecting empty pageViews, received: %+v", pageViews)
	}
}

func TestMapPageViews(t *testing.T) {
	pageViewEntities := entity.PageViewEntities{}
	pageViewEntities = append(pageViewEntities, pageViewEntity)
	pageViewEntities = append(pageViewEntities, pageViewEntity)
	pageViews := MapPageViews(&pageViewEntities)
	if len(pageViews) != 2 {
		t.Errorf("Expecting pageViews with 2 elements, received: %d elements", len(pageViews))
	}
	for i := range []int{1, 2} {
		if !reflect.DeepEqual(pageViews[i], &expectedPageView) {
			t.Errorf("Expecting: %+v, received: %+v", expectedPageView, pageViews[i])
		}
	}
}
