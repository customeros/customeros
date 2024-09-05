import { body, param } from "express-validator";

export const getBrowserRunByIdValidators = [
  param("id")
    .exists({ checkFalsy: true })
    .isNumeric()
    .withMessage("BrowserAutomationRun id param must be a number."),
];
