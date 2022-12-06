import Link from "next/link"
import { signIn, signOut, useSession } from "next-auth/react"
import styles from "./header.module.css"

// The approach used in this component shows how to build a sign in and sign out
// component that works on pages which support both client and server side
// rendering, and avoids any flash incorrect content on initial page load.
export default function Header() {
  const { data: session, status } = useSession()
  const loading = status === "loading"

  return (
    <header>
      <noscript>
        <style>{`.nojs-show { opacity: 1; top: 0; }`}</style>
      </noscript>

        <nav className="navbar navbar--fixed-top">
            <div className="navbar__inner">
                <div className="navbar__items">
                    <button aria-label="Toggle navigation bar" aria-expanded="false"
                            className="navbar__toggle clean-btn" type="button">
                        <svg width="30" height="30" viewBox="0 0 30 30" aria-hidden="true">
                            <path stroke="currentColor" strokeLinecap="round" strokeMiterlimit="10" strokeWidth="2"
                                  d="M4 7h22M4 15h22M4 23h22"></path>
                        </svg>
                    </button>
                    <a className="navbar__brand" target="_self" href="/">
                        <div className="navbar__logo">
                            <img src="../static/img/OpenlineLogoLightMode.svg" alt="Openline Logo" className="themedImage_ToTc themedImage--light_HNdA" />
                        </div>
                        <b className="navbar__title text--truncate">Openline</b></a>
                    <div className="navbar__item dropdown dropdown--hoverable"><a href="#" aria-haspopup="true"
                                                                                  aria-expanded="false" role="button"
                                                                                  className="navbar__link">Developers</a>
                        <ul className="dropdown__menu">
                            <li><a className="dropdown__link" href="https://openline.ai/docs">Guides</a></li>
                            <li><a className="dropdown__link" href="https://openline.ai//docs/reference">API reference</a></li>
                            <li><a className="dropdown__link" href="https://openline.ai//docs/contribute">Community contributions</a></li>
                        </ul>
                    </div>
                    <a className="navbar__item navbar__link" href="https://openline.ai/blog">Blog</a></div>
                <div className="navbar__items navbar__items--right">
                    <div className="button">
                            {!session && (
                                <>
                                    <a href={`/api/auth/signin`}
                                        onClick={(e) => {
                                            e.preventDefault()
                                            signIn()
                                        }}
                                    >
                                        Sign in
                                    </a>
                                </>
                            )}
                            {session?.user && (
                                <>
                                    <a
                                        href={`/api/auth/signout`}
                                        onClick={(e) => {
                                            e.preventDefault()
                                            signOut()
                                        }}
                                    >
                                        Sign out
                                    </a>
                                </>
                            )}
                    </div>
                </div>
            </div>
            <div role="presentation" className="navbar-sidebar__backdrop"></div>



        </nav>
    </header>
  )
}
