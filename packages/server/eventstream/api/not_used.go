
func (h *WebsiteEventsHandler) assignEventToSession(ctx context.Context, trackerData *postgres_entity.WebTrackerEvents) error {
	span, _ := opentracing.StartSpanFromContext(ctx, "WebsiteTrackerEventsHandler.assignEventsToSession")
	defer span.Finish()
	tracing.TagComponentRest(span)

	span.LogKV(
		"ip", trackerData.IP,
		"hostname", trackerData.Hostname,
		"visitor_id", trackerData.VisitorID,
		"event_type", trackerData.EventType,
	)

	query := postgres_entity.WebSession{
		Tenant:    trackerData.Tenant,
		IP:        trackerData.IP,
		Hostname:  trackerData.Hostname,
		VisitorID: trackerData.VisitorID,
		IsActive:  true,
	}

	session, err := h.services.Repositories.PostgresRepositories.WebSessionRepository.FindSession(ctx, query, nil)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}

	if session == nil && enum.IsActiveWebTrackerEvent(trackerData.EventType) {
		session, err = h.createWebSession(ctx, trackerData)
		if err != nil {
			tracing.TraceErr(span, err)
			return err
		}
	}

	if session == nil {
		err = errors.New("session not found and not created")
		if trackerData.EventType != enum.WebTrackerPageExit.String() {
			tracing.TraceErr(span, err)
		}
		return err
	}

	trackerData.SessionID = session.ID
	err = h.updateSessionLastActivity(ctx, trackerData.SessionID, trackerData.EventType)
	if err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (h *WebsiteEventsHandler) updateSessionLastActivity(ctx context.Context, sessionID, eventType string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebsiteTrackerEventsHandler.updateSessionLastActivityTimestamp")
	defer span.Finish()
	tracing.TagComponentRest(span)
	span.LogKV("sessionID", sessionID, "eventType", eventType)

	_, err := h.services.Repositories.PostgresRepositories.WebSessionRepository.UpdateLastActivity(ctx, sessionID, eventType)
	if err != nil {
		err = errors.New("unable to update web session")
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (h *WebsiteEventsHandler) createWebSession(ctx context.Context, trackerData *postgres_entity.WebTrackerEvents) (*postgres_entity.WebSession, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebsiteTrackerEventsHandler.createWebSession")
	defer span.Finish()
	tracing.TagComponentRest(span)

	query := postgres_entity.WebSession{
		Tenant:        trackerData.Tenant,
		IP:            trackerData.IP,
		VisitorID:     trackerData.VisitorID,
		Hostname:      trackerData.Hostname,
		Referrer:      &trackerData.Referrer,
		StartTime:     utils.Now(),
		LastEventType: trackerData.EventType,
		LastActivity:  utils.Now(),
		IsActive:      true,
	}

	params := h.parseURLParams(trackerData.Search)

	if params != nil {
		paramString, err := h.setReferrerQueryParams(ctx, params)
		if err != nil {
			tracing.TraceErr(span, err)
			return nil, err
		}
		if paramString != nil {
			query.QueryParams = paramString
		}
	}

	newSession, err := h.services.Repositories.PostgresRepositories.WebSessionRepository.Create(ctx, query)
	if err != nil {
		tracing.TraceErr(span, err)
		return nil, err
	}
	if newSession == nil {
		err = errors.New("unable to create new web session")
		tracing.TraceErr(span, err)
		return nil, err
	}
	return newSession, nil
}

func (h *WebsiteEventsHandler) setReferrerQueryParams(ctx context.Context, queryParams []ReferrerQueryParams) (*string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebsiteTrackerEventsHandler.SetReferrerQueryParams")
	defer span.Finish()

	bytes, err := json.Marshal(queryParams)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal QueryParam: %w", err)
	}
	results := string(bytes)
	return &results, nil
}

func (h *WebsiteEventsHandler) isTrustedIP(ctx context.Context, ipAddress string) bool {
	span, ctx := opentracing.StartSpanFromContext(ctx, "WebsiteTrackerEventsHandler.isTrustedIp")
	defer span.Finish()
	tracing.TagTenant(span, common.GetTenantFromContext(ctx))
	span.LogKV("ipAddress", ipAddress)

	ipThreats, err := h.services.CommonServices.VerifyService.Threats(ctx, ipAddress)
	if err != nil || ipThreats == nil {
		tracing.TraceErr(span, err)
		return true
	}

	return !ipThreats.IsThreat
}
