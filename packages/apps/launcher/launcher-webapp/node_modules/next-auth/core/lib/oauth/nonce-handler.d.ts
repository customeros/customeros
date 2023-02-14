import type { InternalOptions } from "../../types";
import type { Cookie } from "../cookie";
/**
 * Returns nonce if the provider supports it
 * and saves it in a cookie */
export declare function createNonce(options: InternalOptions<"oauth">): Promise<undefined | {
    value: string;
    cookie: Cookie;
}>;
/**
 * Returns nonce from if the provider supports nonce,
 * and clears the container cookie afterwards.
 */
export declare function useNonce(nonce: string | undefined, options: InternalOptions<"oauth">): Promise<{
    value: string;
    cookie: Cookie;
} | undefined>;
