import { Transport } from '@store/transport';

export class MailboxesService {
  private static instance: MailboxesService;
  private transport: Transport;

  private constructor(transport: Transport) {
    this.transport = transport;
  }

  static getInstance(transport: Transport) {
    if (!MailboxesService.instance) {
      MailboxesService.instance = new MailboxesService(transport);
    }

    return MailboxesService.instance;
  }
}
