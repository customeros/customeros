import { Tracer } from '@infra/tracer';
import { RootStore } from '@store/root';
import { action, computed, observable } from 'mobx';

export const cooldownPeriods = [1, 4, 8, 12, 24, 999999];
export const cooldownPeriodsMap: Record<number, string> = {
  1: '1h',
  4: '4h',
  8: '8h',
  12: '12h',
  24: '24h',
  999999: `Don't notify again`,
};

export class AddSlackChannelUsecase {
  private root = RootStore.getInstance();
  @observable accessor isOpen = false;
  @observable accessor inputValue = '';
  @observable accessor selectedChannel = '';
  @observable accessor cooldownPeriod = 12;
  @observable private accessor _options = [
    { label: 'Canal1', value: '123' },
    { label: 'Canal2', value: '456' },
  ];

  constructor(private agentId: string) {
    this.toggle = this.toggle.bind(this);
    this.close = this.close.bind(this);
    this.selectChannel = this.selectChannel.bind(this);
    this.setInputValue = this.setInputValue.bind(this);
    this.enableSlack = this.enableSlack.bind(this);
    this.disableSlack = this.disableSlack.bind(this);
  }

  @computed
  get options() {
    return this._options.filter((option) =>
      option.label.toLowerCase().includes(this.inputValue.toLowerCase()),
    );
  }

  @computed
  get selectedOption() {
    return this._options.find(
      (option) => option.value === this.selectedChannel,
    );
  }

  @computed
  get selectedChannelName() {
    return this._options.find((option) => option.value === this.selectedChannel)
      ?.label;
  }

  @computed
  get isSlackEnabled() {
    return this.root.settings.slack.enabled;
  }

  @action
  toggle(open: boolean) {
    this.isOpen = open;
  }

  @action
  close() {
    const span = Tracer.span('AddSlackChannelUsecase.close', {
      isOpen: this.isOpen,
    });

    this.isOpen = false;

    span.end(`isOpen = ${this.isOpen}`);
  }

  @action
  selectChannel(channel: string) {
    const span = Tracer.span('AddSlackChannelUsecase.selectChannel', {
      channel,
    });

    this.selectedChannel = channel;
    this.close();

    span.end(`this.selectedChannel = ${this.selectedChannel}`);
  }

  @action
  setInputValue(value: string) {
    this.inputValue = value;
  }

  @action
  setCooldownPeriod(period: number) {
    this.cooldownPeriod = period;
  }

  enableSlack() {
    this.root.settings.slack.enableSync();
  }

  disableSlack() {
    this.root.settings.slack.disableSync();
  }
}
