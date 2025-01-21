package agent_capability

//
// func (a *agentVisitorIDService) isNewCompanyVisit(ctx context.Context, tenant, domain *string) (bool, error) {
// 	span, ctx := tracing.StartTracerSpan(ctx, "AgentVisitorIDService.isNewCompanyVisit")
// 	defer span.Finish()
//
// 	if tenant == nil || domain == nil {
// 		err := errors.New("neither tenant or domain can be nil")
// 		tracing.TraceErr(span, err)
// 		span.LogKV("tenant", tenant)
// 		span.LogKV("domain", domain)
// 		return false, err
// 	}
//
// 	query := postgres_entity.WebSession{
// 		Tenant: *tenant,
// 		Domain: domain,
// 	}
//
// 	results, err := a.postgresRepositories.WebSessionRepository.FindAllSessions(ctx, query, nil)
// 	if err != nil {
// 		tracing.TraceErr(span, err)
// 		return false, err
// 	}
//
// 	if len(results) == 0 {
// 		return true, nil
// 	}
// 	return false, nil
// }
//
// func (a *agentVisitorIDService) isNewWebsiteVisitor(ctx context.Context, tenant, visitorId string) (bool, error) {
// 	span, ctx := tracing.StartTracerSpan(ctx, "AgentVisitorIDService.isNewWebsiteVisitor")
// 	defer span.Finish()
//
// 	query := postgres_entity.WebSession{
// 		Tenant:    tenant,
// 		VisitorID: visitorId,
// 	}
//
// 	results, err := a.postgresRepositories.WebSessionRepository.FindAllSessions(ctx, query, nil)
// 	if err != nil {
// 		tracing.TraceErr(span, err)
// 		return false, err
// 	}
//
// 	if len(results) == 0 {
// 		return true, nil
// 	}
// 	return false, nil
// }
//
// func (a *agentVisitorIDService) getUniquePageViews(ctx context.Context, pageViews []postgres_entity.WebTrackerEvents) ([]string, error) {
// 	span, ctx := tracing.StartTracerSpan(ctx, "AgentVisitorIDService.getUniquePageViews")
// 	defer span.Finish()
//
// 	// Use map to track unique pages
// 	uniquePageMap := make(map[string]struct{})
// 	for _, page := range pageViews {
// 		if page.Pathname != "" {
// 			uniquePageMap[page.Pathname] = struct{}{}
// 		}
// 	}
//
// 	// Convert map keys to slice
// 	uniquePages := make([]string, 0, len(uniquePageMap))
// 	for pathname := range uniquePageMap {
// 		uniquePages = append(uniquePages, pathname)
// 	}
//
// 	return uniquePages, nil
// }
//
// func (a *agentVisitorIDService) buildTimelineMessage(ctx context.Context, eventData *data_fields.WebsiteVisitEvent, isNewOrg, isNewVisitor bool) (*string, error) {
// 	span, ctx := tracing.StartTracerSpan(ctx, "AgentVisitorIDService.buildTimelineMessage")
// 	defer span.Finish()
//
// 	query := postgres_entity.WebTrackerEvents{
// 		Tenant:    eventData.Tenant,
// 		SessionID: eventData.SessionID,
// 		EventType: enum.WebTrackerPageView.String(),
// 	}
// 	pageViews, err := a.postgresRepositories.WebTrackerEventsRepository.FindAll(ctx, query, nil)
// 	if err != nil {
// 		tracing.TraceErr(span, err)
// 		return nil, err
// 	}
// 	if pageViews == nil {
// 		err := errors.New("could not locate any records for session id: " + eventData.SessionID)
// 		return nil, err
// 	}
//
// 	uniquePageViews, err := a.getUniquePageViews(ctx, pageViews)
// 	if err != nil {
// 		tracing.TraceErr(span, err)
// 		return nil, err
// 	}
//
// 	sessionDuration, err := a.calculateSessionDuration(ctx, &postgres_entity.WebSession{
// 		ID: eventData.SessionID,
// 	})
// 	if err != nil {
// 		tracing.TraceErr(span, err)
// 		return nil, err
// 	}
//
// 	// Build base message
// 	var baseMessage string
// 	switch {
// 	case isNewOrg:
// 		baseMessage = fmt.Sprintf("First Visit: %s", sessionDuration)
// 	case isNewVisitor:
// 		baseMessage = fmt.Sprintf("New Visitor: %s", sessionDuration)
// 	case !isNewVisitor:
// 		baseMessage = fmt.Sprintf("Repeat Visitor: %s", sessionDuration)
// 	default:
// 		baseMessage = fmt.Sprintf("Web Visitor: %s", sessionDuration)
// 	}
//
// 	// Append pages with indentation
// 	var fullMessage strings.Builder
// 	fullMessage.WriteString(baseMessage)
// 	for _, page := range uniquePageViews {
// 		fullMessage.WriteString(fmt.Sprintf("\n    • %s", page))
// 	}
//
// 	timelineMessage := fullMessage.String()
// 	return &timelineMessage, nil
// }
//
// func (a *agentVisitorIDService) calculateSessionDuration(ctx context.Context, session *postgres_entity.WebSession) (string, error) {
// 	span, ctx := tracing.StartTracerSpan(ctx, "AgentVisitorIDService.calculateSessionDuration")
// 	defer span.Finish()
//
// 	if session.EndTime.IsZero() {
// 		session, err := a.postgresRepositories.WebSessionRepository.FindSession(ctx, *session, nil)
// 		if err != nil {
// 			tracing.TraceErr(span, err)
// 			return "", err
// 		}
// 		if session == nil {
// 			return "", nil
// 		}
// 	}
//
// 	duration := session.EndTime.Sub(session.StartTime)
// 	minutes := duration.Minutes()
//
// 	switch {
// 	case minutes < 1:
// 		return "<1 minute", nil
// 	case minutes == 1:
// 		return "1 minute", nil
// 	default:
// 		return fmt.Sprintf("%.0f minutes", minutes), nil
// 	}
// }
//
// func (a *agentVisitorIDService) buildWebVisitorSlackNotification(ctx context.Context, session *postgres_entity.WebSession, orgID string) (*string, error) {
// 	span, ctx := opentracing.StartSpanFromContext(ctx, "AgentVisitorIDService.buildWebVisitorSlackNotification")
// 	defer span.Finish()
// 	tracing.SetDefaultListenerSpanTags(ctx, span)
//
// 	// get org data from global org table
// 	globalOrg, err := a.postgresRepositories.GlobalOrganizationRepository.GetByPrimaryDomain(ctx, *session.Domain)
// 	if err != nil {
// 		tracing.TraceErr(span, err)
// 		return nil, err
// 	}
//
// 	if globalOrg == nil {
// 		err = a.postgresRepositories.GlobalOrganizationWebsiteToProcessRepository.AddWebsiteToProcess(ctx, *session.Domain)
// 		if err != nil {
// 			tracing.TraceErr(span, err)
// 		}
// 	}
//
// 	// Build company info section
// 	var companyLines []string
// 	primaryDomain := *session.Domain
// 	if globalOrg != nil {
// 		primaryDomain = globalOrg.PrimaryDomain
// 	}
// 	website := "https://" + primaryDomain
// 	name := *session.Domain
// 	if globalOrg != nil {
// 		name = globalOrg.Name
// 	}
//
// 	companyLines = append(companyLines, fmt.Sprintf("<%s|*%s*>", website, name))
//
// 	if globalOrg != nil && globalOrg.Description != "" {
// 		companyLines = append(companyLines, fmt.Sprintf("%s", globalOrg.Description))
// 	}
//
// 	if website != "" {
// 		companyLines = append(companyLines, fmt.Sprintf("*Website:* <%s|%s>", website, primaryDomain))
// 	}
//
// 	if globalOrg != nil && globalOrg.LinkedInUrl != "" && globalOrg.LinkedInAlias != "" {
// 		companyLines = append(companyLines, fmt.Sprintf("*LinkedIn:* <%s|/%s>", globalOrg.LinkedInUrl, globalOrg.LinkedInAlias))
// 	}
//
// 	if globalOrg != nil && globalOrg.City != "" && globalOrg.CountryA2 != "" {
// 		companyLines = append(companyLines, fmt.Sprintf("*Location:* %s, %s", globalOrg.City, globalOrg.CountryA2))
// 	}
//
// 	if session.Referrer != nil {
// 		referrer := strings.TrimPrefix(*session.Referrer, "https://")
// 		referrer = strings.TrimPrefix(referrer, "http://")
// 		referrer = strings.TrimPrefix(referrer, "www.")
// 		referrer = strings.Trim(referrer, "/")
// 		companyLines = append(companyLines, fmt.Sprintf("*Source:* <%s|%s>", *session.Referrer, referrer))
// 	} else {
// 		companyLines = append(companyLines, "*Source:* Direct")
// 	}
//
// 	companyContent := strings.Join(companyLines, "\n")
//
// 	// Build session info section
// 	var sessionLines []string
//
// 	if session.EndTime != nil && !session.EndTime.IsZero() {
// 		duration, err := a.calculateSessionDuration(ctx, session)
// 		if err != nil {
// 			tracing.TraceErr(span, err)
// 			return nil, err
// 		}
//
// 		sessionLines = append(sessionLines, fmt.Sprintf("*Session Duration:* %s minutes", duration))
// 	}
//
// 	// Get page views
// 	query := postgres_entity.WebTrackerEvents{
// 		SessionID: session.ID,
// 		EventType: enum.WebTrackerPageView.String(),
// 	}
// 	pageViews, err := a.postgresRepositories.WebTrackerEventsRepository.FindAll(ctx, query, nil)
// 	if err != nil {
// 		tracing.TraceErr(span, err)
// 		return nil, err
// 	}
//
// 	uniquePages, err := a.getUniquePageViews(ctx, pageViews)
// 	if err != nil {
// 		tracing.TraceErr(span, err)
// 		return nil, err
// 	}
//
// 	sessionLines = append(sessionLines, fmt.Sprintf("*Pages Viewed:* %d", len(uniquePages)))
// 	for _, page := range uniquePages {
// 		sessionLines = append(sessionLines, fmt.Sprintf("• <%s%s|%s>", website, page, page))
// 	}
//
// 	sessionContent := strings.Join(sessionLines, "\n")
//
// 	// Handle logo accessory
// 	var logoAccessory string
// 	if globalOrg != nil && globalOrg.LogoUrl != "" {
// 		logoAccessory = fmt.Sprintf(`,
//            "accessory": {
//                "type": "image",
//                "image_url": "%s",
//                "alt_text": "%s logo"
//            }`, globalOrg.LogoUrl, name)
// 	}
//
// 	// Build the final layout
// 	layoutBlocks := fmt.Sprintf(`[
//        {
//            "type": "header",
//            "text": {
//                "type": "plain_text",
//                "text": "A visitor from %s is on your website",
//                "emoji": true
//            }
//        },
//        {
//            "type": "divider"
//        },
//        {
//            "type": "section",
//            "text": {
//                "type": "mrkdwn",
//                "text": "%s"
//            }%s
//        },
//        {
//            "type": "divider"
//        },
//        {
//            "type": "section",
//            "text": {
//                "type": "mrkdwn",
//                "text": "%s"
//            }
//        },
//        {
//            "type": "divider"
//        },
//        {
//            "type": "actions",
//            "elements": [
//                {
//                    "type": "button",
//                    "text": {
//                        "type": "plain_text",
//                        "text": "View in CustomerOS"
//                    },
//                    "url": "https://app.customeros.ai/organization/%s?tab=about",
//                    "value": "click_me_123",
//                    "action_id": "actionId-0"
//                }
//            ]
//        }
//    ]`,
// 		name,
// 		companyContent,
// 		logoAccessory,
// 		sessionContent,
// 		orgID)
//
// 	return &layoutBlocks, nil
// }
