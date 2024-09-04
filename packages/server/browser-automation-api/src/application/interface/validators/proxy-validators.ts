import { body, param } from "express-validator";

export const getProxyValidator = [
  param("id")
    .exists({ checkFalsy: true })
    .isNumeric()
    .withMessage("Proxy id param must be a number."),
];

export const addProxyValidator = [
  body("url")
    .notEmpty()
    .withMessage("url field is required.")
    .isString()
    .withMessage("url field must be a string."),
  body("username")
    .notEmpty()
    .withMessage("username field is required.")
    .isString()
    .withMessage("username field must be a string."),
  body("password")
    .notEmpty()
    .withMessage("password field is required.")
    .isString()
    .withMessage("password field must be a string."),
];
