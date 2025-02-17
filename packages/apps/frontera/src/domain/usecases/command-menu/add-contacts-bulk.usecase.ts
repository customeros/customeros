import { Tracer } from '@infra/tracer';
import { action, observable } from 'mobx';
import { CommonService } from '@domain/services/common/common.service';

export class AddContactsBulkUsecase {
  private commonService = new CommonService();

  constructor() {
    this.checkBrowserExtensionStatus();
  }

  @observable static accessor isBrowserExtensionEnabled: boolean = true;

  @action
  private async setIsBrowserExtensionEnabled(status: boolean) {
    AddContactsBulkUsecase.isBrowserExtensionEnabled = status;
  }

  private async checkBrowserExtensionStatus() {
    const span = Tracer.span(
      'AddContactsBulkUsecase.checkBrowserExtensionStatus',
    );
    const [res, err] = await this.commonService.getBrowserAutomationConfig();

    if (err || !res?.data?.data) {
      this.setIsBrowserExtensionEnabled(false);
      span.end();

      return;
    }

    const { data } = res.data;

    if (data.sessionStatus === 'VALID') {
      this.setIsBrowserExtensionEnabled(true);
    }

    span.end({
      browserConfigId: data.id,
      sessionStatus: data.sessionStatus,
    });
  }
}
