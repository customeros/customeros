import { action, observable, runInAction } from 'mobx';

export class TraceRecord {
  id: string;
  name: string;
  result: unknown;
  startTime: number;
  duration: number | null;
  parentId: string | null;
  children: TraceRecord[];
  expanded: boolean = false;
  attributes: Record<string, unknown> | null;

  constructor(
    name: string,
    startTime: number,
    parentId: string | null,
    attributes: Record<string, unknown> | null,
  ) {
    this.id = crypto.randomUUID();
    this.name = name;

    this.startTime = startTime;
    this.duration = null;
    this.parentId = parentId;
    this.children = [];
    this.attributes = attributes;
  }
}

export class Tracer {
  @observable static accessor enabled: boolean = true;
  @observable static accessor expandedTraces: Set<string> = new Set();
  @observable static accessor displayCallStack: boolean = false;
  @observable static accessor traces: TraceRecord[] = [];
  @observable private static accessor activeTraces: TraceRecord[] = [];

  private start: number = -1;
  @observable private accessor traceRecord = new TraceRecord('', 0, null, null);

  constructor(private name: string, attributes?: Record<string, unknown>) {
    if (!('tracer' in globalThis)) return;
    if (!Tracer.enabled) return;

    this.start = performance.now();

    this.traceRecord = new TraceRecord(
      name,
      this.start,
      Tracer.activeTraces.length > 0
        ? Tracer.activeTraces[Tracer.activeTraces.length - 1].id
        : null,
      attributes ?? null,
    );

    runInAction(() => {
      Tracer.traces.push(this.traceRecord);

      if (this.traceRecord.parentId) {
        const parent = Tracer.traces.find(
          (t) => t.id === this.traceRecord.parentId,
        );

        if (parent) {
          parent.children.push(this.traceRecord);
        }
      }
      Tracer.activeTraces.push(this.traceRecord);
    });
  }

  logArgs(args: Record<string, unknown>) {
    runInAction(() => {
      this.traceRecord.attributes = { ...this.traceRecord.attributes, ...args };
    });
  }

  end(result?: unknown) {
    if (!('tracer' in globalThis)) return;
    if (!Tracer.enabled) return;

    const duration = (performance.now() - this.start).toFixed(2);

    runInAction(() => {
      this.traceRecord.duration = parseFloat(duration);
      this.traceRecord.result = result;
      Tracer.activeTraces.pop();
    });

    return result;
  }

  public static span(name: string, attributes?: Record<string, unknown>) {
    return new Tracer(name, attributes);
  }

  public static async run<T>(
    name: string,
    attributes: Record<string, unknown>,
    fn: () => Promise<T>,
  ): Promise<T> {
    const trace = new Tracer(name, attributes);

    try {
      const result = await fn();

      return trace.end(result) as T;
    } catch (error) {
      trace.end(error);
      throw error;
    }
  }

  @action
  public static toggleExpand(traceId: string) {
    if (Tracer.expandedTraces.has(traceId)) {
      Tracer.expandedTraces.delete(traceId);
    } else {
      Tracer.expandedTraces.add(traceId);
    }
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

  @action
  public static clear() {
    Tracer.traces = [];
    Tracer.activeTraces = [];
  }
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
(globalThis as any).tracer = Tracer;
