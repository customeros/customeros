import Layout from "../components/layout"
import {inspect} from "util";
import styles from './index.module.css'
import Link from "next/link";
import {AbortController} from "next/dist/compiled/@edge-runtime/primitives/abort-controller";

export default function IndexPage() {
  return (
    <Layout>
        <div className="container" style={{ paddingTop: '5rem', paddingBottom: '5rem' }}>
            <h1 className={styles.title}>Welcome to Customer OS</h1>
        </div>
        <section>
            <div className="apps-container">
                <div className="apps-container-list">
                    <a href={`http://localhost:3006`}>
                        <div className="apps-container-list-item">
                            <h3 className={styles.subtitle}>Oasis</h3>
                            <p className={styles.centeredtext}>Short description here</p>
                        </div>
                    </a>
                    <a href={`http://localhost:3000`}>
                        <div className="apps-container-list-item">
                            <h3 className={styles.subtitle}>Contacts</h3>
                            <p className={styles.centeredtext}>Short description here</p>
                        </div>
                    </a>
                    <a href={`http://localhost:3000`}>
                        <div className="apps-container-list-item">
                            <h3 className={styles.subtitle}>Settings</h3>
                            <p className={styles.centeredtext}>N/A For the moment</p>
                        </div>
                    </a>
                    <a href={`http://localhost:9011`}>
                        <div className="apps-container-list-item">
                            <h3 className={styles.subtitle}>Auth</h3>
                            <p className={styles.centeredtext}>Short description here</p>
                        </div>
                    </a>
                </div>
            </div>
        </section>
        <div className="container" style={{ paddingTop: '5rem', paddingBottom: '5rem' }}>
            <h2 className={styles.subtitle}>Browse on Github</h2>
            <div className={styles.centeredtext}>
                <a href="https://github.com/openline-ai" rel="noreferrer" target="_blank">
                    <img src='../static/img/GithubButton.png' width={101} height={101} alt="Github Logo" />
                </a>
            </div>
        </div>
   </Layout>

  )
}
