interface TimeToRenewalForm {
  contractEnds: Date | null;
  serviceStarts: Date | null;
  contractSigned: Date | null;
  contractRenews: Date | null;
}

export class ContractDTO implements TimeToRenewalForm {
  contractSigned: Date | null;
  contractEnds: Date | null;
  serviceStarts: Date | null;
  contractRenews: Date | null;

  constructor(data?: TimeToRenewalForm | null) {
    this.contractRenews = data?.contractRenews
      ? new Date(data.contractRenews)
      : null;
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

  static toPayload(data: TimeToRenewalForm): TimeToRenewalForm {
    return {
      contractEnds: data?.contractEnds,
      serviceStarts: data?.serviceStarts,
      contractSigned: data?.contractSigned,
      contractRenews: data?.contractRenews,
    };
  }
}
