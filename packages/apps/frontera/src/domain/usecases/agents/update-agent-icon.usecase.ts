import Fuse from 'fuse.js';
import { Tracer } from '@infra/tracer.ts';
import { RootStore } from '@store/root.ts';
import { reaction, observable } from 'mobx';
import { AgentService } from '@domain/services/agent/agent.service.ts';

import { IconName } from '@ui/media/Icon';

export interface IconDefinition {
  name: IconName;
  keywords: string[];
}

export class UpdateAgentIconUsecase {
  private root = RootStore.getInstance();
  private service = new AgentService();
  private fuse: Fuse<IconDefinition>;

  @observable accessor searchQuery: string = '';
  @observable accessor isSaving: boolean = false;
  @observable accessor iconOptions: Array<IconName> = [];
  private iconDefinitions: IconDefinition[] = [
    {
      name: 'chart-breakout-circle',
      keywords: ['chart', 'breakout', 'circle', 'analytics', 'graph'],
    },
    {
      name: 'clock',
      keywords: ['time', 'schedule', 'watch', 'timer'],
    },
    {
      name: 'trend-up-01',
      keywords: ['trend', 'up', 'growth', 'increase', 'analytics'],
    },
    {
      name: 'bar-chart-05',
      keywords: ['bar', 'chart', 'graph', 'statistics', 'analytics'],
    },
    {
      name: 'bar-chart-06',
      keywords: ['bar', 'chart', 'graph', 'statistics', 'analytics'],
    },
    {
      name: 'bar-chart-circle-02',
      keywords: ['bar', 'chart', 'circle', 'statistics', 'graph'],
    },
    {
      name: 'bar-chart-circle-03',
      keywords: ['bar', 'chart', 'circle', 'statistics', 'graph'],
    },
    {
      name: 'chart-breakout-square',
      keywords: ['chart', 'breakout', 'square', 'analytics'],
    },
    {
      name: 'trend-down-01',
      keywords: ['trend', 'down', 'decrease', 'analytics'],
    },
    {
      name: 'presentation-chart-02',
      keywords: ['presentation', 'chart', 'slides', 'analytics'],
    },
    {
      name: 'presentation-chart-03',
      keywords: ['presentation', 'chart', 'slides', 'analytics'],
    },
    {
      name: 'bell-04',
      keywords: ['bell', 'notification', 'alert', 'ring'],
    },
    {
      name: 'bell-off-03',
      keywords: ['bell', 'off', 'mute', 'silent', 'notification'],
    },
    {
      name: 'thumbs-up',
      keywords: ['thumbs', 'up', 'like', 'positive', 'approval'],
    },
    {
      name: 'thumbs-down',
      keywords: ['thumbs', 'down', 'dislike', 'negative', 'disapproval'],
    },
    {
      name: 'announcement-03',
      keywords: ['announcement', 'megaphone', 'broadcast', 'alert'],
    },
    {
      name: 'infinity',
      keywords: ['infinity', 'endless', 'loop', 'continuous'],
    },
    {
      name: 'inbox-unread',
      keywords: ['inbox', 'unread', 'mail', 'message', 'notification'],
    },
    {
      name: 'mail-01',
      keywords: ['mail', 'email', 'message', 'envelope'],
    },
    {
      name: 'mail-05',
      keywords: ['mail', 'email', 'message', 'envelope'],
    },
    {
      name: 'message-chat-circle',
      keywords: ['message', 'chat', 'circle', 'communication'],
    },
    {
      name: 'message-text-circle-01',
      keywords: ['message', 'text', 'circle', 'chat'],
    },
    {
      name: 'phone-call-01',
      keywords: ['phone', 'call', 'telephone', 'contact'],
    },
    {
      name: 'send-03',
      keywords: ['send', 'paper plane', 'message', 'share'],
    },
    {
      name: 'send-01',
      keywords: ['send', 'paper plane', 'message', 'share'],
    },
    {
      name: 'brackets',
      keywords: ['brackets', 'code', 'programming', 'development'],
    },
    {
      name: 'code-02',
      keywords: ['code', 'programming', 'development', 'software'],
    },
    {
      name: 'code-square-02',
      keywords: ['code', 'square', 'programming', 'development'],
    },
    {
      name: 'cpu-chip-01',
      keywords: ['cpu', 'chip', 'processor', 'hardware', 'computer'],
    },
    {
      name: 'data',
      keywords: ['data', 'database', 'storage', 'information'],
    },
    {
      name: 'dataflow-03',
      keywords: ['dataflow', 'connection', 'network', 'flow'],
    },
    {
      name: 'dataflow-04',
      keywords: ['dataflow', 'connection', 'network', 'flow'],
    },
    {
      name: 'puzzle-piece-01',
      keywords: ['puzzle', 'piece', 'solution', 'game'],
    },
    {
      name: 'variable',
      keywords: ['variable', 'code', 'programming', 'math'],
    },
    {
      name: 'award-01',
      keywords: ['award', 'trophy', 'achievement', 'prize'],
    },
    {
      name: 'beaker-01',
      keywords: ['beaker', 'science', 'laboratory', 'chemistry'],
    },
    {
      name: 'beaker-02',
      keywords: ['beaker', 'science', 'laboratory', 'chemistry'],
    },
    {
      name: 'book-open-01',
      keywords: ['book', 'open', 'reading', 'education'],
    },
    {
      name: 'briefcase-02',
      keywords: ['briefcase', 'work', 'business', 'job'],
    },
    {
      name: 'certificate-01',
      keywords: ['certificate', 'achievement', 'diploma', 'award'],
    },
    {
      name: 'glasses-02',
      keywords: ['glasses', 'eyewear', 'vision', 'reading'],
    },
    {
      name: 'graduation-hat-01',
      keywords: ['graduation', 'hat', 'education', 'academic'],
    },
    {
      name: 'stand',
      keywords: ['stand', 'podium', 'platform', 'presentation'],
    },
    {
      name: 'telescope',
      keywords: ['telescope', 'astronomy', 'vision', 'search'],
    },
    {
      name: 'trophy-01',
      keywords: ['trophy', 'award', 'achievement', 'winner'],
    },
    {
      name: 'clipboard-check',
      keywords: ['clipboard', 'check', 'task', 'complete'],
    },
    {
      name: 'sticker-circle',
      keywords: ['sticker', 'circle', 'label', 'tag'],
    },
    {
      name: 'paperclip',
      keywords: ['paperclip', 'attachment', 'file', 'document'],
    },
  ];

