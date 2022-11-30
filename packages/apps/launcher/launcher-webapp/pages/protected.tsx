import { useSession } from "next-auth/react"
import Layout from "../components/layout"
import AccessDenied from "../components/access-denied"

export default function ProtectedPage() {
  const { data: session } = useSession()

  // If no session exists, display access denied message
  if (!session) {
    return (
        <Layout>
          <AccessDenied />
        </Layout>
    )
  }

  // If session exists, display content
  return (
      <Layout>
        <h1>Protected Page</h1>
        <p>
          <strong>Link to some apps</strong>
        </p>
      </Layout>
  )
}
