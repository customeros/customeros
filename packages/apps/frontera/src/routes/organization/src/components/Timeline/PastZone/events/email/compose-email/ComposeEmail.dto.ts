type Option = {
  value: string;
  label: string;
  provider?: string;
};
export interface ComposeEmailDtoI {
  from: Option;
  subject: string;
  content: string;
  to: Array<Option>;
  cc: Array<Option>;
  bcc: Array<Option>;
  fromProvider: string;
  // files: Array<any>;
}

export class ComposeEmailDto implements ComposeEmailDtoI {
  from: Option;
  fromProvider: string;
  to: Array<Option>;
  cc: Array<Option>;
  bcc: Array<Option>;
  subject: string;
  content: string;
  // files: Array<any>;

  // TODO: type this correctly
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  constructor(data?: any) {
    this.from = data?.from || '';
    this.fromProvider = data?.fromProvider || '';
    this.to = data?.to || [];
    this.cc = data?.cc || [];
    this.bcc = data?.bcc || [];
    this.subject = data?.subject || '';
    this.content = data?.content || '';
    // this.files = data?.files || [];
  }
}
