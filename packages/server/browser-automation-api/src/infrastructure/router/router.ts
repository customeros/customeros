import express from "express";

export class Router {
  public instance: ReturnType<typeof express.Router>;

  constructor() {
    this.instance = express.Router();
  }
}
