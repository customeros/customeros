"use strict";
// Copyright © 2022 Ory Corp
// SPDX-License-Identifier: Apache-2.0
Object.defineProperty(exports, "__esModule", { value: true });
exports.config = void 0;
// @ory/integrations offers a package for integrating with Next.js.
var next_edge_1 = require("@ory/integrations/next-edge");
Object.defineProperty(exports, "config", { enumerable: true, get: function () { return next_edge_1.config; } });
// And create the Ory Network API "bridge".
exports.default = (0, next_edge_1.createApiHandler)({
    fallbackToPlayground: true,
});
