import axios from 'axios';


export interface OAuthToken {
    accessToken: string;
    refreshToken: string;
    expiresAt: Date;
    scope: string;
    providerAccountId: string;
    idToken: string
}

export interface SignInRequest {
    email: string;
    provider: string;
    oAuthToken: OAuthToken;
}


export function UserSignIn(
    data: SignInRequest,
): Promise<any> {
    return new Promise((resolve, reject) =>
        axios
            .post(process.env.USER_ADMIN_API_URL as string + "/signin", data,
                {
                    headers : {
                        'X-Openline-API-KEY': process.env.USER_ADMIN_API_KEY as string,
                    }
                })
            .then((response: any) => {
                if (response.data) {
                    resolve(response.data);
                } else {
                    reject(response.error);
                }
            })
            .catch((reason) => {
                reject(reason);
            }),
    );
}
