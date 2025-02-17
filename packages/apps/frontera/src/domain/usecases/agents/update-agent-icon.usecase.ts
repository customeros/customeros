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
      name: 'refresh-ccw-02',
      keywords: ['refresh'],
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
      name: 'container',
      keywords: ['container', 'box', 'storage', 'information'],
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
      keywords: ['beaker', 'science', 'laboratory', 'chemistry', 'experiment'],
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
      name: 'ruler',
      keywords: ['ruler', 'centimeter', 'measure'],
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
    {
      name: 'file-05',
      keywords: ['file', 'document', 'paper', 'page', 'new'],
    },
    {
      name: 'file-03',
      keywords: ['file', 'document', 'chart', 'statistics', 'graph'],
    },
    {
      name: 'clipboard',
      keywords: ['clipboard', 'paste', 'copy', 'board', 'list'],
    },
    {
      name: 'bank',
      keywords: ['bank', 'finance', 'building', 'money', 'institution'],
    },
    {
      name: 'bank-note-01',
      keywords: ['bank note', 'money', 'cash', 'currency', 'payment'],
    },
    {
      name: 'coins-02',
      keywords: ['coins', 'money', 'currency', 'finance', 'cash'],
    },
    {
      name: 'coins-stacked-01',
      keywords: ['coins', 'stack', 'money', 'savings', 'finance'],
    },
    {
      name: 'currency-bitcoin-circle',
      keywords: ['bitcoin', 'cryptocurrency', 'digital', 'currency', 'crypto'],
    },
    {
      name: 'currency-dollar-circle',
      keywords: ['dollar', 'usd', 'currency', 'money', 'finance'],
    },
    {
      name: 'currency-euro-circle',
      keywords: ['euro', 'eur', 'currency', 'money', 'finance'],
    },
    {
      name: 'currency-pound-circle',
      keywords: ['pound', 'gbp', 'sterling', 'currency', 'british'],
    },
    {
      name: 'currency-ruble-circle',
      keywords: ['ruble', 'rub', 'russian', 'currency', 'money'],
    },
    {
      name: 'currency-rupee-circle',
      keywords: ['rupee', 'inr', 'indian', 'currency', 'money'],
    },
    {
      name: 'currency-yen-circle',
      keywords: ['yen', 'jpy', 'japanese', 'currency', 'money'],
    },
    {
      name: 'diamond-01',
      keywords: ['diamond', 'jewel', 'precious', 'luxury', 'value'],
    },
    {
      name: 'gift-01',
      keywords: ['gift', 'present', 'package', 'surprise', 'box'],
    },
    {
      name: 'sale-03',
      keywords: ['sale', 'discount', 'offer', 'price', 'deal'],
    },
    {
      name: 'shopping-bag-01',
      keywords: ['shopping', 'bag', 'purchase', 'retail', 'store'],
    },
    {
      name: 'scales-01',
      keywords: ['scales', 'balance', 'justice', 'weight', 'measure'],
    },
    {
      name: 'activity',
      keywords: ['activity', 'pulse', 'health', 'monitor', 'heartbeat'],
    },
    {
      name: 'activity-heart',
      keywords: ['activity', 'pulse', 'health', 'monitor', 'heartbeat'],
    },
    {
      name: 'anchor',
      keywords: ['anchor', 'ship', 'marine', 'nautical', 'stability'],
    },
    {
      name: 'asterisk-02',
      keywords: ['asterisk', 'star', 'symbol', 'reference', 'mark'],
    },
    {
      name: 'bookmark',
      keywords: ['bookmark', 'save', 'favorite', 'marker', 'tag'],
    },
    {
      name: 'building-05',
      keywords: ['building', 'office', 'work', 'business', 'company'],
    },
    {
      name: 'cake',
      keywords: ['cake', 'birthday', 'celebration', 'party', 'dessert'],
    },
    {
      name: 'flame',
      keywords: ['flame', 'fire', 'hot', 'burn', 'energy'],
    },
    {
      name: 'home-02',
      keywords: ['home', 'house', 'building', 'residence', 'property'],
    },
    {
      name: 'home-smile',
      keywords: ['home', 'smile', 'happy', 'house', 'friendly'],
    },
    {
      name: 'loading-03',
      keywords: ['loading', 'spinner', 'wait', 'process', 'refresh'],
    },
    {
      name: 'medical-cross',
      keywords: ['medical', 'cross', 'health', 'hospital', 'emergency'],
    },
    {
      name: 'speedometer-03',
      keywords: ['speedometer', 'speed', 'gauge', 'measure', 'performance'],
    },
    {
      name: 'target-01',
      keywords: ['target', 'aim', 'goal', 'objective', 'focus'],
    },
    {
      name: 'target-03',
      keywords: ['target', 'bullseye', 'aim', 'precision', 'accuracy'],
    },
    {
      name: 'virus',
      keywords: ['virus', 'disease', 'infection', 'health', 'medical'],
    },
    {
      name: 'zap',
      keywords: ['zap', 'lightning', 'thunder', 'power', 'energy'],
    },
    {
      name: 'zap-circle',
      keywords: ['zap', 'lightning', 'power', 'energy', 'circle'],
    },
    {
      name: 'camera-lens',
      keywords: ['camera', 'lens', 'photography', 'focus', 'capture'],
    },
    {
      name: 'image-03',
      keywords: ['image', 'picture', 'photo', 'media', 'gallery'],
    },
    {
      name: 'image-02',
      keywords: ['image', 'picture', 'photo', 'media', 'gallery'],
    },

    {
      name: 'grid-03',
      keywords: ['grid', 'layout', 'pattern', 'organize', 'structure'],
    },
    {
      name: 'intersect-circle',
      keywords: ['intersect', 'circle', 'overlap', 'combine', 'merge'],
    },
    {
      name: 'layers-three-01',
      keywords: ['layers', 'stack', 'multiple', 'organize', 'depth'],
    },
    {
      name: 'key-01',
      keywords: ['key', 'lock', 'security', 'access', 'password'],
    },
    {
      name: 'fingerprint-04',
      keywords: [
        'fingerprint',
        'biometric',
        'security',
        'identity',
        'authentication',
      ],
    },
    {
      name: 'lock-01',
      keywords: ['lock', 'security', 'protection', 'private', 'secure'],
    },
    {
      name: 'shield-dollar',
      keywords: ['shield', 'dollar', 'protection', 'security', 'finance'],
    },
    {
      name: 'shield-01',
      keywords: ['shield', 'protection', 'security', 'defense', 'guard'],
    },
    {
      name: 'star-06',
      keywords: ['star', 'favorite', 'rating', 'bookmark', 'important'],
    },
    {
      name: 'calendar',
      keywords: ['calendar', 'date', 'schedule', 'time', 'event'],
    },
    {
      name: 'alarm-clock',
      keywords: ['alarm', 'clock', 'time', 'schedule', 'alert'],
    },
    {
      name: 'clock-check',
      keywords: ['clock', 'check', 'time', 'complete', 'done'],
    },
    {
      name: 'hourglass-02',
      keywords: ['hourglass', 'time', 'wait', 'loading', 'progress'],
    },
    {
      name: 'hourglass-03',
      keywords: ['hourglass', 'time', 'wait', 'loading', 'progress'],
    },
    {
      name: 'calendar-check-01',
      keywords: ['calendar', 'check', 'schedule', 'event', 'complete'],
    },
    {
      name: 'face-smile',
      keywords: ['face', 'smile', 'happy', 'emotion', 'positive'],
    },
    {
      name: 'face-frown',
      keywords: ['face', 'frown', 'sad', 'emotion', 'negative'],
    },
    {
      name: 'face-neutral',
      keywords: ['face', 'neutral', 'emotion', 'expression', 'emoticon'],
    },
    {
      name: 'cloud-01',
      keywords: ['cloud', 'weather', 'storage', 'data', 'sky'],
    },
    {
      name: 'cloud-raining-02',
      keywords: ['cloud', 'rain', 'weather', 'storm', 'precipitation'],
    },
    {
      name: 'cloud-sun-01',
      keywords: ['cloud', 'sun', 'weather', 'partly cloudy', 'day'],
    },
    {
      name: 'hurricane-01',
      keywords: ['hurricane', 'storm', 'weather', 'cyclone', 'disaster'],
    },
    {
      name: 'sun',
      keywords: ['sun', 'weather', 'day', 'light', 'bright'],
    },
    {
      name: 'sun-setting-01',
      keywords: ['sun', 'sunset', 'dusk', 'evening', 'weather'],
    },
    {
      name: 'wind-02',
      keywords: ['wind', 'weather', 'breeze', 'air', 'nature'],
    },
    {
      name: 'umbrella-03',
      keywords: ['umbrella', 'rain', 'weather', 'protection', 'shelter'],
    },
    {
      name: 'thermometer-cold',
      keywords: ['thermometer', 'cold', 'temperature', 'weather', 'freeze'],
    },
    {
      name: 'thermometer-warm',
      keywords: ['thermometer', 'warm', 'temperature', 'weather', 'heat'],
    },
    {
      name: 'scissors-02',
      keywords: ['scissors', 'cut', 'tool', 'edit', 'craft'],
    },
    {
      name: 'building-02',
      keywords: ['building', 'architecture', 'office', 'structure', 'company'],
    },
    {
      name: 'colors',
      keywords: ['colors', 'palette', 'design', 'art', 'creative'],
    },
    {
      name: 'feather',
      keywords: ['feather', 'write', 'edit', 'light', 'pen'],
    },
    {
      name: 'magic-wand-01',
      keywords: ['magic', 'wand', 'effect', 'wizard', 'spell'],
    },
    {
      name: 'pen-tool-02',
      keywords: ['pen', 'tool', 'design', 'draw', 'vector'],
    },
    {
      name: 'edit-03',
      keywords: ['edit', 'pencil'],
    },
    {
      name: 'coins-swap-01',
      keywords: ['coins', 'swap', 'exchange', 'currency', 'money'],
    },
    {
      name: 'shopping-cart-03',
      keywords: ['shopping', 'cart', 'buy', 'purchase', 'store'],
    },
    {
      name: 'heart-rounded',
      keywords: ['heart', 'love', 'like', 'favorite', 'emotion'],
    },
    {
      name: 'life-buoy-01',
      keywords: ['life buoy', 'help', 'support', 'safety', 'rescue'],
    },
    {
      name: 'life-buoy-02',
      keywords: ['life buoy', 'help', 'support', 'safety', 'rescue'],
    },
    {
      name: 'image-05',
      keywords: ['image', 'picture', 'photo', 'media', 'gallery'],
    },
    {
      name: 'bus',
      keywords: ['bus', 'transport', 'vehicle', 'travel', 'public'],
    },
    {
      name: 'car-01',
      keywords: ['car', 'vehicle', 'transport', 'automobile', 'travel'],
    },
    {
      name: 'flag-04',
      keywords: ['flag', 'banner', 'mark', 'country', 'symbol'],
    },
    {
      name: 'globe-05',
      keywords: ['globe', 'world', 'earth', 'international', 'global'],
    },
    {
      name: 'luggage-03',
      keywords: ['luggage', 'bag', 'travel', 'suitcase', 'baggage'],
    },
    {
      name: 'map-01',
      keywords: ['map', 'location', 'navigation', 'direction', 'geography'],
    },
    {
      name: 'marker-pin-04',
      keywords: ['marker', 'pin', 'location', 'map', 'place'],
    },
    {
      name: 'plane',
      keywords: ['plane', 'flight', 'travel', 'airplane', 'transport'],
    },
    {
      name: 'rocket-02',
      keywords: ['rocket', 'launch', 'startup', 'space', 'fast'],
    },
    {
      name: 'route',
      keywords: ['route', 'path', 'direction', 'navigation', 'journey'],
    },
    {
      name: 'train',
      keywords: ['train', 'transport', 'railway', 'travel', 'public'],
    },
    {
      name: 'truck-02',
      keywords: ['truck', 'transport', 'delivery', 'shipping', 'cargo'],
    },
    {
      name: 'battery-charging-01',
      keywords: ['battery', 'charging', 'power', 'energy', 'electric'],
    },
    {
      name: 'clapperboard',
      keywords: ['clapperboard', 'movie', 'film', 'cinema', 'video'],
    },
    {
      name: 'disc-02',
      keywords: ['disc', 'cd', 'dvd', 'music', 'storage'],
    },
    {
      name: 'gaming-pad-01',
      keywords: ['gaming', 'pad', 'controller', 'game', 'play'],
    },
    {
      name: 'music-note-01',
      keywords: ['music', 'note', 'sound', 'audio', 'melody'],
    },
    {
      name: 'play-circle',
      keywords: ['play', 'circle', 'media', 'video', 'start'],
    },
    {
      name: 'repeat-01',
      keywords: ['repeat', 'loop', 'cycle', 'reload', 'refresh'],
    },
    {
      name: 'repeat-03',
      keywords: ['repeat', 'loop', 'cycle', 'reload', 'refresh'],
    },
    {
      name: 'signal-01',
      keywords: ['signal', 'wifi', 'connection', 'network', 'wireless'],
    },
    {
      name: 'radar',
      keywords: ['signal', 'radar'],
    },
    {
      name: 'voicemail',
      keywords: ['voicemail', 'message', 'phone', 'audio', 'communication'],
    },
    {
      name: 'webcam-01',
      keywords: ['webcam', 'camera', 'video', 'conference', 'stream'],
    },
    {
      name: 'arrows-triangle',
      keywords: ['arrows', 'triangle', 'direction', 'move', 'navigation'],
    },
    {
      name: 'phone-hang-up',
      keywords: ['phone', 'hang up', 'call', 'end', 'telephone'],
    },
    {
      name: 'database-01',
      keywords: ['database', 'storage', 'data', 'server', 'information'],
    },
    {
      name: 'eraser',
      keywords: ['eraser', 'delete', 'remove', 'clear', 'edit'],
    },
    {
      name: 'hand',
      keywords: ['hand', 'gesture', 'touch', 'interact', 'human'],
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
