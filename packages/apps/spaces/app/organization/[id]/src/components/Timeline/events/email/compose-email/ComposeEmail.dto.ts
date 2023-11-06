type Option = {
  value: string;
  label: string;
};
export interface ComposeEmailDtoI {
  subject: string;
  content: string;
  to: Array<Option>;
  cc: Array<Option>;
  bcc: Array<Option>;
  // files: Array<any>;
}

export class ComposeEmailDto implements ComposeEmailDtoI {
  to: Array<Option>;
  cc: Array<Option>;
  bcc: Array<Option>;
  subject: string;
  content: string;
  // files: Array<any>;

  // TODO: type this correctly
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  constructor(data?: any) {
    this.to = data?.to || [];
    this.cc = data?.cc || [];
    this.bcc = data?.bcc || [];
    this.subject = data?.subject || '';
    this.content = data?.content || '';
    // this.files = data?.files || [];
  }

  // TODO: type this correctly
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  static toForm(data: any) {
    return new ComposeEmailDto(data);
  }

  static toPayload(data: ComposeEmailDtoI) {
    return {
      to: data.to,
      cc: data.cc,
      bcc: data.bcc,
      subject: data.subject,
      content: data.content,
      // files: data.files,
      // TODO: type this correctly
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
    } as any;
  }
}
