import { CronJob, CronJobParams } from "cron";

import { logger } from "../logger";
import { ErrorParser, StandardError } from "@/util/error";

export type JobParams = CronJobParams & {
  runOnce?: boolean;
  onComplete?: () => void | Promise<void>;
  onTick: (completeTick: () => void) => void | Promise<void>;
};

export class Scheduler {
  private static instance: Scheduler;
  private jobs: Map<number | string, CronJob> = new Map();
  private runningJobs: Map<number | string, boolean> = new Map();

  private constructor() {}

  public static getInstance(): Scheduler {
    if (!Scheduler.instance) {
      Scheduler.instance = new Scheduler();
    }
    return Scheduler.instance;
  }

  public schedule(id: number | string, jobParams: JobParams) {
    const { cronTime, onTick, runOnce, onComplete, start } = jobParams;
    const job = new CronJob(
      cronTime,
      async () => {
        this.startTick(id);
        await onTick(() => this.completeTick(id));
      },
      onComplete,
      start,
    );

    if (runOnce) {
      // set a flag to indicate that the job is running
      job.addCallback(() => {
        this.runningJobs.set(id, true);
      });
      // stop the cron from scheduling the job again
      // and remove it from the running jobs
      job.addCallback(() => {
        job.stop();
        this.jobs.delete(id);
      });
    }

    this.jobs.set(id, job);
  }

  public startJobs() {
    this.jobs.forEach((job) => job.start());
  }

  public async stopJobs() {
    try {
      logger.info("Gracefully stopping jobs.", {
        source: "Scheduler",
      });
      let jobsStopped = 0;

      this.jobs.forEach((job) => {
        const nextInvocation = job.nextDate().diffNow().milliseconds;
        if (nextInvocation > 1000) {
          job.stop();
          jobsStopped++;
        }
      });
      logger.info(`${jobsStopped} scheduled jobs stopped.`, {
        source: "Scheduler",
      });

      return await this.checkRunningJobs();
    } catch (err) {
      const error = ErrorParser.parse(err);
      logger.error("Error stopping jobs", {
        error: error.message,
        details: error.details,
      });
      throw error;
    }
  }

  private checkRunningJobs(): Promise<boolean | undefined> {
    return new Promise((resolve, reject) => {
      const timeout = setTimeout(() => {
        reject(
          new StandardError({
            code: "INTERNAL_ERROR",
            message: "Timeout waiting for scheduled/running jobs to complete.",
            severity: "high",
          }),
        );
      }, 60 * 1000);

      const interval = setInterval(() => {
        logger.info(`Waiting on running jobs(${this.runningJobs.size}).`, {
          source: "Scheduler",
        });

        if (this.runningJobs.size === 0) {
          clearTimeout(timeout);
          clearInterval(interval);
          logger.info("All running jobs completed.", {
            source: "Scheduler",
          });
          resolve(true);
        }
      }, 1000);
    });
  }

  private startTick(id: number | string) {
    this.runningJobs.set(id, true);
  }

  private completeTick(id: number | string) {
    this.runningJobs.delete(id);
  }
}
