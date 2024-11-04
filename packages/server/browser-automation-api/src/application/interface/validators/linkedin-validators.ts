import { body } from "express-validator";

export const connectValidators = [
  body("profileUrl")
    .notEmpty()
    .withMessage("profileUrl field is required")
    .isString()
    .withMessage("profileUrl field must be a string")
    .isURL()
    .withMessage("Invalid URL"),
  body("message")
    .optional()
    .isString()
    .withMessage("message field must be a string")
    .escape(),
  body("dryRun")
    .optional()
    .isBoolean()
    .withMessage("dryRun field must be a boolean"),
];

export const sendMessageValidators = [
  body("profileUrl")
    .notEmpty()
    .withMessage("profileUrl field is required")
    .isString()
    .withMessage("profileUrl field must be a string")
    .isURL()
    .withMessage("Invalid URL"),
  body("message")
    .notEmpty()
    .withMessage("message field is required")
    .isString()
    .withMessage("message field must be a string")
    .escape(),
  body("dryRun")
    .optional()
    .isBoolean()
    .withMessage("dryRun field must be a boolean"),
];

export const connectionStatusValidators = [
  body("profileUrl")
    .notEmpty()
    .withMessage("profileUrl field is required")
    .isString()
    .withMessage("profileUrl field must be a string")
    .isURL()
    .withMessage("Invalid URL"),
];

export const messagesRetrievalValidators = [
  body("profileUrl")
    .notEmpty()
    .withMessage("profileUrl field is required")
    .isString()
    .withMessage("profileUrl field must be a string")
    .isURL()
    .withMessage("Invalid URL"),
];
