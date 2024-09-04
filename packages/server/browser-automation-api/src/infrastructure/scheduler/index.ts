import { CronJob, CronJobParams } from "cron";

export type JobParams = CronJobParams & {
  runOnce?: boolean;
  onComplete?: () => void | Promise<void>;
};

export class Scheduler {
  private static instance: Scheduler;
  private jobs: Map<number | string, CronJob> = new Map();

  private constructor() {}

  public static getInstance(): Scheduler {
    if (!Scheduler.instance) {
      Scheduler.instance = new Scheduler();
    }
    return Scheduler.instance;
  }

  public schedule(id: number | string, jobParams: JobParams) {
    const { cronTime, onTick, runOnce, onComplete, start } = jobParams;
    const job = new CronJob(cronTime, onTick, onComplete, start);

    if (runOnce) {
      const teardown = () => {
        job.stop();
        this.jobs.delete(id);
      };
      job.addCallback(teardown);
    }

    this.jobs.set(id, job);
  }

  public startJobs() {
    this.jobs.forEach((job) => job.start());
  }

  public stopJobs() {
    this.jobs.forEach((job) => job.stop());
  }
}
