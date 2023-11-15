/* eslint-disable */
//TODO: fix any and remove eslint disable - this will be fixed then update mutation input type will be available
interface TimeToRenewalForm {
  renewalCycle: any | null;
  contractEnds: Date | null;
  serviceStarts: Date | null;
  contractSigned: Date | null;
}

export class ContractDTO implements TimeToRenewalForm {
  contractSigned: Date | null;
  contractEnds: Date | null;
  serviceStarts: Date | null;
  renewalCycle: any | null;

  constructor(data?: TimeToRenewalForm | null) {
    this.renewalCycle = data?.renewalCycle ? new Date(data.renewalCycle) : null;
    this.contractSigned = data?.contractSigned
      ? new Date(data.contractSigned)
      : null;
    this.contractEnds = data?.contractEnds ? new Date(data.contractEnds) : null;
    this.serviceStarts = data?.serviceStarts
      ? new Date(data.serviceStarts)
      : null;
  }

  static toForm(data?: TimeToRenewalForm | null): TimeToRenewalForm {
    return new ContractDTO(data);
  }

  static toPayload(data: TimeToRenewalForm): any {
    return {
      serviceStartedAt: data?.serviceStarts,
      signedAt: data?.contractSigned,
      contractEnds: data?.contractEnds,
      renewalCycle: data?.renewalCycle,
    };
  }
}
