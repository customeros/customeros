import { action, observable } from 'mobx';

export class Tracer {
  @observable static accessor enabled: boolean = false;
  @observable static accessor displayCallStack: boolean = false;
  private start: number = -1;

  constructor(private name: string, attributes?: Record<string, unknown>) {
    if (!('tracer' in globalThis)) return;
    if (!Tracer.enabled) return;

    this.start = performance.now();

    // eslint-disable-next-line no-console
    console.groupCollapsed(`▶️ ${name}`);
    // eslint-disable-next-line no-console
    Tracer.displayCallStack && console.trace('🛠 Call stack:');
    attributes && this.logArgs(attributes);
  }

  logArgs(args: Record<string, unknown>) {
    console.info('📥 Attributes:', args);
  }

  end(result?: unknown) {
    if (!('tracer' in globalThis)) return;
    if (!Tracer.enabled) return;

    const duration = (performance.now() - this.start).toFixed(2);

    result && console.info('✅ Result:', result);
    console.info(`⏳ Duration: ${duration} ms`);
    // eslint-disable-next-line no-console
    console.groupEnd();

    return result;
  }

  public static span(name: string, attributes?: Record<string, unknown>) {
    return new Tracer(name, attributes);
  }

  @action
  public static enable() {
    Tracer.enabled = true;
  }

  @action
  public static disable() {
    Tracer.enabled = false;
  }

  @action
  public static showCallstack() {
    Tracer.displayCallStack = true;
  }

  @action
  public static hideCallstack() {
    Tracer.displayCallStack = false;
  }
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
(globalThis as any).tracer = Tracer;
