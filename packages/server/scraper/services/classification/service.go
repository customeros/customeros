package classification

func (s *NewWebSessionProducer) scrapeAndClassifyWebpage(ctx context.Context, url string) error {
	spans, ctx := telemetry.StartServiceSpan(ctx, "NewWebSessionProducer.scrapeAndClassifyWebpage")
	defer spans.Finish()

	content, err := s.webscraperService.Scrape(ctx, url)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	_, err = s.webscraperService.ClassifyWebpageCategory(ctx, url, &content)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	if content == "" {
		return nil
	}

	_, err = s.webscraperService.ClassifyContentStage(ctx, url, &content)
	if err != nil {
		spans.TraceError(err)
		return err
	}

	_, err = s.webscraperService.ClassifyWebpageTopics(ctx, url, &content)
	if err != nil {
		spans.TraceError(err)
		return err
	}
	return nil
}
