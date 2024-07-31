import { match } from 'ts-pattern';
import { RootStore } from '@store/root';
import { Transport } from '@store/transport';
import { runInAction, makeAutoObservable } from 'mobx';

import {
  Note,
  Issue,
  Order,
  Action,
  Meeting,
  Analysis,
  PageView,
  LogEntry,
  TimelineEvent,
  InteractionEvent,
  InteractionSession,
} from '@graphql/types';

import mock from './mock.json';
import { NoteStore } from './Note/Note.store';
import { NotesStore } from './Note/Notes.store';
import { OrderStore } from './Order/Order.store';
import { IssueStore } from './Issues/Issue.store';
import { OrdersStore } from './Order/Orders.store';
import { IssuesStore } from './Issues/Issues.store';
import { ActionStore } from './Actions/Action.store';
import { ActionsStore } from './Actions/Actions.store';
import { MeetingStore } from './Meetings/Meeting.store';
import { AnalysisStore } from './Analyses/Analysis.store';
import { AnalysesStore } from './Analyses/Analyses.store';
import { MeetingsStore } from './Meetings/Meetings.store';
import { LogEntryStore } from './LogEntry/LogEntry.store';
import { PageViewStore } from './PageViews/PageView.store';
import { PageViewsStore } from './PageViews/PageViews.store';
import { LogEntriesStore } from './LogEntry/LogEntries.store';
import { TimelineEventsService } from './__service__/TimelineEvents.service';
import { InteractionEventStore } from './InteractionEvents/InteractionEvent.store';
import { InteractionEventsStore } from './InteractionEvents/InteractionEvents.store';
import { InteractionSessionStore } from './InteractionSessions/InteractionSession.store';
import { InteractionSessionsStore } from './InteractionSessions/InteractionsSessions.store';

type TimelineEventStore =
  | NoteStore
  | OrderStore
  | IssueStore
  | ActionStore
  | AnalysisStore
  | MeetingStore
  | PageViewStore
  | LogEntryStore
  | InteractionEventStore
  | InteractionSessionStore;

export class TimelineEventsStore {
  notes: NotesStore;
  orders: OrdersStore;
  issues: IssuesStore;
  actions: ActionsStore;
  analyses: AnalysesStore;
  meetings: MeetingsStore;
  pageViews: PageViewsStore;
  logEntries: LogEntriesStore;
  interactionEvents: InteractionEventsStore;
  interactionSessions: InteractionSessionsStore;
  private service: TimelineEventsService;
  isLoading = false;
  error: string | null = null;
  value: Map<string, TimelineEventStore[]> = new Map();

  constructor(public root: RootStore, public transport: Transport) {
    this.notes = new NotesStore(this.root, this.transport);
    this.orders = new OrdersStore(this.root, this.transport);
    this.issues = new IssuesStore(this.root, this.transport);
    this.actions = new ActionsStore(this.root, this.transport);
    this.analyses = new AnalysesStore(this.root, this.transport);
    this.meetings = new MeetingsStore(this.root, this.transport);
    this.pageViews = new PageViewsStore(this.root, this.transport);
    this.logEntries = new LogEntriesStore(this.root, this.transport);
    this.interactionEvents = new InteractionEventsStore(
      this.root,
      this.transport,
    );
    this.interactionSessions = new InteractionSessionsStore(
      this.root,
      this.transport,
    );
    this.service = TimelineEventsService.getInstance(this.transport);

    makeAutoObservable(this);
  }

  bootstrapTimeline(organizationId: string) {
    if (this.value.has(organizationId)) {
      return;
    }

    this.invalidateTimeline(organizationId);
  }

  async invalidateTimeline(organizationId: string) {
    if (this.root.demoMode) {
      runInAction(() => {
        const mockedTimeline = this.makeTimeline(
          (mock as unknown as Record<string, TimelineEvent[]>)[organizationId],
        );

        this.value.set(organizationId, mockedTimeline as TimelineEventStore[]);
      });

      return;
    }

    try {
      this.isLoading = true;

      const { organization } = await this.service.getTimeline({
        from: new Date(),
        organizationId,
        size: 100,
      });

      runInAction(() => {
        const timeline = this.makeTimeline(
          (organization?.timelineEvents as TimelineEvent[]) || [],
        );

        this.value.set(organizationId, timeline as TimelineEventStore[]);
      });
    } catch (e) {
      runInAction(() => {
        this.error = (e as Error).message;
      });
    } finally {
      runInAction(() => {
        this.isLoading = false;
      });
    }
  }

  getByOrganizationId(organizationId: string) {
    return this.value.get(organizationId);
  }

  private makeTimeline(timelineEvents: TimelineEvent[]) {
    return timelineEvents.map((event) => {
      return match(event)
        .with({ __typename: 'Note' }, (note) => {
          this.notes.load([note as Note]);

          return this.notes.value.get((note as Note).id);
        })
        .with({ __typename: 'Order' }, (order) => {
          this.orders.load([order as Order]);

          return this.orders.value.get((order as Order).id);
        })
        .with({ __typename: 'Issue' }, (issue) => {
          this.issues.load([issue as unknown as Issue]);

          return this.issues.value.get((issue as unknown as Issue).id);
        })
        .with({ __typename: 'Action' }, (action) => {
          this.actions.load([action as Action]);

          return this.actions.value.get((action as Action).id);
        })
        .with({ __typename: 'Analysis' }, (analysis) => {
          this.analyses.load([analysis as Analysis]);

          return this.analyses.value.get((analysis as Analysis).id);
        })
        .with({ __typename: 'Meeting' }, (meeting) => {
          this.meetings.load([meeting as Meeting]);

          return this.meetings.value.get((meeting as Meeting).id);
        })
        .with({ __typename: 'PageView' }, (pageView) => {
          this.pageViews.load([pageView as PageView]);

          return this.pageViews.value.get((pageView as PageView).id);
        })
        .with({ __typename: 'LogEntry' }, (logEntry) => {
          this.logEntries.load([logEntry as unknown as LogEntry]);

          return this.logEntries.value.get(
            (logEntry as unknown as LogEntry).id,
          );
        })
        .with({ __typename: 'InteractionEvent' }, (interactionEvent) => {
          this.interactionEvents.load([
            interactionEvent as unknown as InteractionEvent,
          ]);

          return this.interactionEvents.value.get(
            (interactionEvent as unknown as InteractionEvent).id,
          );
        })
        .with({ __typename: 'InteractionSession' }, (interactionSession) => {
          this.interactionSessions.load([
            interactionSession as InteractionSession,
          ]);

          return this.interactionSessions.value.get(
            (interactionSession as InteractionSession).id,
          );
        })
        .otherwise(() => {});
    });
  }
}