  constructor(private id?: string) {
    this.execute = this.execute.bind(this);
    this.setSearchQuery = this.setSearchQuery.bind(this);
    this.getIconOptions = this.getIconOptions.bind(this);

    this.fuse = new Fuse(this.iconDefinitions, {
      keys: [
        { name: 'name', weight: 2 },
        { name: 'keywords', weight: 1 },
      ],
      threshold: 0.3,
      ignoreLocation: true,
      useExtendedSearch: true,
    });

    reaction(
      () => this.searchQuery,
      () => this.getIconOptions(),
      {
        fireImmediately: true,
      },
    );
  }

  get agent() {
    return this.id ? this.root.agents.getById(this.id) : null;
  }

  private getIconOptions() {
    if (!this.searchQuery.length) {
      this.iconOptions = this.iconDefinitions.map((result) => result.name);

      return;
    }

    const searchResults = this.fuse.search(this.searchQuery);

    this.iconOptions = searchResults.map((result) => result.item.name);
  }

  public setSearchQuery(query: string) {
    this.searchQuery = query;
  }

  async execute(selectedIcon?: IconName) {
    if (!selectedIcon) {
      console.error('UpdateAgentIconUsecase.execute: No icon  selected');

      return;
    }

    this.isSaving = true;

    const agent = this.agent;

    if (!agent) {
      console.error('UpdateAgentIconUsecase.execute: Could not find agent');
      this.isSaving = false;

      return;
    }

    const span = Tracer.span('UpdateAgentIconUsecase.execute', {
      payload: {
        agent,
        icon: selectedIcon,
      },
    });

    const prevIcon = agent.value.icon as IconName;

    agent.setIcon(selectedIcon);

    const [res, err] = await this.service.saveAgent(agent);

    if (err) {
      console.error(
        'UpdateAgentIconUsecase.execute: Could not update agent icon',
        err,
      );
      agent.setIcon(prevIcon);

      this.root.ui.toastError('We could not update agent icon', `${err}`);

      return;
    }

    if (res?.agent_Save) {
      agent.put(res.agent_Save);
    }

    this.isSaving = false;
    span.end();
  }
}
