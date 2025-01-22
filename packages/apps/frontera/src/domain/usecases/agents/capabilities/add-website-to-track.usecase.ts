import { action, observable } from 'mobx';

export class AddWebsiteToTrackUsecase {
  @observable accessor website: string = '';
  @observable accessor isOpen: boolean = true;

  constructor() {
    this.toggle = this.toggle.bind(this);
    this.open = this.open.bind(this);
    this.close = this.close.bind(this);
    this.execute = this.execute.bind(this);
    this.setWebsite = this.setWebsite.bind(this);
  }

  @action
  setWebsite(website: string) {
    this.website = website;
  }

  @action
  open() {
    this.isOpen = true;
  }

  @action
  toggle(open: boolean) {
    this.isOpen = open;
  }

  @action
  close() {
    this.isOpen = false;
  }

  execute() {
    // perform mutation to add
    alert('Added website to track ' + this.website);
    this.close();
  }
}
