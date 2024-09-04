import express from "express";
import bodyParser from "body-parser";
import swaggerUi from "swagger-ui-express";

import { logger } from "../logger";
import swaggerDocument from "./swagger.json";

export class Server {
  private port: number;
  public instance: ReturnType<typeof express>;
  private server: ReturnType<typeof this.instance.listen>;

  constructor(port: number) {
    this.port = port;
    this.instance = express();
    this.instance.use(bodyParser.json());
    this.instance.use(
      "/docs",
      swaggerUi.serve,
      swaggerUi.setup(swaggerDocument),
    );

    this.server = this.instance.listen(this.port, () => {
      logger.info(`Server is running on port ${this.port}`);
    });
  }

  public stop() {
    this.server.close(() => {
      console.log("Server has stopped");
    });
  }
}
