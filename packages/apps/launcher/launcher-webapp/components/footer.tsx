import Link from "next/link"
import styles from "./footer.module.css"
import packageJSON from "../package.json"

export default function Footer() {
  return (
      <footer className="footer footer--dark">
        <div className="container container-fluid">
          <div className="row footer__links">
            <div className="col footer__col">
              <div className="footer__title">Docs</div>
              <ul className="footer__items clean-list">
                <li className="footer__item"><a className="footer__link-item" href="/docs">Getting Started</a></li>
              </ul>
            </div>
            <div className="col footer__col">
              <div className="footer__title">Community</div>
              <ul className="footer__items clean-list">
                <li className="footer__item"><a href="https://stackoverflow.com/questions/tagged/openline"
                                                target="_blank" rel="noopener noreferrer" className="footer__link-item">Stack
                  Overflow
                  <svg width="13.5" height="13.5" aria-hidden="true" viewBox="0 0 24 24"
                       className="iconExternalLink_nPIU">
                    <path fill="currentColor"
                          d="M21 13v10h-21v-19h12v2h-10v15h17v-8h2zm3-12h-10.988l4.035 4-6.977 7.07 2.828 2.828 6.977-7.07 4.125 4.172v-11z"></path>
                  </svg>
                </a></li>
                <li className="footer__item"><a
                    href="https://join.slack.com/t/openline-ai/shared_invite/zt-1i6umaw6c-aaap4VwvGHeoJ1zz~ngCKQ"
                    target="_blank" rel="noopener noreferrer" className="footer__link-item">Slack
                  <svg width="13.5" height="13.5" aria-hidden="true" viewBox="0 0 24 24"
                       className="iconExternalLink_nPIU">
                    <path fill="currentColor"
                          d="M21 13v10h-21v-19h12v2h-10v15h17v-8h2zm3-12h-10.988l4.035 4-6.977 7.07 2.828 2.828 6.977-7.07 4.125 4.172v-11z"></path>
                  </svg>
                </a></li>
                <li className="footer__item"><a href="https://twitter.com/openlineAI" target="_blank"
                                                rel="noopener noreferrer" className="footer__link-item">Twitter
                  <svg width="13.5" height="13.5" aria-hidden="true" viewBox="0 0 24 24"
                       className="iconExternalLink_nPIU">
                    <path fill="currentColor"
                          d="M21 13v10h-21v-19h12v2h-10v15h17v-8h2zm3-12h-10.988l4.035 4-6.977 7.07 2.828 2.828 6.977-7.07 4.125 4.172v-11z"></path>
                  </svg>
                </a></li>
              </ul>
            </div>
            <div className="col footer__col">
              <div className="footer__title">More</div>
              <ul className="footer__items clean-list">
                <li className="footer__item"><a href="https://github.com/openline-ai" target="_blank"
                                                rel="noopener noreferrer" className="footer__link-item">GitHub
                  <svg width="13.5" height="13.5" aria-hidden="true" viewBox="0 0 24 24"
                       className="iconExternalLink_nPIU">
                    <path fill="currentColor"
                          d="M21 13v10h-21v-19h12v2h-10v15h17v-8h2zm3-12h-10.988l4.035 4-6.977 7.07 2.828 2.828 6.977-7.07 4.125 4.172v-11z"></path>
                  </svg>
                </a></li>
                <li className="footer__item"><a className="footer__link-item" href="/legal">Legal</a></li>
              </ul>
            </div>
          </div>
          <div className="footer__bottom text--center">
            <div className="footer__copyright">Copyright © 2022 Openline Technologies, Inc. Built with ❤️ by the
              Openline community.
            </div>
          </div>
        </div>
      </footer>
  )
}
