import { relations } from "drizzle-orm/relations";
import { browserAutomationRuns, browserAutomationRunResults, browserAutomationRunErrors } from "./schema";

export const browserAutomationRunResultsRelations = relations(browserAutomationRunResults, ({one}) => ({
	browserAutomationRun: one(browserAutomationRuns, {
		fields: [browserAutomationRunResults.runId],
		references: [browserAutomationRuns.id]
	}),
}));

export const browserAutomationRunsRelations = relations(browserAutomationRuns, ({many}) => ({
	browserAutomationRunResults: many(browserAutomationRunResults),
	browserAutomationRunErrors: many(browserAutomationRunErrors),
}));

export const browserAutomationRunErrorsRelations = relations(browserAutomationRunErrors, ({one}) => ({
	browserAutomationRun: one(browserAutomationRuns, {
		fields: [browserAutomationRunErrors.runId],
		references: [browserAutomationRuns.id]
	}),
}));