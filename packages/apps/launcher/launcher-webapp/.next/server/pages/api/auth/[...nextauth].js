"use strict";
/*
 * ATTENTION: An "eval-source-map" devtool has been used.
 * This devtool is neither made for production nor for readable output files.
 * It uses "eval()" calls to create a separate source file with attached SourceMaps in the browser devtools.
 * If you are trying to read the output file, select a different devtool (https://webpack.js.org/configuration/devtool/)
 * or disable the default devtool with "devtool: false".
 * If you are looking for production-ready output files, see mode: "production" (https://webpack.js.org/configuration/mode/).
 */
(() => {
var exports = {};
exports.id = "pages/api/auth/[...nextauth]";
exports.ids = ["pages/api/auth/[...nextauth]"];
exports.modules = {

/***/ "next-auth":
/*!****************************!*\
  !*** external "next-auth" ***!
  \****************************/
/***/ ((module) => {

module.exports = require("next-auth");

/***/ }),

/***/ "next-auth/providers/fusionauth":
/*!*************************************************!*\
  !*** external "next-auth/providers/fusionauth" ***!
  \*************************************************/
/***/ ((module) => {

module.exports = require("next-auth/providers/fusionauth");

/***/ }),

/***/ "(api)/./pages/api/auth/[...nextauth].ts":
/*!*****************************************!*\
  !*** ./pages/api/auth/[...nextauth].ts ***!
  \*****************************************/
/***/ ((__unused_webpack_module, __webpack_exports__, __webpack_require__) => {

eval("__webpack_require__.r(__webpack_exports__);\n/* harmony export */ __webpack_require__.d(__webpack_exports__, {\n/* harmony export */   \"authOptions\": () => (/* binding */ authOptions),\n/* harmony export */   \"default\": () => (__WEBPACK_DEFAULT_EXPORT__)\n/* harmony export */ });\n/* harmony import */ var next_auth__WEBPACK_IMPORTED_MODULE_0__ = __webpack_require__(/*! next-auth */ \"next-auth\");\n/* harmony import */ var next_auth__WEBPACK_IMPORTED_MODULE_0___default = /*#__PURE__*/__webpack_require__.n(next_auth__WEBPACK_IMPORTED_MODULE_0__);\n/* harmony import */ var next_auth_providers_fusionauth__WEBPACK_IMPORTED_MODULE_1__ = __webpack_require__(/*! next-auth/providers/fusionauth */ \"next-auth/providers/fusionauth\");\n/* harmony import */ var next_auth_providers_fusionauth__WEBPACK_IMPORTED_MODULE_1___default = /*#__PURE__*/__webpack_require__.n(next_auth_providers_fusionauth__WEBPACK_IMPORTED_MODULE_1__);\n\n\nconst authOptions = {\n    providers: [\n        next_auth_providers_fusionauth__WEBPACK_IMPORTED_MODULE_1___default()({\n            id: \"fusionauth\",\n            name: \"Openline\",\n            clientId: process.env.NEXTAUTH_OAUTH_CLIENT_ID,\n            clientSecret: process.env.NEXTAUTH_OAUTH_CLIENT_SECRET,\n            tenantId: process.env.NEXTAUTH_OAUTH_TENANT_ID,\n            issuer: process.env.NEXTAUTH_OAUTH_SERVER_URL,\n            client: {\n                authorization_signed_response_alg: \"HS256\",\n                id_token_signed_response_alg: \"HS256\"\n            }\n        })\n    ],\n    theme: {\n        colorScheme: \"dark\"\n    },\n    callbacks: {\n        async jwt ({ token  }) {\n            token.userRole = \"admin\";\n            return token;\n        }\n    }\n};\n/* harmony default export */ const __WEBPACK_DEFAULT_EXPORT__ = (next_auth__WEBPACK_IMPORTED_MODULE_0___default()(authOptions));\n//# sourceURL=[module]\n//# sourceMappingURL=data:application/json;charset=utf-8;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiKGFwaSkvLi9wYWdlcy9hcGkvYXV0aC9bLi4ubmV4dGF1dGhdLnRzLmpzIiwibWFwcGluZ3MiOiI7Ozs7Ozs7OztBQUFxRDtBQUNHO0FBRWpELE1BQU1FLGNBQStCO0lBQzFDQyxXQUFXO1FBQ1RGLHFFQUFVQSxDQUFDO1lBQ1RHLElBQUk7WUFDSkMsTUFBTTtZQUNOQyxVQUFVQyxRQUFRQyxHQUFHLENBQUNDLHdCQUF3QjtZQUM5Q0MsY0FBY0gsUUFBUUMsR0FBRyxDQUFDRyw0QkFBNEI7WUFDdERDLFVBQVVMLFFBQVFDLEdBQUcsQ0FBQ0ssd0JBQXdCO1lBQzlDQyxRQUFRUCxRQUFRQyxHQUFHLENBQUNPLHlCQUF5QjtZQUM3Q0MsUUFBUTtnQkFDTkMsbUNBQW1DO2dCQUNuQ0MsOEJBQThCO1lBQ2hDO1FBQ0Y7S0FDRDtJQUNEQyxPQUFPO1FBQ0xDLGFBQWE7SUFDZjtJQUNBQyxXQUFXO1FBQ1QsTUFBTUMsS0FBSSxFQUFFQyxNQUFLLEVBQUUsRUFBRTtZQUNuQkEsTUFBTUMsUUFBUSxHQUFHO1lBQ2pCLE9BQU9EO1FBQ1Q7SUFDRjtBQUNGLEVBQUM7QUFFRCxpRUFBZXZCLGdEQUFRQSxDQUFDRSxZQUFZQSxFQUFBIiwic291cmNlcyI6WyJ3ZWJwYWNrOi8vLy4vcGFnZXMvYXBpL2F1dGgvWy4uLm5leHRhdXRoXS50cz8yZThiIl0sInNvdXJjZXNDb250ZW50IjpbImltcG9ydCBOZXh0QXV0aCwgeyBOZXh0QXV0aE9wdGlvbnMgfSBmcm9tIFwibmV4dC1hdXRoXCJcbmltcG9ydCBGdXNpb25BdXRoIGZyb20gXCJuZXh0LWF1dGgvcHJvdmlkZXJzL2Z1c2lvbmF1dGhcIjtcblxuZXhwb3J0IGNvbnN0IGF1dGhPcHRpb25zOiBOZXh0QXV0aE9wdGlvbnMgPSB7XG4gIHByb3ZpZGVyczogW1xuICAgIEZ1c2lvbkF1dGgoe1xuICAgICAgaWQ6IFwiZnVzaW9uYXV0aFwiLFxuICAgICAgbmFtZTogXCJPcGVubGluZVwiLFxuICAgICAgY2xpZW50SWQ6IHByb2Nlc3MuZW52Lk5FWFRBVVRIX09BVVRIX0NMSUVOVF9JRCBhcyBzdHJpbmcsXG4gICAgICBjbGllbnRTZWNyZXQ6IHByb2Nlc3MuZW52Lk5FWFRBVVRIX09BVVRIX0NMSUVOVF9TRUNSRVQgYXMgc3RyaW5nLFxuICAgICAgdGVuYW50SWQ6IHByb2Nlc3MuZW52Lk5FWFRBVVRIX09BVVRIX1RFTkFOVF9JRCBhcyBzdHJpbmcsXG4gICAgICBpc3N1ZXI6IHByb2Nlc3MuZW52Lk5FWFRBVVRIX09BVVRIX1NFUlZFUl9VUkwsXG4gICAgICBjbGllbnQ6IHtcbiAgICAgICAgYXV0aG9yaXphdGlvbl9zaWduZWRfcmVzcG9uc2VfYWxnOiAnSFMyNTYnLFxuICAgICAgICBpZF90b2tlbl9zaWduZWRfcmVzcG9uc2VfYWxnOiAnSFMyNTYnXG4gICAgICB9XG4gICAgfSksXG4gIF0sXG4gIHRoZW1lOiB7XG4gICAgY29sb3JTY2hlbWU6IFwiZGFya1wiLFxuICB9LFxuICBjYWxsYmFja3M6IHtcbiAgICBhc3luYyBqd3QoeyB0b2tlbiB9KSB7XG4gICAgICB0b2tlbi51c2VyUm9sZSA9IFwiYWRtaW5cIlxuICAgICAgcmV0dXJuIHRva2VuXG4gICAgfSxcbiAgfSxcbn1cblxuZXhwb3J0IGRlZmF1bHQgTmV4dEF1dGgoYXV0aE9wdGlvbnMpXG4iXSwibmFtZXMiOlsiTmV4dEF1dGgiLCJGdXNpb25BdXRoIiwiYXV0aE9wdGlvbnMiLCJwcm92aWRlcnMiLCJpZCIsIm5hbWUiLCJjbGllbnRJZCIsInByb2Nlc3MiLCJlbnYiLCJORVhUQVVUSF9PQVVUSF9DTElFTlRfSUQiLCJjbGllbnRTZWNyZXQiLCJORVhUQVVUSF9PQVVUSF9DTElFTlRfU0VDUkVUIiwidGVuYW50SWQiLCJORVhUQVVUSF9PQVVUSF9URU5BTlRfSUQiLCJpc3N1ZXIiLCJORVhUQVVUSF9PQVVUSF9TRVJWRVJfVVJMIiwiY2xpZW50IiwiYXV0aG9yaXphdGlvbl9zaWduZWRfcmVzcG9uc2VfYWxnIiwiaWRfdG9rZW5fc2lnbmVkX3Jlc3BvbnNlX2FsZyIsInRoZW1lIiwiY29sb3JTY2hlbWUiLCJjYWxsYmFja3MiLCJqd3QiLCJ0b2tlbiIsInVzZXJSb2xlIl0sInNvdXJjZVJvb3QiOiIifQ==\n//# sourceURL=webpack-internal:///(api)/./pages/api/auth/[...nextauth].ts\n");

/***/ })

};
;

// load runtime
var __webpack_require__ = require("../../../webpack-api-runtime.js");
__webpack_require__.C(exports);
var __webpack_exec__ = (moduleId) => (__webpack_require__(__webpack_require__.s = moduleId))
var __webpack_exports__ = (__webpack_exec__("(api)/./pages/api/auth/[...nextauth].ts"));
module.exports = __webpack_exports__;

})();